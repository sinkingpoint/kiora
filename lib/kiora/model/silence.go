package model

import (
	"fmt"
	"math"
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"
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

func NewSilence(creator, comment string, matchers []Matcher, startTime, endTime time.Time) (Silence, error) {
	silence := Silence{
		ID:        uuid.New().String(),
		Creator:   creator,
		Comment:   comment,
		Matchers:  matchers,
		StartTime: startTime,
		EndTime:   endTime,
	}

	return silence, silence.validate()
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

func (s *Silence) Fields() map[string]any {
	return map[string]any{
		"__id__":        s.ID,
		"__creator__":   s.Creator,
		"__comment__":   s.Comment,
		"__starts_at__": s.StartTime,
		"__ends_at__":   s.EndTime,
		"__duration__":  s.EndTime.Sub(s.StartTime),
	}
}

func (s *Silence) Field(name string) (any, error) {
	switch name {
	case "__id__":
		return s.ID, nil
	case "__creator__":
		return s.Creator, nil
	case "__comment__":
		return s.Comment, nil
	case "__starts_at__":
		return s.StartTime, nil
	case "__ends_at__":
		return s.EndTime, nil
	case "__duration__":
		if s.EndTime.IsZero() {
			return time.Duration(math.MaxInt64), nil
		}

		return s.EndTime.Sub(s.StartTime), nil
	}

	return "", fmt.Errorf("silence %q doesn't exist", name)
}
