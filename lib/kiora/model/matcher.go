package model

import (
	"bytes"
	"encoding/json"
	"errors"
	"regexp"
	"strings"
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

func (m *Matcher) UnmarshalText(raw string) error {
	var parts []string

	if strings.Contains(raw, "=~") {
		parts = strings.Split(raw, "=~")
		m.IsRegex = true
		m.IsNegative = false
	} else if strings.Contains(raw, "!~") {
		parts = strings.Split(raw, "!~")
		m.IsRegex = true
		m.IsNegative = true
	} else if strings.Contains(raw, "!=") {
		parts = strings.Split(raw, "!=")
		m.IsRegex = false
		m.IsNegative = true
	} else {
		parts = strings.Split(raw, "=")
		m.IsRegex = false
		m.IsNegative = false
	}

	if len(parts) != 2 {
		return errors.New("invalid matcher")
	}

	m.Label = parts[0]
	m.Value = parts[1]
	return nil
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

func (m *Matcher) Matches(labels Labels) bool {
	if _, ok := labels[m.Label]; !ok {
		return false
	}

	result := false
	if m.IsRegex {
		result = m.regex.MatchString(labels[m.Label])
	} else {
		result = labels[m.Label] == m.Value
	}

	if m.IsNegative {
		return !result
	} else {
		return result
	}
}
