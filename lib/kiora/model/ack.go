package model

// AlertAcknowledgement is the metadata provided when an operator acknowledges an alert.
type AlertAcknowledgement struct {
	By      string `json:"creator"`
	Comment string `json:"comment"`
}
