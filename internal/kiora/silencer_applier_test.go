package kiora_test

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/sinkingpoint/kiora/internal/kiora"
	"github.com/sinkingpoint/kiora/lib/kiora/model"
	"github.com/sinkingpoint/kiora/mocks/mock_kioradb"
	"github.com/stretchr/testify/assert"
)

// Check that alerts that are already silenced just skip.
func TestSilencer_AlreadySilenced(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	db := mock_kioradb.NewMockDB(ctrl)
	alert := model.Alert{
		Labels: map[string]string{
			"foo": "bar",
		},
		Status: model.AlertStatusSilenced,
	}

	applier := kiora.NewSilenceApplier()
	assert.NoError(t, applier.ProcessAlert(context.Background(), nil, db, &alert, &alert))
}

// Test that an alert that matches a silenced gets pushed down the pipeline with the Status set to Silenced.
func TestSilencer_Silences(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	db := mock_kioradb.NewMockDB(ctrl)
	alert := model.Alert{
		Labels: map[string]string{
			"foo": "bar",
		},
		Status: model.AlertStatusFiring,
	}

	silence := model.Silence{
		Matchers: []model.Matcher{
			&model.LabelValueEqualMatcher{
				Label: "foo",
				Value: "bar",
			},
		},
	}

	db.EXPECT().GetSilences(gomock.Any(), map[string]string{
		"foo": "bar",
	}).Times(1).Return([]model.Silence{silence}, nil)

	applier := kiora.NewSilenceApplier()
	assert.NoError(t, applier.ProcessAlert(context.Background(), nil, db, nil, &alert))

	assert.Equal(t, model.AlertStatusSilenced, alert.Status)
}
