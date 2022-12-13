package model

import "github.com/sinkingpoint/kiora/internal/dto/kioraproto"

// deserializeStringMapFromProto takes a Map[Text, Text] from a proto definition
// and deserialises it into a map[string]string which we can actually work with
func deserializeStringMapFromProto(m kioraproto.Map) (map[string]string, error) {
	values := make(map[string]string)

	entries, err := m.Entries()
	if err != nil {
		return nil, err
	}

	for i := 0; i < entries.Len(); i++ {
		entry := entries.At(i)
		ptr, err := entry.Key()
		if err != nil {
			return nil, err
		}

		key := ptr.Text()

		ptr, err = entry.Value()
		if err != nil {
			return nil, err
		}

		value := ptr.Text()

		values[key] = value
	}

	return values, nil
}
