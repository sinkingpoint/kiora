package model

import (
	"bytes"
	"encoding/json"
	"regexp"
)

type Matcher struct {
	Label      string `json:"label"`
	Value      string `json:"value"`
	IsRegex    bool   `json:"isRegex"`
	IsNegative bool   `json:"isNegative"`
	regex      *regexp.Regexp
}

func LabelValueRegexMatcher(label string, regex string) (Matcher, error) {
	r, err := regexp.Compile(regex)
	if err != nil {
		return Matcher{}, err
	}

	return Matcher{
		Label:   label,
		Value:   regex,
		IsRegex: true,
		regex:   r,
	}, nil
}

func LabelValueEqualMatcher(label string, value string) Matcher {
	return Matcher{
		Label: label,
		Value: value,
	}
}

func (m *Matcher) Negate() *Matcher {
	m.IsNegative = !m.IsNegative
	return m
}

func (m *Matcher) UnmarshalJSON(b []byte) error {
	raw := struct {
		Label      string `json:"label"`
		Value      string `json:"value"`
		IsRegex    bool   `json:"isRegex"`
		IsNegative bool   `json:"isNegative"`
	}{}

	d := json.NewDecoder(bytes.NewReader(b))
	d.DisallowUnknownFields()
	if err := d.Decode(&raw); err != nil {
		return err
	}

	m.Label = raw.Label
	m.Value = raw.Value
	m.IsRegex = raw.IsRegex
	m.IsNegative = raw.IsNegative

	if m.IsRegex {
		regex, err := regexp.Compile(m.Value)
		if err != nil {
			return err
		}

		m.regex = regex
	}

	return nil
}

func (m *Matcher) MatchesAlert(alert *Alert) bool {
	if _, ok := alert.Labels[m.Label]; !ok {
		return false
	}

	result := false
	if m.IsRegex {
		result = m.regex.MatchString(alert.Labels[m.Label])
	} else {
		result = alert.Labels[m.Label] == m.Value
	}

	if m.IsNegative {
		return !result
	} else {
		return result
	}
}
