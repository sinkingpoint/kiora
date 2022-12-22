package model

import (
	"bytes"
	"encoding/json"
	"errors"
	"time"
)

type Matcher struct {
	Label    string `json:"label"`
	Value    string `json:"value"`
	Negative bool   `json:"negative"`
	Regex    bool   `json:"regex"`
}

type Silence struct {
	Creator   string    `json:"creator"`
	Comment   string    `json:"comment"`
	StartTime time.Time `json:"startTime"`
	EndTime   time.Time `json:"endTime"`
	Matchers  []Matcher `json:"matchers"`
}

func (s *Silence) UnmarshalJSON(b []byte) error {
	rawSilence := struct {
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

	defaultTime := time.Time{}
	if rawSilence.StartTime == defaultTime {
		return errors.New("missing start time in alert")
	}

	if rawSilence.EndTime == defaultTime {
		return errors.New("missing end time in alert")
	}

	if !rawSilence.EndTime.After(rawSilence.StartTime) {
		return errors.New("start time must be before end time")
	}

	if len(rawSilence.Matchers) == 0 {
		return errors.New("silence must have matchers")
	}

	s.Creator = rawSilence.Creator
	s.Comment = rawSilence.Comment
	s.StartTime = rawSilence.StartTime
	s.EndTime = rawSilence.EndTime
	s.Matchers = rawSilence.Matchers

	return nil
}
