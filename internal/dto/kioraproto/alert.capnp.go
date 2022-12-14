// Code generated by capnpc-go. DO NOT EDIT.

package kioraproto

import (
	capnp "capnproto.org/go/capnp/v3"
	text "capnproto.org/go/capnp/v3/encoding/text"
	schemas "capnproto.org/go/capnp/v3/schemas"
)

type Alerts capnp.Struct

// Alerts_TypeID is the unique identifier for the type Alerts.
const Alerts_TypeID = 0xf3c77a394ac60718

func NewAlerts(s *capnp.Segment) (Alerts, error) {
	st, err := capnp.NewStruct(s, capnp.ObjectSize{DataSize: 0, PointerCount: 1})
	return Alerts(st), err
}

func NewRootAlerts(s *capnp.Segment) (Alerts, error) {
	st, err := capnp.NewRootStruct(s, capnp.ObjectSize{DataSize: 0, PointerCount: 1})
	return Alerts(st), err
}

func ReadRootAlerts(msg *capnp.Message) (Alerts, error) {
	root, err := msg.Root()
	return Alerts(root.Struct()), err
}

func (s Alerts) String() string {
	str, _ := text.Marshal(0xf3c77a394ac60718, capnp.Struct(s))
	return str
}

func (s Alerts) EncodeAsPtr(seg *capnp.Segment) capnp.Ptr {
	return capnp.Struct(s).EncodeAsPtr(seg)
}

func (Alerts) DecodeFromPtr(p capnp.Ptr) Alerts {
	return Alerts(capnp.Struct{}.DecodeFromPtr(p))
}

func (s Alerts) ToPtr() capnp.Ptr {
	return capnp.Struct(s).ToPtr()
}
func (s Alerts) IsValid() bool {
	return capnp.Struct(s).IsValid()
}

func (s Alerts) Message() *capnp.Message {
	return capnp.Struct(s).Message()
}

func (s Alerts) Segment() *capnp.Segment {
	return capnp.Struct(s).Segment()
}
func (s Alerts) Alerts() (Alert_List, error) {
	p, err := capnp.Struct(s).Ptr(0)
	return Alert_List(p.List()), err
}

func (s Alerts) HasAlerts() bool {
	return capnp.Struct(s).HasPtr(0)
}

func (s Alerts) SetAlerts(v Alert_List) error {
	return capnp.Struct(s).SetPtr(0, v.ToPtr())
}

// NewAlerts sets the alerts field to a newly
// allocated Alert_List, preferring placement in s's segment.
func (s Alerts) NewAlerts(n int32) (Alert_List, error) {
	l, err := NewAlert_List(capnp.Struct(s).Segment(), n)
	if err != nil {
		return Alert_List{}, err
	}
	err = capnp.Struct(s).SetPtr(0, l.ToPtr())
	return l, err
}

// Alerts_List is a list of Alerts.
type Alerts_List = capnp.StructList[Alerts]

// NewAlerts creates a new list of Alerts.
func NewAlerts_List(s *capnp.Segment, sz int32) (Alerts_List, error) {
	l, err := capnp.NewCompositeList(s, capnp.ObjectSize{DataSize: 0, PointerCount: 1}, sz)
	return capnp.StructList[Alerts](l), err
}

// Alerts_Future is a wrapper for a Alerts promised by a client call.
type Alerts_Future struct{ *capnp.Future }

func (f Alerts_Future) Struct() (Alerts, error) {
	p, err := f.Future.Ptr()
	return Alerts(p.Struct()), err
}

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
	AlertStatus_resolved AlertStatus = 1
)

// String returns the enum's constant name.
func (c AlertStatus) String() string {
	switch c {
	case AlertStatus_firing:
		return "firing"
	case AlertStatus_resolved:
		return "resolved"

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
	case "resolved":
		return AlertStatus_resolved

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

const schema_bca48e54c9238692 = "x\xda\xacSOk\x13A\x1c}ov\xd6\xad\x92\xda" +
	"l'R\xc5JJZ\xb1\x82\xd6Z\xff6\x08i\x8a" +
	"B\xa9\x14:n\xf1\xe0A\xd8\xdaU\xa2\xe9\xa6d\xb7" +
	"\xd5\x08E\x14\x14D,\xd8\x9b\x82\x07A\xfc\x00\x82\x9e" +
	"<\x09\xfe)x\xf1;\xf8\x05D/\x1ede\xa2&" +
	"[\xe2\xd1\xc3.3o\xdf\xef\xf7\xe6\xbd\xf9\xed\xe8\x1a" +
	"'\xe4\xe1\xeeC6\x84>eoI^\xd0\x1b=\xfe" +
	"z\xe6\x01\xdc\x8cH\xd6\xef\x0dn\xcc\xad=\x7f\x03P" +
	"\xad\x8a\xb7\xea\xaep\x00u[\x9c\x00\x93Bw\xe6\xeb" +
	"\xd4\xe9\x1d\xef\xa03\xb4\xdaL\xbbI\xb9/\xd6\xd5#" +
	"\xb3:\xf2P\xe4\x09&\xe5\xc2\xd5\xa9W\xdb.~\x86" +
	"\x9ba\x8aL'KuL\xdeQ\xe3\xb2\x0fPeY" +
	"\x02\xd4S\xd9\x97|\xbf08\xb5\x9b\xb3_\xe0fS" +
	"\xc5\xb60\xfc\xc7\xf2\xa5z&\x9d&\xf3:\x98\xect" +
	"\xdeO\x8f\xdf\xfc\xf0\xad\xa37\xa0~\xc8'\x8a\xb6Y" +
	"\xfd\x94%\xd4\x12\xbf\x1a\xd4\xe3\x91K\xbeX\x0a\x97\x8a" +
	"e\xb3\xf1b?^\x8e0K\xea.\x0a\xc0u\x8b\x00" +
	"\xe9n\x9d\x06J\x97+\xf5Jx%\xa9\x07Q\xad\xba" +
	"\x12,\x00hu\xe0\xdf\x0e\x8cMm\xce\x92\x80$\xe0" +
	"\xae\x16\x01}\xc3\xa2\xfe$\xe8\x929\x1apc\x1e\xd0" +
	"\x1f-z\xfd\x14\xa4\xc8\x19%\xb5\x8bE\xc0\xcb\xd1\xa2" +
	"7@A\xd7b\x8e\x16\xa0\xf6\xf0\x1c\xe0\xf5\x1b|\xd8" +
	"\xe0R\xe4(\x01\xb5\x97\x93\x807`\xf0\x03\x14,U" +
	"\xfd\xf9\xa0\x1a1\xdb\x0e\x08\x98\xa0\xcb\xbc\x96\x82i\xd0" +
	"\xe5>\xddE\x1a\x8b\x96yg-2\x03\xd1|\xb2`" +
	"\xe2\x87a-\xf6\xe3\x0a\x9cZ\xf8\x1f\xfa\x95\xa2f\xa6" +
	"\xeciO\x13\xc8\x1e0\x89b\xbf\x1e\xcfU\x16\xc1\x80" +
	"6\x04m\xf0V\x10.\xccU\x16[\xfb\xcd\x01\xcf\xf8" +
	"K\xd0\x92L\x0d\x04\xc7\xf2g\xc2\xb8\xde\xd0\xb2\x15y" +
	"\xf7$`\x8e\xa2O\x0a\xd3/\xaeW\x82\x88\xdb\xc1Y" +
	"\x8b\xcc\xb6K\xff\xe9\xc6\\\x8f\xe1\x96\xbb\xe8\xda\x05\xd7" +
	"\x1es\xce\x06\x8d\xfcy\xbf\xba\x1ct\x1ce\xe4\xb7r" +
	"sVZ\xe2\xfb\x0b\x80\x1e\xb2\xa8GS\xf7}p\x0c" +
	"\xd0\xc3\x16\xf5QA\xe7Z\xd0`/7%\xc8^0" +
	"\xbfbD\xd8kw~\xea\x9c2+\x8e\x8cl\xcas" +
	"\xf1\x8f\xe7!\xc1R\x93\x9e\xb2\xdc\xfa9A\x03\xfe\x0a" +
	"\x00\x00\xff\xff\xf92\xda\x7f"

func init() {
	schemas.Register(schema_bca48e54c9238692,
		0x8b4db636305301a6,
		0xc5154448f10c0d22,
		0xd15e0ab5486a2241,
		0xe450011b48235af4,
		0xf3c77a394ac60718)
}
