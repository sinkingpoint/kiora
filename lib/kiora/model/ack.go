package model

import "fmt"

// AlertAcknowledgement is the metadata provided when an operator acknowledges an alert.
type AlertAcknowledgement struct {
	Creator string `json:"creator"`
	Comment string `json:"comment"`
}

func (a *AlertAcknowledgement) Fields() map[string]any {
	return map[string]any{
		"__creator__": a.Creator,
		"__comment__": a.Comment,
	}
}

func (a *AlertAcknowledgement) Field(name string) (any, error) {
	switch name {
	case "__creator__":
		return a.Creator, nil
	case "__comment__":
		return a.Comment, nil
	default:
		return "", fmt.Errorf("field %q doesn't exist", name)
	}
}
