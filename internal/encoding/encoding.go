package encoding

import "encoding/json"

type Encoder interface {
	Marshal(o any) ([]byte, error)
}

var encodingTable = map[string]Encoder{
	"json": JSONEncoder{},
}

func LookupEncoding(name string) Encoder {
	return encodingTable[name]
}

type JSONEncoder struct{}

func (j JSONEncoder) Marshal(o any) ([]byte, error) {
	return json.Marshal(o)
}
