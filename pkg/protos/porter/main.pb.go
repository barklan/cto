// Copyright 2015 gRPC authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.26.0
// 	protoc        v3.17.3
// source: pkg/protos/porter/main.proto

package porter

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

type ProjectAlertRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Project string `protobuf:"bytes,1,opt,name=project,proto3" json:"project,omitempty"`
	Message string `protobuf:"bytes,2,opt,name=message,proto3" json:"message,omitempty"`
}

func (x *ProjectAlertRequest) Reset() {
	*x = ProjectAlertRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_protos_porter_main_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ProjectAlertRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ProjectAlertRequest) ProtoMessage() {}

func (x *ProjectAlertRequest) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_protos_porter_main_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ProjectAlertRequest.ProtoReflect.Descriptor instead.
func (*ProjectAlertRequest) Descriptor() ([]byte, []int) {
	return file_pkg_protos_porter_main_proto_rawDescGZIP(), []int{0}
}

func (x *ProjectAlertRequest) GetProject() string {
	if x != nil {
		return x.Project
	}
	return ""
}

func (x *ProjectAlertRequest) GetMessage() string {
	if x != nil {
		return x.Message
	}
	return ""
}

type NewIssueRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Project   string `protobuf:"bytes,1,opt,name=project,proto3" json:"project,omitempty"`
	Key       string `protobuf:"bytes,2,opt,name=key,proto3" json:"key,omitempty"`
	Timestamp string `protobuf:"bytes,3,opt,name=timestamp,proto3" json:"timestamp,omitempty"`
	Env       string `protobuf:"bytes,4,opt,name=env,proto3" json:"env,omitempty"`
	Service   string `protobuf:"bytes,5,opt,name=service,proto3" json:"service,omitempty"`
	Flag      string `protobuf:"bytes,6,opt,name=flag,proto3" json:"flag,omitempty"`
}

func (x *NewIssueRequest) Reset() {
	*x = NewIssueRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_protos_porter_main_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *NewIssueRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*NewIssueRequest) ProtoMessage() {}

func (x *NewIssueRequest) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_protos_porter_main_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use NewIssueRequest.ProtoReflect.Descriptor instead.
func (*NewIssueRequest) Descriptor() ([]byte, []int) {
	return file_pkg_protos_porter_main_proto_rawDescGZIP(), []int{1}
}

func (x *NewIssueRequest) GetProject() string {
	if x != nil {
		return x.Project
	}
	return ""
}

func (x *NewIssueRequest) GetKey() string {
	if x != nil {
		return x.Key
	}
	return ""
}

func (x *NewIssueRequest) GetTimestamp() string {
	if x != nil {
		return x.Timestamp
	}
	return ""
}

func (x *NewIssueRequest) GetEnv() string {
	if x != nil {
		return x.Env
	}
	return ""
}

func (x *NewIssueRequest) GetService() string {
	if x != nil {
		return x.Service
	}
	return ""
}

func (x *NewIssueRequest) GetFlag() string {
	if x != nil {
		return x.Flag
	}
	return ""
}

type Message struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Message string `protobuf:"bytes,1,opt,name=message,proto3" json:"message,omitempty"`
}

func (x *Message) Reset() {
	*x = Message{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_protos_porter_main_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Message) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Message) ProtoMessage() {}

func (x *Message) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_protos_porter_main_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Message.ProtoReflect.Descriptor instead.
func (*Message) Descriptor() ([]byte, []int) {
	return file_pkg_protos_porter_main_proto_rawDescGZIP(), []int{2}
}

func (x *Message) GetMessage() string {
	if x != nil {
		return x.Message
	}
	return ""
}

var File_pkg_protos_porter_main_proto protoreflect.FileDescriptor

var file_pkg_protos_porter_main_proto_rawDesc = []byte{
	0x0a, 0x1c, 0x70, 0x6b, 0x67, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x73, 0x2f, 0x70, 0x6f, 0x72,
	0x74, 0x65, 0x72, 0x2f, 0x6d, 0x61, 0x69, 0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x06,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x73, 0x22, 0x49, 0x0a, 0x13, 0x50, 0x72, 0x6f, 0x6a, 0x65, 0x63,
	0x74, 0x41, 0x6c, 0x65, 0x72, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x18, 0x0a,
	0x07, 0x70, 0x72, 0x6f, 0x6a, 0x65, 0x63, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07,
	0x70, 0x72, 0x6f, 0x6a, 0x65, 0x63, 0x74, 0x12, 0x18, 0x0a, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61,
	0x67, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67,
	0x65, 0x22, 0x9b, 0x01, 0x0a, 0x0f, 0x4e, 0x65, 0x77, 0x49, 0x73, 0x73, 0x75, 0x65, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x18, 0x0a, 0x07, 0x70, 0x72, 0x6f, 0x6a, 0x65, 0x63, 0x74,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x70, 0x72, 0x6f, 0x6a, 0x65, 0x63, 0x74, 0x12,
	0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65,
	0x79, 0x12, 0x1c, 0x0a, 0x09, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x18, 0x03,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x12,
	0x10, 0x0a, 0x03, 0x65, 0x6e, 0x76, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x65, 0x6e,
	0x76, 0x12, 0x18, 0x0a, 0x07, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x18, 0x05, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x07, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x66,
	0x6c, 0x61, 0x67, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x66, 0x6c, 0x61, 0x67, 0x22,
	0x23, 0x0a, 0x07, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x6d, 0x65,
	0x73, 0x73, 0x61, 0x67, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x6d, 0x65, 0x73,
	0x73, 0x61, 0x67, 0x65, 0x32, 0xb5, 0x01, 0x0a, 0x06, 0x50, 0x6f, 0x72, 0x74, 0x65, 0x72, 0x12,
	0x3e, 0x0a, 0x0c, 0x50, 0x72, 0x6f, 0x6a, 0x65, 0x63, 0x74, 0x41, 0x6c, 0x65, 0x72, 0x74, 0x12,
	0x1b, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x73, 0x2e, 0x50, 0x72, 0x6f, 0x6a, 0x65, 0x63, 0x74,
	0x41, 0x6c, 0x65, 0x72, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x0f, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x73, 0x2e, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x22, 0x00, 0x12,
	0x33, 0x0a, 0x0d, 0x49, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x41, 0x6c, 0x65, 0x72, 0x74,
	0x12, 0x0f, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x73, 0x2e, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67,
	0x65, 0x1a, 0x0f, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x73, 0x2e, 0x4d, 0x65, 0x73, 0x73, 0x61,
	0x67, 0x65, 0x22, 0x00, 0x12, 0x36, 0x0a, 0x08, 0x4e, 0x65, 0x77, 0x49, 0x73, 0x73, 0x75, 0x65,
	0x12, 0x17, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x73, 0x2e, 0x4e, 0x65, 0x77, 0x49, 0x73, 0x73,
	0x75, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x0f, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x73, 0x2e, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x22, 0x00, 0x42, 0x2a, 0x5a, 0x28,
	0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x62, 0x61, 0x72, 0x6b, 0x6c,
	0x61, 0x6e, 0x2f, 0x63, 0x74, 0x6f, 0x2f, 0x70, 0x6b, 0x67, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x73, 0x2f, 0x70, 0x6f, 0x72, 0x74, 0x65, 0x72, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_pkg_protos_porter_main_proto_rawDescOnce sync.Once
	file_pkg_protos_porter_main_proto_rawDescData = file_pkg_protos_porter_main_proto_rawDesc
)

func file_pkg_protos_porter_main_proto_rawDescGZIP() []byte {
	file_pkg_protos_porter_main_proto_rawDescOnce.Do(func() {
		file_pkg_protos_porter_main_proto_rawDescData = protoimpl.X.CompressGZIP(file_pkg_protos_porter_main_proto_rawDescData)
	})
	return file_pkg_protos_porter_main_proto_rawDescData
}

var file_pkg_protos_porter_main_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_pkg_protos_porter_main_proto_goTypes = []interface{}{
	(*ProjectAlertRequest)(nil), // 0: protos.ProjectAlertRequest
	(*NewIssueRequest)(nil),     // 1: protos.NewIssueRequest
	(*Message)(nil),             // 2: protos.Message
}
var file_pkg_protos_porter_main_proto_depIdxs = []int32{
	0, // 0: protos.Porter.ProjectAlert:input_type -> protos.ProjectAlertRequest
	2, // 1: protos.Porter.InternalAlert:input_type -> protos.Message
	1, // 2: protos.Porter.NewIssue:input_type -> protos.NewIssueRequest
	2, // 3: protos.Porter.ProjectAlert:output_type -> protos.Message
	2, // 4: protos.Porter.InternalAlert:output_type -> protos.Message
	2, // 5: protos.Porter.NewIssue:output_type -> protos.Message
	3, // [3:6] is the sub-list for method output_type
	0, // [0:3] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_pkg_protos_porter_main_proto_init() }
func file_pkg_protos_porter_main_proto_init() {
	if File_pkg_protos_porter_main_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_pkg_protos_porter_main_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ProjectAlertRequest); i {
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
		file_pkg_protos_porter_main_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*NewIssueRequest); i {
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
		file_pkg_protos_porter_main_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Message); i {
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
			RawDescriptor: file_pkg_protos_porter_main_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_pkg_protos_porter_main_proto_goTypes,
		DependencyIndexes: file_pkg_protos_porter_main_proto_depIdxs,
		MessageInfos:      file_pkg_protos_porter_main_proto_msgTypes,
	}.Build()
	File_pkg_protos_porter_main_proto = out.File
	file_pkg_protos_porter_main_proto_rawDesc = nil
	file_pkg_protos_porter_main_proto_goTypes = nil
	file_pkg_protos_porter_main_proto_depIdxs = nil
}
