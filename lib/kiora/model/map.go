package model

import (
	"github.com/cespare/xxhash"
)

var hashSep = []byte{'\xff'}

type LabelsHash = uint64

// Labels is a utility type encapsulating a map[string]string that can be hashed.
type Labels map[string]string

// Hash takes an xxhash64 across all the labels in the map.
func (s Labels) Hash() LabelsHash {
	hash := xxhash.New()

	for k, v := range s {
		hash.Sum([]byte(k))
		hash.Sum(hashSep)
		hash.Sum([]byte(v))
	}

	return hash.Sum64()
}
