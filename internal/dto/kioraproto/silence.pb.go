// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v3.19.6
// source: silence.proto

package kioraproto

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type PostSilencesRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Silences []*Silence `protobuf:"bytes,1,rep,name=silences,proto3" json:"silences,omitempty"`
}

func (x *PostSilencesRequest) Reset() {
	*x = PostSilencesRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_silence_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PostSilencesRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PostSilencesRequest) ProtoMessage() {}

func (x *PostSilencesRequest) ProtoReflect() protoreflect.Message {
	mi := &file_silence_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PostSilencesRequest.ProtoReflect.Descriptor instead.
func (*PostSilencesRequest) Descriptor() ([]byte, []int) {
	return file_silence_proto_rawDescGZIP(), []int{0}
}

func (x *PostSilencesRequest) GetSilences() []*Silence {
	if x != nil {
		return x.Silences
	}
	return nil
}

type Silence struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ID        string                 `protobuf:"bytes,1,opt,name=ID,proto3" json:"ID,omitempty"`
	Creator   string                 `protobuf:"bytes,2,opt,name=Creator,proto3" json:"Creator,omitempty"`
	Comment   string                 `protobuf:"bytes,3,opt,name=Comment,proto3" json:"Comment,omitempty"`
	StartTime *timestamppb.Timestamp `protobuf:"bytes,4,opt,name=startTime,proto3" json:"startTime,omitempty"`
	EndTime   *timestamppb.Timestamp `protobuf:"bytes,5,opt,name=endTime,proto3" json:"endTime,omitempty"`
	Matchers  []*Matcher             `protobuf:"bytes,6,rep,name=matchers,proto3" json:"matchers,omitempty"`
}

func (x *Silence) Reset() {
	*x = Silence{}
	if protoimpl.UnsafeEnabled {
		mi := &file_silence_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Silence) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Silence) ProtoMessage() {}

func (x *Silence) ProtoReflect() protoreflect.Message {
	mi := &file_silence_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Silence.ProtoReflect.Descriptor instead.
func (*Silence) Descriptor() ([]byte, []int) {
	return file_silence_proto_rawDescGZIP(), []int{1}
}

func (x *Silence) GetID() string {
	if x != nil {
		return x.ID
	}
	return ""
}

func (x *Silence) GetCreator() string {
	if x != nil {
		return x.Creator
	}
	return ""
}

func (x *Silence) GetComment() string {
	if x != nil {
		return x.Comment
	}
	return ""
}

func (x *Silence) GetStartTime() *timestamppb.Timestamp {
	if x != nil {
		return x.StartTime
	}
	return nil
}

func (x *Silence) GetEndTime() *timestamppb.Timestamp {
	if x != nil {
		return x.EndTime
	}
	return nil
}

func (x *Silence) GetMatchers() []*Matcher {
	if x != nil {
		return x.Matchers
	}
	return nil
}

type Matcher struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Key      string `protobuf:"bytes,1,opt,name=Key,proto3" json:"Key,omitempty"`
	Value    string `protobuf:"bytes,2,opt,name=Value,proto3" json:"Value,omitempty"`
	Regex    bool   `protobuf:"varint,3,opt,name=Regex,proto3" json:"Regex,omitempty"`
	Negative bool   `protobuf:"varint,4,opt,name=Negative,proto3" json:"Negative,omitempty"`
}

func (x *Matcher) Reset() {
	*x = Matcher{}
	if protoimpl.UnsafeEnabled {
		mi := &file_silence_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Matcher) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Matcher) ProtoMessage() {}

func (x *Matcher) ProtoReflect() protoreflect.Message {
	mi := &file_silence_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Matcher.ProtoReflect.Descriptor instead.
func (*Matcher) Descriptor() ([]byte, []int) {
	return file_silence_proto_rawDescGZIP(), []int{2}
}

func (x *Matcher) GetKey() string {
	if x != nil {
		return x.Key
	}
	return ""
}

func (x *Matcher) GetValue() string {
	if x != nil {
		return x.Value
	}
	return ""
}

func (x *Matcher) GetRegex() bool {
	if x != nil {
		return x.Regex
	}
	return false
}

func (x *Matcher) GetNegative() bool {
	if x != nil {
		return x.Negative
	}
	return false
}

var File_silence_proto protoreflect.FileDescriptor

var file_silence_proto_rawDesc = []byte{
	0x0a, 0x0d, 0x73, 0x69, 0x6c, 0x65, 0x6e, 0x63, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12,
	0x0a, 0x6b, 0x69, 0x6f, 0x72, 0x61, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1f, 0x67, 0x6f, 0x6f,
	0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x74, 0x69, 0x6d,
	0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x46, 0x0a, 0x13,
	0x50, 0x6f, 0x73, 0x74, 0x53, 0x69, 0x6c, 0x65, 0x6e, 0x63, 0x65, 0x73, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x12, 0x2f, 0x0a, 0x08, 0x73, 0x69, 0x6c, 0x65, 0x6e, 0x63, 0x65, 0x73, 0x18,
	0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x13, 0x2e, 0x6b, 0x69, 0x6f, 0x72, 0x61, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x2e, 0x53, 0x69, 0x6c, 0x65, 0x6e, 0x63, 0x65, 0x52, 0x08, 0x73, 0x69, 0x6c, 0x65,
	0x6e, 0x63, 0x65, 0x73, 0x22, 0xee, 0x01, 0x0a, 0x07, 0x53, 0x69, 0x6c, 0x65, 0x6e, 0x63, 0x65,
	0x12, 0x0e, 0x0a, 0x02, 0x49, 0x44, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x49, 0x44,
	0x12, 0x18, 0x0a, 0x07, 0x43, 0x72, 0x65, 0x61, 0x74, 0x6f, 0x72, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x07, 0x43, 0x72, 0x65, 0x61, 0x74, 0x6f, 0x72, 0x12, 0x18, 0x0a, 0x07, 0x43, 0x6f,
	0x6d, 0x6d, 0x65, 0x6e, 0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x43, 0x6f, 0x6d,
	0x6d, 0x65, 0x6e, 0x74, 0x12, 0x38, 0x0a, 0x09, 0x73, 0x74, 0x61, 0x72, 0x74, 0x54, 0x69, 0x6d,
	0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74,
	0x61, 0x6d, 0x70, 0x52, 0x09, 0x73, 0x74, 0x61, 0x72, 0x74, 0x54, 0x69, 0x6d, 0x65, 0x12, 0x34,
	0x0a, 0x07, 0x65, 0x6e, 0x64, 0x54, 0x69, 0x6d, 0x65, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75,
	0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x07, 0x65, 0x6e, 0x64,
	0x54, 0x69, 0x6d, 0x65, 0x12, 0x2f, 0x0a, 0x08, 0x6d, 0x61, 0x74, 0x63, 0x68, 0x65, 0x72, 0x73,
	0x18, 0x06, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x13, 0x2e, 0x6b, 0x69, 0x6f, 0x72, 0x61, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x2e, 0x4d, 0x61, 0x74, 0x63, 0x68, 0x65, 0x72, 0x52, 0x08, 0x6d, 0x61, 0x74,
	0x63, 0x68, 0x65, 0x72, 0x73, 0x22, 0x63, 0x0a, 0x07, 0x4d, 0x61, 0x74, 0x63, 0x68, 0x65, 0x72,
	0x12, 0x10, 0x0a, 0x03, 0x4b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x4b,
	0x65, 0x79, 0x12, 0x14, 0x0a, 0x05, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x05, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x52, 0x65, 0x67, 0x65,
	0x78, 0x18, 0x03, 0x20, 0x01, 0x28, 0x08, 0x52, 0x05, 0x52, 0x65, 0x67, 0x65, 0x78, 0x12, 0x1a,
	0x0a, 0x08, 0x4e, 0x65, 0x67, 0x61, 0x74, 0x69, 0x76, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x08,
	0x52, 0x08, 0x4e, 0x65, 0x67, 0x61, 0x74, 0x69, 0x76, 0x65, 0x42, 0x31, 0x5a, 0x2f, 0x67, 0x69,
	0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x73, 0x69, 0x6e, 0x6b, 0x69, 0x6e, 0x67,
	0x70, 0x6f, 0x69, 0x6e, 0x74, 0x2f, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x2f, 0x64,
	0x74, 0x6f, 0x2f, 0x6b, 0x69, 0x6f, 0x72, 0x61, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x06, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_silence_proto_rawDescOnce sync.Once
	file_silence_proto_rawDescData = file_silence_proto_rawDesc
)

func file_silence_proto_rawDescGZIP() []byte {
	file_silence_proto_rawDescOnce.Do(func() {
		file_silence_proto_rawDescData = protoimpl.X.CompressGZIP(file_silence_proto_rawDescData)
	})
	return file_silence_proto_rawDescData
}

var file_silence_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_silence_proto_goTypes = []interface{}{
	(*PostSilencesRequest)(nil),   // 0: kioraproto.PostSilencesRequest
	(*Silence)(nil),               // 1: kioraproto.Silence
	(*Matcher)(nil),               // 2: kioraproto.Matcher
	(*timestamppb.Timestamp)(nil), // 3: google.protobuf.Timestamp
}
var file_silence_proto_depIdxs = []int32{
	1, // 0: kioraproto.PostSilencesRequest.silences:type_name -> kioraproto.Silence
	3, // 1: kioraproto.Silence.startTime:type_name -> google.protobuf.Timestamp
	3, // 2: kioraproto.Silence.endTime:type_name -> google.protobuf.Timestamp
	2, // 3: kioraproto.Silence.matchers:type_name -> kioraproto.Matcher
	4, // [4:4] is the sub-list for method output_type
	4, // [4:4] is the sub-list for method input_type
	4, // [4:4] is the sub-list for extension type_name
	4, // [4:4] is the sub-list for extension extendee
	0, // [0:4] is the sub-list for field type_name
}

func init() { file_silence_proto_init() }
func file_silence_proto_init() {
	if File_silence_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_silence_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PostSilencesRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_silence_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Silence); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_silence_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Matcher); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_silence_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_silence_proto_goTypes,
		DependencyIndexes: file_silence_proto_depIdxs,
		MessageInfos:      file_silence_proto_msgTypes,
	}.Build()
	File_silence_proto = out.File
	file_silence_proto_rawDesc = nil
	file_silence_proto_goTypes = nil
	file_silence_proto_depIdxs = nil
}
