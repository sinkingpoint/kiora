package serf

import (
	"context"
	"math/rand"
	"net"
	"os"
	"strconv"
	"sync"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/memberlist"
	"github.com/hashicorp/serf/serf"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/sinkingpoint/kiora/internal/clustering"
	"github.com/sinkingpoint/kiora/internal/clustering/serf/messages"
	"github.com/sinkingpoint/kiora/lib/kiora/model"
	"github.com/sinkingpoint/msgpack/v5"
	"go.opentelemetry.io/otel"
)

var _ = clustering.Broadcaster(&SerfBroadcaster{})

// randomNodeName returns a random, 16 char long node name to use when one isn't given.
func randomNodeName() string {
	var letterRunes = []rune("1234567890abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	b := make([]rune, 16)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

type Config struct {
	ListenURL         string
	BootstrapPeers    []string
	NodeName          string
	ClustererDelegate clustering.ClustererDelegate
	EventDelegate     clustering.EventDelegate

	Logger zerolog.Logger
}

func DefaultConfig() *Config {
	return &Config{
		ListenURL:      "localhost:4279",
		BootstrapPeers: []string{},
		NodeName:       randomNodeName(),
		Logger:         zerolog.New(os.Stdout).Level(zerolog.DebugLevel).With().Str("component", "serf").Logger(),
	}
}

type SerfBroadcaster struct {
	conf *Config

	shutdownOnce sync.Once
	serfCh       chan serf.Event
	serf         *serf.Serf
}

// NewSerfBroadcaster constructs a SerfBroadcaster with the given config, storing models in the given DB.
func NewSerfBroadcaster(conf *Config) (*SerfBroadcaster, error) {
	host, portStr, err := net.SplitHostPort(conf.ListenURL)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get host and port from Serf listen URL")
	}

	port, err := strconv.ParseUint(portStr, 10, 16) // 16 bits because that's the range of port numbers.
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse port from Serf listen URL")
	}

	// TODO(cdouch): 16 here is arbitrary. Should benchmark it.
	serfCh := make(chan serf.Event, 16)

	memberlistConf := memberlist.DefaultLANConfig()
	memberlistConf.BindAddr = host
	memberlistConf.BindPort = int(port)
	memberlistConf.Name = conf.NodeName

	serfConfig := serf.DefaultConfig()
	serfConfig.MemberlistConfig = memberlistConf
	serfConfig.EventCh = serfCh
	serfConfig.NodeName = conf.NodeName
	serfConfig.MaxQueueDepth = 1 << 16      // 64k
	serfConfig.UserEventSizeLimit = 1 << 12 // 4k
	serfConfig.Logger = &HCLogger{
		conf.Logger,
	}

	serf, err := serf.Create(serfConfig)
	if err != nil {
		return nil, errors.Wrap(err, "failed to init serf")
	}

	return &SerfBroadcaster{
		conf:         conf,
		shutdownOnce: sync.Once{},
		serf:         serf,
		serfCh:       serfCh,
	}, nil
}

func (s *SerfBroadcaster) Name() string {
	return "serf"
}

// Run provides a BackgroundService that processes events that come in via the Serf cluster.
func (s *SerfBroadcaster) Run(ctx context.Context) error {
	defer close(s.serfCh)

	if _, err := s.serf.Join(s.conf.BootstrapPeers, false); err != nil {
		return errors.Wrap(err, "failed to join bootstrap peers")
	}

	for {
		select {
		case <-ctx.Done():
			s.Shutdown()
			return nil
		case event := <-s.serfCh:
			s.processEvent(context.Background(), event)
		}
	}
}

func (s *SerfBroadcaster) Shutdown() {
	s.shutdownOnce.Do(func() {
		if err := s.serf.Leave(); err != nil {
			log.Error().Err(err).Msg("failed to leave Serf cluster cleanly")
		}

		s.serf.Shutdown() // nolint:errcheck
	})
}

func (s *SerfBroadcaster) processEvent(ctx context.Context, event serf.Event) {
	switch ev := event.(type) {
	case serf.UserEvent:
		s.processUserEvent(ctx, ev)
	case serf.MemberEvent:
		s.processMemberEvent(ctx, ev)
	default:
		return
	}
}

func (s *SerfBroadcaster) processMemberEvent(ctx context.Context, ev serf.MemberEvent) {
	if s.conf.ClustererDelegate == nil {
		return
	}

	switch ev.Type {
	case serf.EventMemberJoin:
		for _, member := range ev.Members {
			addr := net.JoinHostPort(member.Addr.String(), strconv.Itoa(int(member.Port)))
			s.conf.ClustererDelegate.AddNode(member.Name, addr)
		}
	case serf.EventMemberLeave, serf.EventMemberFailed:
		for _, member := range ev.Members {
			s.conf.ClustererDelegate.RemoveNode(member.Name)
		}
	}
}

func (s *SerfBroadcaster) processUserEvent(ctx context.Context, u serf.UserEvent) {
	ctx, span := otel.Tracer("").Start(ctx, "SerfBroadcaster.processUserEvent")
	defer span.End()

	// If we don't have an EventDelegate, there's nothing to handle these events.
	if s.conf.EventDelegate == nil {
		return
	}

	msg := messages.GetMessage(u.Name)
	if msg == nil {
		log.Error().Str("message name", u.Name).Msg("unhandled message type")
		return
	}

	if err := msgpack.Unmarshal(u.Payload, msg); err != nil {
		log.Err(err).Str("message name", u.Name).Msg("failed to unmarshal message")
		return
	}

	var err error
	switch msg := msg.(type) {
	case *messages.Alert:
		s.conf.EventDelegate.ProcessAlert(ctx, msg.Alert)
	case *messages.Acknowledgement:
		s.conf.EventDelegate.ProcessAlertAcknowledgement(ctx, msg.AlertID, msg.Acknowledgement)
	case *messages.Silence:
		s.conf.EventDelegate.ProcessSilence(ctx, msg.Silence)
	default:
		log.Error().Str("message name", u.Name).Msg("unhandled message type")
		return
	}

	if err != nil {
		log.Error().Str("message name", u.Name).Msg("failed to process message")
	}
}

func (s *SerfBroadcaster) broadcast(ctx context.Context, msg messages.Message) error {
	_, span := otel.Tracer("").Start(ctx, "SerfBroadcaster.broadcast")
	defer span.End()

	bytes, err := msgpack.Marshal(&msg)
	if err != nil {
		return errors.Wrap(err, "failed to marshal alerts")
	}

	if err := s.serf.UserEvent(msg.Name(), bytes, false); err != nil {
		return errors.Wrap(err, "failed to send broadcast")
	}

	return nil
}

// BroadcastAlerts sends alerts over the Serf gossip channel to the cluster.
func (s *SerfBroadcaster) BroadcastAlerts(ctx context.Context, alerts ...model.Alert) error {
	ctx, span := otel.Tracer("").Start(ctx, "SerfBroadcaster.BroadcastAlerts")
	defer span.End()

	var broadcastError error

	// Note: We break the alerts into individual messages in order to attempt to avoid Serf message size limits.
	for _, a := range alerts {
		msg := messages.Alert{
			Alert: a,
		}

		if err := s.broadcast(ctx, &msg); err != nil {
			broadcastError = multierror.Append(broadcastError, err)
		}
	}

	return broadcastError
}

func (s *SerfBroadcaster) BroadcastAlertAcknowledgement(ctx context.Context, alertID string, ack model.AlertAcknowledgement) error {
	msg := messages.Acknowledgement{
		AlertID:         alertID,
		Acknowledgement: ack,
	}

	return s.broadcast(ctx, &msg)
}

func (s *SerfBroadcaster) BroadcastSilences(ctx context.Context, silences ...model.Silence) error {
	var broadcastError error

	for _, silence := range silences {
		msg := messages.Silence{
			Silence: silence,
		}

		if err := s.broadcast(ctx, &msg); err != nil {
			broadcastError = multierror.Append(broadcastError, err)
		}
	}

	return broadcastError
}
