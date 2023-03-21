package model

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/sinkingpoint/kiora/internal/stubs"
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
		ID        string    `json:"id"`
		Creator   string    `json:"creator"`
		Comment   string    `json:"comment"`
		Matchers  []Matcher `json:"matchers"`
		StartTime time.Time `json:"startsAt"`
		EndTime   time.Time `json:"endsAt"`
	}{}

	decoder := json.NewDecoder(bytes.NewReader(b))
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&rawSilence); err != nil {
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

func (s *Silence) IsActive() bool {
	return s.StartTime.Before(stubs.Time.Now()) && (s.EndTime.IsZero() || s.EndTime.After(stubs.Time.Now()))
}

func (s *Silence) Matches(l Labels) bool {
	for _, matcher := range s.Matchers {
		if !matcher.Matches(l) {
			return false
		}
	}

	return true
}

func (s *Silence) Field(name string) (string, error) {
	switch name {
	case "creator":
		return s.Creator, nil
	case "comment":
		return s.Comment, nil
	case "startsAt":
		return s.StartTime.Format(time.RFC3339), nil
	case "endsAt":
		return s.EndTime.Format(time.RFC3339), nil
	}

	return "", fmt.Errorf("silence %q doesn't exist", name)
}
