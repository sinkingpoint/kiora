package model

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"
)

type Silence struct {
	// ID is the unique identifier of the silence.
	ID string `json:"id"`

	// By is the user that created the silence.
	Creator string `json:"creator"`

	// Comment is a comment about the silence.
	Comment string `json:"comment"`

	// StartTime is the time at which the silence starts.
	StartTime time.Time `json:"startsAt"`

	// EndTime is the time at which the silence ends.
	EndTime time.Time `json:"endsAt"`

	// Matchers is a list of matchers that must all match an alert for it to be silenced.
	Matchers []Matcher `json:"matchers"`
}

func (s *Silence) validate() error {
	if s.StartTime.IsZero() {
		return errors.New("silence is missing a start time")
	}

	if !s.EndTime.IsZero() && s.EndTime.Before(s.StartTime) {
		return errors.New("end time is before start time")
	}

	// NOTE: this precludes the ability for a silence to match all alerts, which might be a valid use case.
	// But if you get here trying to do that, please don't.
	if len(s.Matchers) == 0 {
		return errors.New("silence must have at least one matcher")
	}

	return nil
}

func (s *Silence) UnmarshalJSON(b []byte) error {
	rawSilence := struct {
		Creator   string    `json:"creator"`
		Comment   string    `json:"comment"`
		Matchers  []Matcher `json:"matchers"`
		StartTime time.Time `json:"startsAt"`
		EndTime   time.Time `json:"endsAt"`
	}{}

	if err := json.Unmarshal(b, &rawSilence); err != nil {
		return err
	}

	s.ID = uuid.New().String()
	s.Creator = rawSilence.Creator
	s.Comment = rawSilence.Comment
	s.Matchers = rawSilence.Matchers
	s.StartTime = rawSilence.StartTime
	s.EndTime = rawSilence.EndTime

	return s.validate()
}
