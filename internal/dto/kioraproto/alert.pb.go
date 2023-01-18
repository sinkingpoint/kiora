// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v3.19.6
// source: alert.proto

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

// AlertStatus is the status of the alert coming in.
type AlertStatus int32

const (
	AlertStatus_firing     AlertStatus = 0
	AlertStatus_resolved   AlertStatus = 1
	AlertStatus_processing AlertStatus = 2
)

// Enum value maps for AlertStatus.
var (
	AlertStatus_name = map[int32]string{
		0: "firing",
		1: "resolved",
		2: "processing",
	}
	AlertStatus_value = map[string]int32{
		"firing":     0,
		"resolved":   1,
		"processing": 2,
	}
)

func (x AlertStatus) Enum() *AlertStatus {
	p := new(AlertStatus)
	*p = x
	return p
}

func (x AlertStatus) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (AlertStatus) Descriptor() protoreflect.EnumDescriptor {
	return file_alert_proto_enumTypes[0].Descriptor()
}

func (AlertStatus) Type() protoreflect.EnumType {
	return &file_alert_proto_enumTypes[0]
}

func (x AlertStatus) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use AlertStatus.Descriptor instead.
func (AlertStatus) EnumDescriptor() ([]byte, []int) {
	return file_alert_proto_rawDescGZIP(), []int{0}
}

type PostAlertsMessage struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Alerts []*Alert `protobuf:"bytes,1,rep,name=alerts,proto3" json:"alerts,omitempty"`
}

func (x *PostAlertsMessage) Reset() {
	*x = PostAlertsMessage{}
	if protoimpl.UnsafeEnabled {
		mi := &file_alert_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PostAlertsMessage) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PostAlertsMessage) ProtoMessage() {}

func (x *PostAlertsMessage) ProtoReflect() protoreflect.Message {
	mi := &file_alert_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PostAlertsMessage.ProtoReflect.Descriptor instead.
func (*PostAlertsMessage) Descriptor() ([]byte, []int) {
	return file_alert_proto_rawDescGZIP(), []int{0}
}

func (x *PostAlertsMessage) GetAlerts() []*Alert {
	if x != nil {
		return x.Alerts
	}
	return nil
}

type Alert struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Labels      map[string]string      `protobuf:"bytes,1,rep,name=labels,proto3" json:"labels,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	Annotations map[string]string      `protobuf:"bytes,2,rep,name=annotations,proto3" json:"annotations,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	Status      AlertStatus            `protobuf:"varint,3,opt,name=status,proto3,enum=kioraproto.AlertStatus" json:"status,omitempty"`
	StartTime   *timestamppb.Timestamp `protobuf:"bytes,4,opt,name=startTime,proto3" json:"startTime,omitempty"`
	EndTime     *timestamppb.Timestamp `protobuf:"bytes,5,opt,name=endTime,proto3" json:"endTime,omitempty"`
}

func (x *Alert) Reset() {
	*x = Alert{}
	if protoimpl.UnsafeEnabled {
		mi := &file_alert_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Alert) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Alert) ProtoMessage() {}

func (x *Alert) ProtoReflect() protoreflect.Message {
	mi := &file_alert_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Alert.ProtoReflect.Descriptor instead.
func (*Alert) Descriptor() ([]byte, []int) {
	return file_alert_proto_rawDescGZIP(), []int{1}
}

func (x *Alert) GetLabels() map[string]string {
	if x != nil {
		return x.Labels
	}
	return nil
}

func (x *Alert) GetAnnotations() map[string]string {
	if x != nil {
		return x.Annotations
	}
	return nil
}

func (x *Alert) GetStatus() AlertStatus {
	if x != nil {
		return x.Status
	}
	return AlertStatus_firing
}

func (x *Alert) GetStartTime() *timestamppb.Timestamp {
	if x != nil {
		return x.StartTime
	}
	return nil
}

func (x *Alert) GetEndTime() *timestamppb.Timestamp {
	if x != nil {
		return x.EndTime
	}
	return nil
}

var File_alert_proto protoreflect.FileDescriptor

var file_alert_proto_rawDesc = []byte{
	0x0a, 0x0b, 0x61, 0x6c, 0x65, 0x72, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0a, 0x6b,
	0x69, 0x6f, 0x72, 0x61, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1f, 0x67, 0x6f, 0x6f, 0x67, 0x6c,
	0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x74, 0x69, 0x6d, 0x65, 0x73,
	0x74, 0x61, 0x6d, 0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x3e, 0x0a, 0x11, 0x50, 0x6f,
	0x73, 0x74, 0x41, 0x6c, 0x65, 0x72, 0x74, 0x73, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12,
	0x29, 0x0a, 0x06, 0x61, 0x6c, 0x65, 0x72, 0x74, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32,
	0x11, 0x2e, 0x6b, 0x69, 0x6f, 0x72, 0x61, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x41, 0x6c, 0x65,
	0x72, 0x74, 0x52, 0x06, 0x61, 0x6c, 0x65, 0x72, 0x74, 0x73, 0x22, 0xa0, 0x03, 0x0a, 0x05, 0x41,
	0x6c, 0x65, 0x72, 0x74, 0x12, 0x35, 0x0a, 0x06, 0x6c, 0x61, 0x62, 0x65, 0x6c, 0x73, 0x18, 0x01,
	0x20, 0x03, 0x28, 0x0b, 0x32, 0x1d, 0x2e, 0x6b, 0x69, 0x6f, 0x72, 0x61, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x2e, 0x41, 0x6c, 0x65, 0x72, 0x74, 0x2e, 0x4c, 0x61, 0x62, 0x65, 0x6c, 0x73, 0x45, 0x6e,
	0x74, 0x72, 0x79, 0x52, 0x06, 0x6c, 0x61, 0x62, 0x65, 0x6c, 0x73, 0x12, 0x44, 0x0a, 0x0b, 0x61,
	0x6e, 0x6e, 0x6f, 0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0b,
	0x32, 0x22, 0x2e, 0x6b, 0x69, 0x6f, 0x72, 0x61, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x41, 0x6c,
	0x65, 0x72, 0x74, 0x2e, 0x41, 0x6e, 0x6e, 0x6f, 0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x45,
	0x6e, 0x74, 0x72, 0x79, 0x52, 0x0b, 0x61, 0x6e, 0x6e, 0x6f, 0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e,
	0x73, 0x12, 0x2f, 0x0a, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x18, 0x03, 0x20, 0x01, 0x28,
	0x0e, 0x32, 0x17, 0x2e, 0x6b, 0x69, 0x6f, 0x72, 0x61, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x41,
	0x6c, 0x65, 0x72, 0x74, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x52, 0x06, 0x73, 0x74, 0x61, 0x74,
	0x75, 0x73, 0x12, 0x38, 0x0a, 0x09, 0x73, 0x74, 0x61, 0x72, 0x74, 0x54, 0x69, 0x6d, 0x65, 0x18,
	0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d,
	0x70, 0x52, 0x09, 0x73, 0x74, 0x61, 0x72, 0x74, 0x54, 0x69, 0x6d, 0x65, 0x12, 0x34, 0x0a, 0x07,
	0x65, 0x6e, 0x64, 0x54, 0x69, 0x6d, 0x65, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e,
	0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e,
	0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x07, 0x65, 0x6e, 0x64, 0x54, 0x69,
	0x6d, 0x65, 0x1a, 0x39, 0x0a, 0x0b, 0x4c, 0x61, 0x62, 0x65, 0x6c, 0x73, 0x45, 0x6e, 0x74, 0x72,
	0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03,
	0x6b, 0x65, 0x79, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38, 0x01, 0x1a, 0x3e, 0x0a,
	0x10, 0x41, 0x6e, 0x6e, 0x6f, 0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x45, 0x6e, 0x74, 0x72,
	0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03,
	0x6b, 0x65, 0x79, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38, 0x01, 0x2a, 0x37, 0x0a,
	0x0b, 0x41, 0x6c, 0x65, 0x72, 0x74, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x0a, 0x0a, 0x06,
	0x66, 0x69, 0x72, 0x69, 0x6e, 0x67, 0x10, 0x00, 0x12, 0x0c, 0x0a, 0x08, 0x72, 0x65, 0x73, 0x6f,
	0x6c, 0x76, 0x65, 0x64, 0x10, 0x01, 0x12, 0x0e, 0x0a, 0x0a, 0x70, 0x72, 0x6f, 0x63, 0x65, 0x73,
	0x73, 0x69, 0x6e, 0x67, 0x10, 0x02, 0x42, 0x31, 0x5a, 0x2f, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62,
	0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x73, 0x69, 0x6e, 0x6b, 0x69, 0x6e, 0x67, 0x70, 0x6f, 0x69, 0x6e,
	0x74, 0x2f, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x2f, 0x64, 0x74, 0x6f, 0x2f, 0x6b,
	0x69, 0x6f, 0x72, 0x61, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x33,
}

var (
	file_alert_proto_rawDescOnce sync.Once
	file_alert_proto_rawDescData = file_alert_proto_rawDesc
)

func file_alert_proto_rawDescGZIP() []byte {
	file_alert_proto_rawDescOnce.Do(func() {
		file_alert_proto_rawDescData = protoimpl.X.CompressGZIP(file_alert_proto_rawDescData)
	})
	return file_alert_proto_rawDescData
}

var file_alert_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_alert_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_alert_proto_goTypes = []interface{}{
	(AlertStatus)(0),              // 0: kioraproto.AlertStatus
	(*PostAlertsMessage)(nil),     // 1: kioraproto.PostAlertsMessage
	(*Alert)(nil),                 // 2: kioraproto.Alert
	nil,                           // 3: kioraproto.Alert.LabelsEntry
	nil,                           // 4: kioraproto.Alert.AnnotationsEntry
	(*timestamppb.Timestamp)(nil), // 5: google.protobuf.Timestamp
}
var file_alert_proto_depIdxs = []int32{
	2, // 0: kioraproto.PostAlertsMessage.alerts:type_name -> kioraproto.Alert
	3, // 1: kioraproto.Alert.labels:type_name -> kioraproto.Alert.LabelsEntry
	4, // 2: kioraproto.Alert.annotations:type_name -> kioraproto.Alert.AnnotationsEntry
	0, // 3: kioraproto.Alert.status:type_name -> kioraproto.AlertStatus
	5, // 4: kioraproto.Alert.startTime:type_name -> google.protobuf.Timestamp
	5, // 5: kioraproto.Alert.endTime:type_name -> google.protobuf.Timestamp
	6, // [6:6] is the sub-list for method output_type
	6, // [6:6] is the sub-list for method input_type
	6, // [6:6] is the sub-list for extension type_name
	6, // [6:6] is the sub-list for extension extendee
	0, // [0:6] is the sub-list for field type_name
}

func init() { file_alert_proto_init() }
func file_alert_proto_init() {
	if File_alert_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_alert_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PostAlertsMessage); i {
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
		file_alert_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Alert); i {
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
			RawDescriptor: file_alert_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_alert_proto_goTypes,
		DependencyIndexes: file_alert_proto_depIdxs,
		EnumInfos:         file_alert_proto_enumTypes,
		MessageInfos:      file_alert_proto_msgTypes,
	}.Build()
	File_alert_proto = out.File
	file_alert_proto_rawDesc = nil
	file_alert_proto_goTypes = nil
	file_alert_proto_depIdxs = nil
}
