package encoding

import (
	"encoding/json"
)

type Encoder interface {
	Marshal(o any) ([]byte, error)
}

var encodingTable = map[string]Encoder{
	"json":        NewJSONEncoder(false),
	"json-pretty": NewJSONEncoder(true),
}

func LookupEncoding(name string) Encoder {
	return encodingTable[name]
}

type JSONEncoder struct {
	Pretty bool
}

func NewJSONEncoder(pretty bool) JSONEncoder {
	return JSONEncoder{
		Pretty: pretty,
	}
}

func (j JSONEncoder) Marshal(o any) ([]byte, error) {
	if j.Pretty {
		return json.MarshalIndent(o, "", "  ")
	} else {
		return json.Marshal(o)
	}
}
