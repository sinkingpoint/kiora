package serf

import (
	"context"
	"net"
	"strconv"

	"github.com/hashicorp/memberlist"
	"github.com/hashicorp/serf/serf"
	"github.com/pkg/errors"
	"github.com/sinkingpoint/kiora/internal/clustering"
	"github.com/sinkingpoint/kiora/internal/clustering/serf/messages"
	"github.com/sinkingpoint/kiora/lib/kiora/kioradb"
	"github.com/sinkingpoint/kiora/lib/kiora/model"
	"github.com/vmihailenco/msgpack/v5"
)

var _ = clustering.Broadcaster(&SerfBroadcaster{})

type Config struct {
	ListenURL      string
	BootstrapPeers []string
}

func DefaultConfig() *Config {
	return &Config{
		ListenURL:      "localhost:4279",
		BootstrapPeers: []string{},
	}
}

type SerfBroadcaster struct {
	conf *Config

	db     kioradb.DB
	serfCh chan serf.Event
	serf   *serf.Serf
}

func NewSerfBroadcaster(conf *Config, db kioradb.DB) (*SerfBroadcaster, error) {
	host, portStr, err := net.SplitHostPort(conf.ListenURL)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get host and port from Serf listen URL")
	}

	port, err := strconv.ParseInt(portStr, 10, 16)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse port from Serf listen URL")
	}

	serfCh := make(chan serf.Event, 16)

	memberlistConf := memberlist.DefaultLANConfig()
	memberlistConf.BindAddr = host
	memberlistConf.BindPort = int(port)

	serfConfig := serf.DefaultConfig()
	serfConfig.MemberlistConfig = memberlistConf
	serfConfig.EventCh = serfCh

	serf, err := serf.Create(serfConfig)
	if err != nil {
		return nil, errors.Wrap(err, "failed to init serf")
	}

	return &SerfBroadcaster{
		conf:   conf,
		serf:   serf,
		serfCh: serfCh,
		db:     db,
	}, nil
}

func (s *SerfBroadcaster) Run(ctx context.Context) error {
	defer close(s.serfCh)

	if _, err := s.serf.Join(s.conf.BootstrapPeers, false); err != nil {
		return errors.Wrap(err, "failed to join bootstrap peers")
	}

	for {
		select {
		case <-ctx.Done():
			return nil
		case event := <-s.serfCh:
			s.processEvent(event)
		}
	}
}

func (s *SerfBroadcaster) processEvent(event serf.Event) {

}

func (s *SerfBroadcaster) BroadcastAlerts(ctx context.Context, alerts ...model.Alert) error {
	msg := messages.AlertMessage{}
	bytes, err := msgpack.Marshal(&msg)
	if err != nil {
		return errors.Wrap(err, "failed to marshal alerts")
	}

	return s.serf.UserEvent(msg.Name(), bytes, false)
}
