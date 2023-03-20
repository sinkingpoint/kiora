package model

import "fmt"

// AlertAcknowledgement is the metadata provided when an operator acknowledges an alert.
type AlertAcknowledgement struct {
	By      string `json:"creator"`
	Comment string `json:"comment"`
}

func (a *AlertAcknowledgement) Field(name string) (string, error) {
	switch name {
	case "from", "creator", "author", "by":
		return a.By, nil
	case "comment":
		return a.Comment, nil
	default:
		return "", fmt.Errorf("field %q doesn't exist", name)
	}
}
