package model

import (
	"bytes"
	"sort"

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

// Subset returns a new Labels object with only the keys specified in labelNames.
func (s Labels) Subset(labelNames ...string) Labels {
	labels := Labels{}
	for _, key := range labelNames {
		labels[key] = s[key]
	}

	return labels
}

func (s Labels) Bytes() []byte {
	keys := []string{}
	for k := range s {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	buf := bytes.Buffer{}
	for _, k := range keys {
		buf.Write([]byte(k))
		buf.Write(hashSep)
		buf.Write([]byte(s[k]))
	}

	return buf.Bytes()
}

func (s Labels) Equal(other Labels) bool {
	if len(s) != len(other) {
		return false
	}

	for k, v := range s {
		if other[k] != v {
			return false
		}
	}

	return true
}
