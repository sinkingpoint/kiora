package raft

import (
	"github.com/sinkingpoint/kiora/internal/dto/kioraproto"
	"github.com/sinkingpoint/kiora/lib/kiora/model"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// newPostAlertsRaftLogMessage construcsta a RaftLogMessage that calls the PostAlerts
// method with the given alerts
func newPostAlertsRaftLogMessage(alerts ...model.Alert) *kioraproto.RaftLogMessage {
	protoAlerts := []*kioraproto.Alert{}
	for _, a := range alerts {
		protoAlerts = append(protoAlerts, &kioraproto.Alert{
			Labels:      a.Labels,
			Annotations: a.Annotations,
			Status:      a.Status.MapToProto(),
			StartTime:   timestamppb.New(a.StartTime),
			EndTime:     timestamppb.New(a.TimeOutDeadline),
		})
	}

	return &kioraproto.RaftLogMessage{
		Log: &kioraproto.RaftLogMessage_Alerts{
			Alerts: &kioraproto.PostAlertsMessage{
				Alerts: protoAlerts,
			},
		},
	}
}

// newPostSilencesRaftLogMessage takes a slice of model.Silences and packages
// them into a protobuf message that can be put into the raft log.
func newPostSilencesRaftLogMessage(silences ...model.Silence) *kioraproto.RaftLogMessage {
	protoSilences := []*kioraproto.Silence{}
	for _, silence := range silences {
		matchers := make([]*kioraproto.Matcher, 0, len(silence.Matchers))
		for _, m := range silence.Matchers {
			proto := m.MarshalProto()
			if proto != nil {
				matchers = append(matchers, proto)
			}
		}

		protoSilences = append(protoSilences, &kioraproto.Silence{
			ID:        silence.ID,
			Creator:   silence.Creator,
			Comment:   silence.Comment,
			StartTime: timestamppb.New(silence.StartTime),
			EndTime:   timestamppb.New(silence.EndTime),
			Matchers:  matchers,
		})
	}

	return &kioraproto.RaftLogMessage{
		Log: &kioraproto.RaftLogMessage_Silences{
			Silences: &kioraproto.PostSilencesRequest{
				Silences: protoSilences,
			},
		},
	}
}
