package model

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/sinkingpoint/kiora/internal/dto/kioraproto"
)

type Matcher struct {
	Label    string `json:"label"`
	Value    string `json:"value"`
	Negative bool   `json:"negative"`
	Regex    bool   `json:"regex"`
}

type Silence struct {
	ID        string    `json:"id"`
	Creator   string    `json:"creator"`
	Comment   string    `json:"comment"`
	StartTime time.Time `json:"startTime"`
	EndTime   time.Time `json:"endTime"`
	Matchers  []Matcher `json:"matchers"`
}

// validate returns an error if the silence fails any data validation checks,
// such as the start and end time being off.
func (s *Silence) validate() error {
	defaultTime := time.Time{}
	if s.StartTime == defaultTime {
		return errors.New("missing start time in silence")
	}

	if s.EndTime == defaultTime {
		return errors.New("missing end time in silence")
	}

	if !s.EndTime.After(s.StartTime) {
		return errors.New("start time must be before end time")
	}

	if len(s.Matchers) == 0 {
		return errors.New("silence must have matchers")
	}

	return nil
}

func (s *Silence) UnmarshalJSON(b []byte) error {
	rawSilence := struct {
		ID        string    `json:"id"`
		Creator   string    `json:"creator"`
		Comment   string    `json:"comment"`
		StartTime time.Time `json:"startTime"`
		EndTime   time.Time `json:"endTime"`
		Matchers  []Matcher `json:"matchers"`
	}{}

	decoder := json.NewDecoder(bytes.NewReader(b))
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&rawSilence); err != nil {
		return err
	}

	if rawSilence.ID != "" {
		s.ID = rawSilence.ID
	} else {
		s.ID = newSilenceID()
	}

	s.Creator = rawSilence.Creator
	s.Comment = rawSilence.Comment
	s.StartTime = rawSilence.StartTime
	s.EndTime = rawSilence.EndTime
	s.Matchers = rawSilence.Matchers

	return s.validate()
}

func newSilenceID() string {
	id := uuid.New()
	return id.String()
}

// DeserializeFromProto creates a model.Silence from a proto silence
func (s *Silence) DeserializeFromProto(proto *kioraproto.Silence) error {
	if proto.ID == "" {
		s.ID = newSilenceID()
	} else if _, err := uuid.Parse(proto.ID); err == nil {
		s.ID = proto.ID
	} else {
		return fmt.Errorf("got an id in the proto that wasn't valid: %q", proto.ID)
	}

	s.Creator = proto.Creator
	s.Comment = proto.Comment
	s.StartTime = proto.StartTime.AsTime()
	s.EndTime = proto.EndTime.AsTime()
	s.Matchers = make([]Matcher, 0, len(proto.Matchers))

	for _, matcher := range proto.Matchers {
		s.Matchers = append(s.Matchers, Matcher{
			Label:    matcher.Key,
			Value:    matcher.Value,
			Regex:    matcher.Regex,
			Negative: matcher.Negative,
		})
	}

	return s.validate()
}
