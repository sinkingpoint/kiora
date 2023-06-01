package unmarshal

// Secret wraps a string in a stringer that redacts the value.
type Secret string

func (s Secret) String() string {
	return "<redacted>"
}
