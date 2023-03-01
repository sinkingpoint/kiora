package raft

import (
	"github.com/sinkingpoint/kiora/internal/dto/kioraproto"
	"github.com/sinkingpoint/kiora/lib/kiora/model"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// newPostAlertsRaftLogMessage construcsta a RaftLogMessage that calls the PostAlerts
// method with the given alerts
func newPostAlertsRaftLogMessage(alerts ...model.Alert) *kioraproto.KioraLogMessage {
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

	return &kioraproto.KioraLogMessage{
		Log: &kioraproto.KioraLogMessage_Alerts{
			Alerts: &kioraproto.PostAlertsMessage{
				Alerts: protoAlerts,
			},
		},
	}
}
