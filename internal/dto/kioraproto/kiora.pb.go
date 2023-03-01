// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v3.19.6
// source: kiora.proto

package kioraproto

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type KioraLogMessage struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Types that are assignable to Log:
	//	*KioraLogMessage_Alerts
	//	*KioraLogMessage_Silences
	Log isKioraLogMessage_Log `protobuf_oneof:"log"`
}

func (x *KioraLogMessage) Reset() {
	*x = KioraLogMessage{}
	if protoimpl.UnsafeEnabled {
		mi := &file_kiora_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *KioraLogMessage) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*KioraLogMessage) ProtoMessage() {}

func (x *KioraLogMessage) ProtoReflect() protoreflect.Message {
	mi := &file_kiora_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use KioraLogMessage.ProtoReflect.Descriptor instead.
func (*KioraLogMessage) Descriptor() ([]byte, []int) {
	return file_kiora_proto_rawDescGZIP(), []int{0}
}

func (m *KioraLogMessage) GetLog() isKioraLogMessage_Log {
	if m != nil {
		return m.Log
	}
	return nil
}

func (x *KioraLogMessage) GetAlerts() *PostAlertsMessage {
	if x, ok := x.GetLog().(*KioraLogMessage_Alerts); ok {
		return x.Alerts
	}
	return nil
}

func (x *KioraLogMessage) GetSilences() *PostSilencesRequest {
	if x, ok := x.GetLog().(*KioraLogMessage_Silences); ok {
		return x.Silences
	}
	return nil
}

type isKioraLogMessage_Log interface {
	isKioraLogMessage_Log()
}

type KioraLogMessage_Alerts struct {
	Alerts *PostAlertsMessage `protobuf:"bytes,1,opt,name=alerts,proto3,oneof"`
}

type KioraLogMessage_Silences struct {
	Silences *PostSilencesRequest `protobuf:"bytes,2,opt,name=silences,proto3,oneof"`
}

func (*KioraLogMessage_Alerts) isKioraLogMessage_Log() {}

func (*KioraLogMessage_Silences) isKioraLogMessage_Log() {}

type KioraLogReply struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *KioraLogReply) Reset() {
	*x = KioraLogReply{}
	if protoimpl.UnsafeEnabled {
		mi := &file_kiora_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *KioraLogReply) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*KioraLogReply) ProtoMessage() {}

func (x *KioraLogReply) ProtoReflect() protoreflect.Message {
	mi := &file_kiora_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use KioraLogReply.ProtoReflect.Descriptor instead.
func (*KioraLogReply) Descriptor() ([]byte, []int) {
	return file_kiora_proto_rawDescGZIP(), []int{1}
}

var File_kiora_proto protoreflect.FileDescriptor

var file_kiora_proto_rawDesc = []byte{
	0x0a, 0x0b, 0x6b, 0x69, 0x6f, 0x72, 0x61, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0a, 0x6b,
	0x69, 0x6f, 0x72, 0x61, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x0b, 0x61, 0x6c, 0x65, 0x72, 0x74,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x0d, 0x73, 0x69, 0x6c, 0x65, 0x6e, 0x63, 0x65, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x90, 0x01, 0x0a, 0x0f, 0x4b, 0x69, 0x6f, 0x72, 0x61, 0x4c,
	0x6f, 0x67, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x37, 0x0a, 0x06, 0x61, 0x6c, 0x65,
	0x72, 0x74, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1d, 0x2e, 0x6b, 0x69, 0x6f, 0x72,
	0x61, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x50, 0x6f, 0x73, 0x74, 0x41, 0x6c, 0x65, 0x72, 0x74,
	0x73, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x48, 0x00, 0x52, 0x06, 0x61, 0x6c, 0x65, 0x72,
	0x74, 0x73, 0x12, 0x3d, 0x0a, 0x08, 0x73, 0x69, 0x6c, 0x65, 0x6e, 0x63, 0x65, 0x73, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x1f, 0x2e, 0x6b, 0x69, 0x6f, 0x72, 0x61, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x2e, 0x50, 0x6f, 0x73, 0x74, 0x53, 0x69, 0x6c, 0x65, 0x6e, 0x63, 0x65, 0x73, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x48, 0x00, 0x52, 0x08, 0x73, 0x69, 0x6c, 0x65, 0x6e, 0x63, 0x65,
	0x73, 0x42, 0x05, 0x0a, 0x03, 0x6c, 0x6f, 0x67, 0x22, 0x0f, 0x0a, 0x0d, 0x4b, 0x69, 0x6f, 0x72,
	0x61, 0x4c, 0x6f, 0x67, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x32, 0x4d, 0x0a, 0x05, 0x4b, 0x69, 0x6f,
	0x72, 0x61, 0x12, 0x44, 0x0a, 0x08, 0x41, 0x70, 0x70, 0x6c, 0x79, 0x4c, 0x6f, 0x67, 0x12, 0x1b,
	0x2e, 0x6b, 0x69, 0x6f, 0x72, 0x61, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x4b, 0x69, 0x6f, 0x72,
	0x61, 0x4c, 0x6f, 0x67, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x1a, 0x19, 0x2e, 0x6b, 0x69,
	0x6f, 0x72, 0x61, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x4b, 0x69, 0x6f, 0x72, 0x61, 0x4c, 0x6f,
	0x67, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x22, 0x00, 0x42, 0x31, 0x5a, 0x2f, 0x67, 0x69, 0x74, 0x68,
	0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x73, 0x69, 0x6e, 0x6b, 0x69, 0x6e, 0x67, 0x70, 0x6f,
	0x69, 0x6e, 0x74, 0x2f, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x2f, 0x64, 0x74, 0x6f,
	0x2f, 0x6b, 0x69, 0x6f, 0x72, 0x61, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x06, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x33,
}

var (
	file_kiora_proto_rawDescOnce sync.Once
	file_kiora_proto_rawDescData = file_kiora_proto_rawDesc
)

func file_kiora_proto_rawDescGZIP() []byte {
	file_kiora_proto_rawDescOnce.Do(func() {
		file_kiora_proto_rawDescData = protoimpl.X.CompressGZIP(file_kiora_proto_rawDescData)
	})
	return file_kiora_proto_rawDescData
}

var file_kiora_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_kiora_proto_goTypes = []interface{}{
	(*KioraLogMessage)(nil),     // 0: kioraproto.KioraLogMessage
	(*KioraLogReply)(nil),       // 1: kioraproto.KioraLogReply
	(*PostAlertsMessage)(nil),   // 2: kioraproto.PostAlertsMessage
	(*PostSilencesRequest)(nil), // 3: kioraproto.PostSilencesRequest
}
var file_kiora_proto_depIdxs = []int32{
	2, // 0: kioraproto.KioraLogMessage.alerts:type_name -> kioraproto.PostAlertsMessage
	3, // 1: kioraproto.KioraLogMessage.silences:type_name -> kioraproto.PostSilencesRequest
	0, // 2: kioraproto.Kiora.ApplyLog:input_type -> kioraproto.KioraLogMessage
	1, // 3: kioraproto.Kiora.ApplyLog:output_type -> kioraproto.KioraLogReply
	3, // [3:4] is the sub-list for method output_type
	2, // [2:3] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_kiora_proto_init() }
func file_kiora_proto_init() {
	if File_kiora_proto != nil {
		return
	}
	file_alert_proto_init()
	file_silence_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_kiora_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*KioraLogMessage); i {
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
		file_kiora_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*KioraLogReply); i {
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
	file_kiora_proto_msgTypes[0].OneofWrappers = []interface{}{
		(*KioraLogMessage_Alerts)(nil),
		(*KioraLogMessage_Silences)(nil),
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_kiora_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_kiora_proto_goTypes,
		DependencyIndexes: file_kiora_proto_depIdxs,
		MessageInfos:      file_kiora_proto_msgTypes,
	}.Build()
	File_kiora_proto = out.File
	file_kiora_proto_rawDesc = nil
	file_kiora_proto_goTypes = nil
	file_kiora_proto_depIdxs = nil
}
