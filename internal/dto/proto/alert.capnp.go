// Code generated by capnpc-go. DO NOT EDIT.

package dto

import (
	capnp "capnproto.org/go/capnp/v3"
	text "capnproto.org/go/capnp/v3/encoding/text"
	schemas "capnproto.org/go/capnp/v3/schemas"
)

type Alert capnp.Struct

// Alert_TypeID is the unique identifier for the type Alert.
const Alert_TypeID = 0xc5154448f10c0d22

func NewAlert(s *capnp.Segment) (Alert, error) {
	st, err := capnp.NewStruct(s, capnp.ObjectSize{DataSize: 24, PointerCount: 2})
	return Alert(st), err
}

func NewRootAlert(s *capnp.Segment) (Alert, error) {
	st, err := capnp.NewRootStruct(s, capnp.ObjectSize{DataSize: 24, PointerCount: 2})
	return Alert(st), err
}

func ReadRootAlert(msg *capnp.Message) (Alert, error) {
	root, err := msg.Root()
	return Alert(root.Struct()), err
}

func (s Alert) String() string {
	str, _ := text.Marshal(0xc5154448f10c0d22, capnp.Struct(s))
	return str
}

func (s Alert) EncodeAsPtr(seg *capnp.Segment) capnp.Ptr {
	return capnp.Struct(s).EncodeAsPtr(seg)
}

func (Alert) DecodeFromPtr(p capnp.Ptr) Alert {
	return Alert(capnp.Struct{}.DecodeFromPtr(p))
}

func (s Alert) ToPtr() capnp.Ptr {
	return capnp.Struct(s).ToPtr()
}
func (s Alert) IsValid() bool {
	return capnp.Struct(s).IsValid()
}

func (s Alert) Message() *capnp.Message {
	return capnp.Struct(s).Message()
}

func (s Alert) Segment() *capnp.Segment {
	return capnp.Struct(s).Segment()
}
func (s Alert) Labels() (Map, error) {
	p, err := capnp.Struct(s).Ptr(0)
	return Map(p.Struct()), err
}

func (s Alert) HasLabels() bool {
	return capnp.Struct(s).HasPtr(0)
}

func (s Alert) SetLabels(v Map) error {
	return capnp.Struct(s).SetPtr(0, capnp.Struct(v).ToPtr())
}

// NewLabels sets the labels field to a newly
// allocated Map struct, preferring placement in s's segment.
func (s Alert) NewLabels() (Map, error) {
	ss, err := NewMap(capnp.Struct(s).Segment())
	if err != nil {
		return Map{}, err
	}
	err = capnp.Struct(s).SetPtr(0, capnp.Struct(ss).ToPtr())
	return ss, err
}

func (s Alert) Annotations() (Map, error) {
	p, err := capnp.Struct(s).Ptr(1)
	return Map(p.Struct()), err
}

func (s Alert) HasAnnotations() bool {
	return capnp.Struct(s).HasPtr(1)
}

func (s Alert) SetAnnotations(v Map) error {
	return capnp.Struct(s).SetPtr(1, capnp.Struct(v).ToPtr())
}

// NewAnnotations sets the annotations field to a newly
// allocated Map struct, preferring placement in s's segment.
func (s Alert) NewAnnotations() (Map, error) {
	ss, err := NewMap(capnp.Struct(s).Segment())
	if err != nil {
		return Map{}, err
	}
	err = capnp.Struct(s).SetPtr(1, capnp.Struct(ss).ToPtr())
	return ss, err
}

func (s Alert) Status() AlertStatus {
	return AlertStatus(capnp.Struct(s).Uint16(0))
}

func (s Alert) SetStatus(v AlertStatus) {
	capnp.Struct(s).SetUint16(0, uint16(v))
}

func (s Alert) StartTime() int64 {
	return int64(capnp.Struct(s).Uint64(8))
}

func (s Alert) SetStartTime(v int64) {
	capnp.Struct(s).SetUint64(8, uint64(v))
}

func (s Alert) EndTime() int64 {
	return int64(capnp.Struct(s).Uint64(16))
}

func (s Alert) SetEndTime(v int64) {
	capnp.Struct(s).SetUint64(16, uint64(v))
}

// Alert_List is a list of Alert.
type Alert_List = capnp.StructList[Alert]

// NewAlert creates a new list of Alert.
func NewAlert_List(s *capnp.Segment, sz int32) (Alert_List, error) {
	l, err := capnp.NewCompositeList(s, capnp.ObjectSize{DataSize: 24, PointerCount: 2}, sz)
	return capnp.StructList[Alert](l), err
}

// Alert_Future is a wrapper for a Alert promised by a client call.
type Alert_Future struct{ *capnp.Future }

func (f Alert_Future) Struct() (Alert, error) {
	p, err := f.Future.Ptr()
	return Alert(p.Struct()), err
}
func (p Alert_Future) Labels() Map_Future {
	return Map_Future{Future: p.Future.Field(0, nil)}
}
func (p Alert_Future) Annotations() Map_Future {
	return Map_Future{Future: p.Future.Field(1, nil)}
}

type AlertStatus uint16

// AlertStatus_TypeID is the unique identifier for the type AlertStatus.
const AlertStatus_TypeID = 0x8b4db636305301a6

// Values of AlertStatus.
const (
	AlertStatus_firing   AlertStatus = 0
	AlertStatus_silenced AlertStatus = 1
)

// String returns the enum's constant name.
func (c AlertStatus) String() string {
	switch c {
	case AlertStatus_firing:
		return "firing"
	case AlertStatus_silenced:
		return "silenced"

	default:
		return ""
	}
}

// AlertStatusFromString returns the enum value with a name,
// or the zero value if there's no such value.
func AlertStatusFromString(c string) AlertStatus {
	switch c {
	case "firing":
		return AlertStatus_firing
	case "silenced":
		return AlertStatus_silenced

	default:
		return 0
	}
}

type AlertStatus_List = capnp.EnumList[AlertStatus]

func NewAlertStatus_List(s *capnp.Segment, sz int32) (AlertStatus_List, error) {
	return capnp.NewEnumList[AlertStatus](s, sz)
}

type Map capnp.Struct

// Map_TypeID is the unique identifier for the type Map.
const Map_TypeID = 0xd15e0ab5486a2241

func NewMap(s *capnp.Segment) (Map, error) {
	st, err := capnp.NewStruct(s, capnp.ObjectSize{DataSize: 0, PointerCount: 1})
	return Map(st), err
}

func NewRootMap(s *capnp.Segment) (Map, error) {
	st, err := capnp.NewRootStruct(s, capnp.ObjectSize{DataSize: 0, PointerCount: 1})
	return Map(st), err
}

func ReadRootMap(msg *capnp.Message) (Map, error) {
	root, err := msg.Root()
	return Map(root.Struct()), err
}

func (s Map) String() string {
	str, _ := text.Marshal(0xd15e0ab5486a2241, capnp.Struct(s))
	return str
}

func (s Map) EncodeAsPtr(seg *capnp.Segment) capnp.Ptr {
	return capnp.Struct(s).EncodeAsPtr(seg)
}

func (Map) DecodeFromPtr(p capnp.Ptr) Map {
	return Map(capnp.Struct{}.DecodeFromPtr(p))
}

func (s Map) ToPtr() capnp.Ptr {
	return capnp.Struct(s).ToPtr()
}
func (s Map) IsValid() bool {
	return capnp.Struct(s).IsValid()
}

func (s Map) Message() *capnp.Message {
	return capnp.Struct(s).Message()
}

func (s Map) Segment() *capnp.Segment {
	return capnp.Struct(s).Segment()
}
func (s Map) Entries() (Map_Entry_List, error) {
	p, err := capnp.Struct(s).Ptr(0)
	return Map_Entry_List(p.List()), err
}

func (s Map) HasEntries() bool {
	return capnp.Struct(s).HasPtr(0)
}

func (s Map) SetEntries(v Map_Entry_List) error {
	return capnp.Struct(s).SetPtr(0, v.ToPtr())
}

// NewEntries sets the entries field to a newly
// allocated Map_Entry_List, preferring placement in s's segment.
func (s Map) NewEntries(n int32) (Map_Entry_List, error) {
	l, err := NewMap_Entry_List(capnp.Struct(s).Segment(), n)
	if err != nil {
		return Map_Entry_List{}, err
	}
	err = capnp.Struct(s).SetPtr(0, l.ToPtr())
	return l, err
}

// Map_List is a list of Map.
type Map_List = capnp.StructList[Map]

// NewMap creates a new list of Map.
func NewMap_List(s *capnp.Segment, sz int32) (Map_List, error) {
	l, err := capnp.NewCompositeList(s, capnp.ObjectSize{DataSize: 0, PointerCount: 1}, sz)
	return capnp.StructList[Map](l), err
}

// Map_Future is a wrapper for a Map promised by a client call.
type Map_Future struct{ *capnp.Future }

func (f Map_Future) Struct() (Map, error) {
	p, err := f.Future.Ptr()
	return Map(p.Struct()), err
}

type Map_Entry capnp.Struct

// Map_Entry_TypeID is the unique identifier for the type Map_Entry.
const Map_Entry_TypeID = 0xe450011b48235af4

func NewMap_Entry(s *capnp.Segment) (Map_Entry, error) {
	st, err := capnp.NewStruct(s, capnp.ObjectSize{DataSize: 0, PointerCount: 2})
	return Map_Entry(st), err
}

func NewRootMap_Entry(s *capnp.Segment) (Map_Entry, error) {
	st, err := capnp.NewRootStruct(s, capnp.ObjectSize{DataSize: 0, PointerCount: 2})
	return Map_Entry(st), err
}

func ReadRootMap_Entry(msg *capnp.Message) (Map_Entry, error) {
	root, err := msg.Root()
	return Map_Entry(root.Struct()), err
}

func (s Map_Entry) String() string {
	str, _ := text.Marshal(0xe450011b48235af4, capnp.Struct(s))
	return str
}

func (s Map_Entry) EncodeAsPtr(seg *capnp.Segment) capnp.Ptr {
	return capnp.Struct(s).EncodeAsPtr(seg)
}

func (Map_Entry) DecodeFromPtr(p capnp.Ptr) Map_Entry {
	return Map_Entry(capnp.Struct{}.DecodeFromPtr(p))
}

func (s Map_Entry) ToPtr() capnp.Ptr {
	return capnp.Struct(s).ToPtr()
}
func (s Map_Entry) IsValid() bool {
	return capnp.Struct(s).IsValid()
}

func (s Map_Entry) Message() *capnp.Message {
	return capnp.Struct(s).Message()
}

func (s Map_Entry) Segment() *capnp.Segment {
	return capnp.Struct(s).Segment()
}
func (s Map_Entry) Key() (capnp.Ptr, error) {
	return capnp.Struct(s).Ptr(0)
}

func (s Map_Entry) HasKey() bool {
	return capnp.Struct(s).HasPtr(0)
}

func (s Map_Entry) SetKey(v capnp.Ptr) error {
	return capnp.Struct(s).SetPtr(0, v)
}
func (s Map_Entry) Value() (capnp.Ptr, error) {
	return capnp.Struct(s).Ptr(1)
}

func (s Map_Entry) HasValue() bool {
	return capnp.Struct(s).HasPtr(1)
}

func (s Map_Entry) SetValue(v capnp.Ptr) error {
	return capnp.Struct(s).SetPtr(1, v)
}

// Map_Entry_List is a list of Map_Entry.
type Map_Entry_List = capnp.StructList[Map_Entry]

// NewMap_Entry creates a new list of Map_Entry.
func NewMap_Entry_List(s *capnp.Segment, sz int32) (Map_Entry_List, error) {
	l, err := capnp.NewCompositeList(s, capnp.ObjectSize{DataSize: 0, PointerCount: 2}, sz)
	return capnp.StructList[Map_Entry](l), err
}

// Map_Entry_Future is a wrapper for a Map_Entry promised by a client call.
type Map_Entry_Future struct{ *capnp.Future }

func (f Map_Entry_Future) Struct() (Map_Entry, error) {
	p, err := f.Future.Ptr()
	return Map_Entry(p.Struct()), err
}
func (p Map_Entry_Future) Key() *capnp.Future {
	return p.Future.Field(0, nil)
}
func (p Map_Entry_Future) Value() *capnp.Future {
	return p.Future.Field(1, nil)
}

const schema_bca48e54c9238692 = "x\xda\xacRAk\x14M\x14\xac\xea\x9e\xf9&\x9fl" +
	"\xcctz$\x88\x91\x0d\x1b\xc5\x08\x1ab\x14\x95\xbdl" +
	"\x12\x15\x82\x12H;\x8b\xa0\xa00IF\x19\xddL\x96" +
	"\x9d\x89q\x0f\"x\xf0\"x\xc8?\x10T\xf0\xe0A" +
	"!\x9e\xf4\xa0\x100\xe0\xc5\xff\xe0\x1f\x10\xbc\x8f\xf4\x8a" +
	"\xbb\x1b\xe2\xd1\xc3\x0ctu\xbd\xaa~\xf5\xdeT\x8d3" +
	"\xce\xa9\xc1\x0d\x09a&\xdc\xff\x8aW\x0c\xa7\xce\xbe_" +
	"x\x0aU\x12\xc5\xe6\x93\xf1\x9d\xfa\xb3\x17\x1f\x00\xea\xe7" +
	"\xfc\xac_\xd3\x03\xf4K\x9e\x03\x8b\xca`\xe9\xc7\xfc\xc5" +
	"\x03\xdb0%\xca\x1e\xd3\x15\x96\xf2\x86\x9bz\xcb\x92O" +
	"\xbfe\x99`1[\xb9;\xbf\xb5\xef\xd67\xa8\x12\xfb" +
	"\xc8\xf4|\xea\xeb\xf2\xb1\xbe)G\x00\x1d\xcb\x1a\xa0?" +
	"\xc9\x91\xe2\xe7\x8d\xf1\xf9C\\\xfc\x0e\xe5\xf7\x15\xbb\xc2" +
	"\xf2?\xcawz[z\x1d\xe6\x06.\x14Q#n\xe5" +
	"\x93\xcb\x91h\xa6\xcd\xea\xac=\x84y\x94\xafgX$" +
	"\xcd\x00\x05\xa0T\x15 \xd5\xff\x97\x81\xda\xed\xa4\x95\xa4" +
	"w\x8a,i\xc4\xe9r\xbc\x02\xa0\xab\xc0?\x0a\xccm" +
	"m \x1d\xc0!\xa0\x1eV\x01\xf3@\xd2|\x15Td" +
	"@\x0b\xee,\x01\xe6\x8bd8JA\x8a\xc0:\xe9\x83" +
	"\xac\x02a@\xc9p\x8c\x82J2\xa0\x04\xf4a^\x05" +
	"\xc2Q\x8bOX\xdc\x11\x01\x1d@\x1f\xe5\x1c\x10\x8eY" +
	"\xfc\x04\x05k\x8dh)nd\xf4{M\x033T," +
	"\x1bG\xb0\x1fT<f\x06H\xdb\xa2\xb4\x7f_\x92%" +
	"\x88\xce\xe7\x83E\x94\xa6ky\x94'\xf0\xd6\xd2\x7f\xa0" +
	"W\xcb:\x99r\xa8\xb7! \x87\xc0\"\xcb\xa3V^" +
	"OV\xc1\x98.\x04]\xf0Q\x9c\xae\xd4\x93\xd5\xeey" +
	"w\xc0\x0bQ\x13\xc6!\xfb\x86\xcc\xe9\xf2\xa54o\xb5" +
	"\x8d\xd3\x8d|p\x0e\xb0O1\xe7\x85\xd5\xcb[I\x9c" +
	"q?\xb8(I\xbfW\xfa\xd7n\xecx,wv\x80" +
	"\xca\xad(w\xda\xbb\x12\xb7\xcb\xd7\xa2\xc6z\xbc\xe7)" +
	"\x93\xbf\x9d;\xbb\xd25?^\x01\xcc\x11I3\xd57" +
	"\xef\x93\xd3\x80\x99\x904g\x04\xbd{q\x9b\xc3\xdc\x95" +
	" \x87\xc1\xf2}k\xc2aw\xef\xd5\xaf\x00\x00\x00\xff" +
	"\xff\xda\xb9\xbd\xf6"

func init() {
	schemas.Register(schema_bca48e54c9238692,
		0x8b4db636305301a6,
		0xc5154448f10c0d22,
		0xd15e0ab5486a2241,
		0xe450011b48235af4)
}
