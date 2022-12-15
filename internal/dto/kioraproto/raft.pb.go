// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v3.19.4
// source: raft.proto

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

type RaftLogMessage struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Types that are assignable to Log:
	//	*RaftLogMessage_Alerts
	//	*RaftLogMessage_Silences
	Log isRaftLogMessage_Log `protobuf_oneof:"log"`
}

func (x *RaftLogMessage) Reset() {
	*x = RaftLogMessage{}
	if protoimpl.UnsafeEnabled {
		mi := &file_raft_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RaftLogMessage) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RaftLogMessage) ProtoMessage() {}

func (x *RaftLogMessage) ProtoReflect() protoreflect.Message {
	mi := &file_raft_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RaftLogMessage.ProtoReflect.Descriptor instead.
func (*RaftLogMessage) Descriptor() ([]byte, []int) {
	return file_raft_proto_rawDescGZIP(), []int{0}
}

func (m *RaftLogMessage) GetLog() isRaftLogMessage_Log {
	if m != nil {
		return m.Log
	}
	return nil
}

func (x *RaftLogMessage) GetAlerts() *PostAlertsMessage {
	if x, ok := x.GetLog().(*RaftLogMessage_Alerts); ok {
		return x.Alerts
	}
	return nil
}

func (x *RaftLogMessage) GetSilences() *PostSilencesRequest {
	if x, ok := x.GetLog().(*RaftLogMessage_Silences); ok {
		return x.Silences
	}
	return nil
}

type isRaftLogMessage_Log interface {
	isRaftLogMessage_Log()
}

type RaftLogMessage_Alerts struct {
	Alerts *PostAlertsMessage `protobuf:"bytes,1,opt,name=alerts,proto3,oneof"`
}

type RaftLogMessage_Silences struct {
	Silences *PostSilencesRequest `protobuf:"bytes,2,opt,name=silences,proto3,oneof"`
}

func (*RaftLogMessage_Alerts) isRaftLogMessage_Log() {}

func (*RaftLogMessage_Silences) isRaftLogMessage_Log() {}

var File_raft_proto protoreflect.FileDescriptor

var file_raft_proto_rawDesc = []byte{
	0x0a, 0x0a, 0x72, 0x61, 0x66, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0a, 0x6b, 0x69,
	0x6f, 0x72, 0x61, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x0b, 0x61, 0x6c, 0x65, 0x72, 0x74, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x0d, 0x73, 0x69, 0x6c, 0x65, 0x6e, 0x63, 0x65, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x22, 0x8f, 0x01, 0x0a, 0x0e, 0x52, 0x61, 0x66, 0x74, 0x4c, 0x6f, 0x67,
	0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x37, 0x0a, 0x06, 0x61, 0x6c, 0x65, 0x72, 0x74,
	0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1d, 0x2e, 0x6b, 0x69, 0x6f, 0x72, 0x61, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x50, 0x6f, 0x73, 0x74, 0x41, 0x6c, 0x65, 0x72, 0x74, 0x73, 0x4d,
	0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x48, 0x00, 0x52, 0x06, 0x61, 0x6c, 0x65, 0x72, 0x74, 0x73,
	0x12, 0x3d, 0x0a, 0x08, 0x73, 0x69, 0x6c, 0x65, 0x6e, 0x63, 0x65, 0x73, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x1f, 0x2e, 0x6b, 0x69, 0x6f, 0x72, 0x61, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e,
	0x50, 0x6f, 0x73, 0x74, 0x53, 0x69, 0x6c, 0x65, 0x6e, 0x63, 0x65, 0x73, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x48, 0x00, 0x52, 0x08, 0x73, 0x69, 0x6c, 0x65, 0x6e, 0x63, 0x65, 0x73, 0x42,
	0x05, 0x0a, 0x03, 0x6c, 0x6f, 0x67, 0x42, 0x31, 0x5a, 0x2f, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62,
	0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x73, 0x69, 0x6e, 0x6b, 0x69, 0x6e, 0x67, 0x70, 0x6f, 0x69, 0x6e,
	0x74, 0x2f, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x2f, 0x64, 0x74, 0x6f, 0x2f, 0x6b,
	0x69, 0x6f, 0x72, 0x61, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x33,
}

var (
	file_raft_proto_rawDescOnce sync.Once
	file_raft_proto_rawDescData = file_raft_proto_rawDesc
)

func file_raft_proto_rawDescGZIP() []byte {
	file_raft_proto_rawDescOnce.Do(func() {
		file_raft_proto_rawDescData = protoimpl.X.CompressGZIP(file_raft_proto_rawDescData)
	})
	return file_raft_proto_rawDescData
}

var file_raft_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_raft_proto_goTypes = []interface{}{
	(*RaftLogMessage)(nil),      // 0: kioraproto.RaftLogMessage
	(*PostAlertsMessage)(nil),   // 1: kioraproto.PostAlertsMessage
	(*PostSilencesRequest)(nil), // 2: kioraproto.PostSilencesRequest
}
var file_raft_proto_depIdxs = []int32{
	1, // 0: kioraproto.RaftLogMessage.alerts:type_name -> kioraproto.PostAlertsMessage
	2, // 1: kioraproto.RaftLogMessage.silences:type_name -> kioraproto.PostSilencesRequest
	2, // [2:2] is the sub-list for method output_type
	2, // [2:2] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_raft_proto_init() }
func file_raft_proto_init() {
	if File_raft_proto != nil {
		return
	}
	file_alert_proto_init()
	file_silence_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_raft_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RaftLogMessage); i {
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
	file_raft_proto_msgTypes[0].OneofWrappers = []interface{}{
		(*RaftLogMessage_Alerts)(nil),
		(*RaftLogMessage_Silences)(nil),
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_raft_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_raft_proto_goTypes,
		DependencyIndexes: file_raft_proto_depIdxs,
		MessageInfos:      file_raft_proto_msgTypes,
	}.Build()
	File_raft_proto = out.File
	file_raft_proto_rawDesc = nil
	file_raft_proto_goTypes = nil
	file_raft_proto_depIdxs = nil
}