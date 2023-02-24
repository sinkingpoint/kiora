package model

import (
	"bytes"

	"github.com/cespare/xxhash"
)

var hashSep = []byte{'\xff'}

type LabelsHash = uint64

// Labels is a utility type encapsulating a map[string]string that can be hashed.
type Labels map[string]string

// Hash takes an xxhash64 across all the labels in the map.
func (s Labels) Hash() LabelsHash {
	hash := xxhash.New()

	hash.Write(s.Bytes())

	return hash.Sum64()
}

func (s Labels) Bytes() []byte {
	buf := bytes.Buffer{}
	for k, v := range s {
		buf.Write([]byte(k))
		buf.Write(hashSep)
		buf.Write([]byte(v))
	}

	return buf.Bytes()
}
