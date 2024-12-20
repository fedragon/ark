// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.35.2
// 	protoc        (unknown)
// source: ark/v1/rpc.proto

package arkv1

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

type Metadata struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Hash      []byte                 `protobuf:"bytes,1,opt,name=hash,proto3" json:"hash,omitempty"`
	Name      string                 `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	Size      int64                  `protobuf:"varint,3,opt,name=size,proto3" json:"size,omitempty"`
	CreatedAt *timestamppb.Timestamp `protobuf:"bytes,10,opt,name=created_at,json=createdAt,proto3" json:"created_at,omitempty"`
}

func (x *Metadata) Reset() {
	*x = Metadata{}
	mi := &file_ark_v1_rpc_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Metadata) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Metadata) ProtoMessage() {}

func (x *Metadata) ProtoReflect() protoreflect.Message {
	mi := &file_ark_v1_rpc_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Metadata.ProtoReflect.Descriptor instead.
func (*Metadata) Descriptor() ([]byte, []int) {
	return file_ark_v1_rpc_proto_rawDescGZIP(), []int{0}
}

func (x *Metadata) GetHash() []byte {
	if x != nil {
		return x.Hash
	}
	return nil
}

func (x *Metadata) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *Metadata) GetSize() int64 {
	if x != nil {
		return x.Size
	}
	return 0
}

func (x *Metadata) GetCreatedAt() *timestamppb.Timestamp {
	if x != nil {
		return x.CreatedAt
	}
	return nil
}

type Chunk struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Data []byte `protobuf:"bytes,1,opt,name=data,proto3" json:"data,omitempty"`
}

func (x *Chunk) Reset() {
	*x = Chunk{}
	mi := &file_ark_v1_rpc_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Chunk) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Chunk) ProtoMessage() {}

func (x *Chunk) ProtoReflect() protoreflect.Message {
	mi := &file_ark_v1_rpc_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Chunk.ProtoReflect.Descriptor instead.
func (*Chunk) Descriptor() ([]byte, []int) {
	return file_ark_v1_rpc_proto_rawDescGZIP(), []int{1}
}

func (x *Chunk) GetData() []byte {
	if x != nil {
		return x.Data
	}
	return nil
}

type UploadFileRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Types that are assignable to File:
	//
	//	*UploadFileRequest_Metadata
	//	*UploadFileRequest_Chunk
	File isUploadFileRequest_File `protobuf_oneof:"file"`
}

func (x *UploadFileRequest) Reset() {
	*x = UploadFileRequest{}
	mi := &file_ark_v1_rpc_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *UploadFileRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UploadFileRequest) ProtoMessage() {}

func (x *UploadFileRequest) ProtoReflect() protoreflect.Message {
	mi := &file_ark_v1_rpc_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UploadFileRequest.ProtoReflect.Descriptor instead.
func (*UploadFileRequest) Descriptor() ([]byte, []int) {
	return file_ark_v1_rpc_proto_rawDescGZIP(), []int{2}
}

func (m *UploadFileRequest) GetFile() isUploadFileRequest_File {
	if m != nil {
		return m.File
	}
	return nil
}

func (x *UploadFileRequest) GetMetadata() *Metadata {
	if x, ok := x.GetFile().(*UploadFileRequest_Metadata); ok {
		return x.Metadata
	}
	return nil
}

func (x *UploadFileRequest) GetChunk() *Chunk {
	if x, ok := x.GetFile().(*UploadFileRequest_Chunk); ok {
		return x.Chunk
	}
	return nil
}

type isUploadFileRequest_File interface {
	isUploadFileRequest_File()
}

type UploadFileRequest_Metadata struct {
	Metadata *Metadata `protobuf:"bytes,1,opt,name=metadata,proto3,oneof"`
}

type UploadFileRequest_Chunk struct {
	Chunk *Chunk `protobuf:"bytes,2,opt,name=chunk,proto3,oneof"`
}

func (*UploadFileRequest_Metadata) isUploadFileRequest_File() {}

func (*UploadFileRequest_Chunk) isUploadFileRequest_File() {}

type UploadFileResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Details string `protobuf:"bytes,1,opt,name=details,proto3" json:"details,omitempty"`
}

func (x *UploadFileResponse) Reset() {
	*x = UploadFileResponse{}
	mi := &file_ark_v1_rpc_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *UploadFileResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UploadFileResponse) ProtoMessage() {}

func (x *UploadFileResponse) ProtoReflect() protoreflect.Message {
	mi := &file_ark_v1_rpc_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UploadFileResponse.ProtoReflect.Descriptor instead.
func (*UploadFileResponse) Descriptor() ([]byte, []int) {
	return file_ark_v1_rpc_proto_rawDescGZIP(), []int{3}
}

func (x *UploadFileResponse) GetDetails() string {
	if x != nil {
		return x.Details
	}
	return ""
}

var File_ark_v1_rpc_proto protoreflect.FileDescriptor

var file_ark_v1_rpc_proto_rawDesc = []byte{
	0x0a, 0x10, 0x61, 0x72, 0x6b, 0x2f, 0x76, 0x31, 0x2f, 0x72, 0x70, 0x63, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x12, 0x06, 0x61, 0x72, 0x6b, 0x2e, 0x76, 0x31, 0x1a, 0x1f, 0x67, 0x6f, 0x6f, 0x67,
	0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x74, 0x69, 0x6d, 0x65,
	0x73, 0x74, 0x61, 0x6d, 0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x81, 0x01, 0x0a, 0x08,
	0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x12, 0x12, 0x0a, 0x04, 0x68, 0x61, 0x73, 0x68,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x04, 0x68, 0x61, 0x73, 0x68, 0x12, 0x12, 0x0a, 0x04,
	0x6e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65,
	0x12, 0x12, 0x0a, 0x04, 0x73, 0x69, 0x7a, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x03, 0x52, 0x04,
	0x73, 0x69, 0x7a, 0x65, 0x12, 0x39, 0x0a, 0x0a, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x5f,
	0x61, 0x74, 0x18, 0x0a, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c,
	0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73,
	0x74, 0x61, 0x6d, 0x70, 0x52, 0x09, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x41, 0x74, 0x22,
	0x1b, 0x0a, 0x05, 0x43, 0x68, 0x75, 0x6e, 0x6b, 0x12, 0x12, 0x0a, 0x04, 0x64, 0x61, 0x74, 0x61,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x04, 0x64, 0x61, 0x74, 0x61, 0x22, 0x72, 0x0a, 0x11,
	0x55, 0x70, 0x6c, 0x6f, 0x61, 0x64, 0x46, 0x69, 0x6c, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x12, 0x2e, 0x0a, 0x08, 0x6d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x10, 0x2e, 0x61, 0x72, 0x6b, 0x2e, 0x76, 0x31, 0x2e, 0x4d, 0x65, 0x74,
	0x61, 0x64, 0x61, 0x74, 0x61, 0x48, 0x00, 0x52, 0x08, 0x6d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74,
	0x61, 0x12, 0x25, 0x0a, 0x05, 0x63, 0x68, 0x75, 0x6e, 0x6b, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x0d, 0x2e, 0x61, 0x72, 0x6b, 0x2e, 0x76, 0x31, 0x2e, 0x43, 0x68, 0x75, 0x6e, 0x6b, 0x48,
	0x00, 0x52, 0x05, 0x63, 0x68, 0x75, 0x6e, 0x6b, 0x42, 0x06, 0x0a, 0x04, 0x66, 0x69, 0x6c, 0x65,
	0x22, 0x2e, 0x0a, 0x12, 0x55, 0x70, 0x6c, 0x6f, 0x61, 0x64, 0x46, 0x69, 0x6c, 0x65, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x64, 0x65, 0x74, 0x61, 0x69, 0x6c,
	0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x64, 0x65, 0x74, 0x61, 0x69, 0x6c, 0x73,
	0x32, 0x51, 0x0a, 0x06, 0x41, 0x72, 0x6b, 0x41, 0x70, 0x69, 0x12, 0x47, 0x0a, 0x0a, 0x55, 0x70,
	0x6c, 0x6f, 0x61, 0x64, 0x46, 0x69, 0x6c, 0x65, 0x12, 0x19, 0x2e, 0x61, 0x72, 0x6b, 0x2e, 0x76,
	0x31, 0x2e, 0x55, 0x70, 0x6c, 0x6f, 0x61, 0x64, 0x46, 0x69, 0x6c, 0x65, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x1a, 0x1a, 0x2e, 0x61, 0x72, 0x6b, 0x2e, 0x76, 0x31, 0x2e, 0x55, 0x70, 0x6c,
	0x6f, 0x61, 0x64, 0x46, 0x69, 0x6c, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22,
	0x00, 0x28, 0x01, 0x42, 0x2a, 0x5a, 0x28, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f,
	0x6d, 0x2f, 0x66, 0x65, 0x64, 0x72, 0x61, 0x67, 0x6f, 0x6e, 0x2f, 0x61, 0x72, 0x6b, 0x2f, 0x67,
	0x65, 0x6e, 0x2f, 0x61, 0x72, 0x6b, 0x2f, 0x76, 0x31, 0x3b, 0x61, 0x72, 0x6b, 0x76, 0x31, 0x62,
	0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_ark_v1_rpc_proto_rawDescOnce sync.Once
	file_ark_v1_rpc_proto_rawDescData = file_ark_v1_rpc_proto_rawDesc
)

func file_ark_v1_rpc_proto_rawDescGZIP() []byte {
	file_ark_v1_rpc_proto_rawDescOnce.Do(func() {
		file_ark_v1_rpc_proto_rawDescData = protoimpl.X.CompressGZIP(file_ark_v1_rpc_proto_rawDescData)
	})
	return file_ark_v1_rpc_proto_rawDescData
}

var file_ark_v1_rpc_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_ark_v1_rpc_proto_goTypes = []any{
	(*Metadata)(nil),              // 0: ark.v1.Metadata
	(*Chunk)(nil),                 // 1: ark.v1.Chunk
	(*UploadFileRequest)(nil),     // 2: ark.v1.UploadFileRequest
	(*UploadFileResponse)(nil),    // 3: ark.v1.UploadFileResponse
	(*timestamppb.Timestamp)(nil), // 4: google.protobuf.Timestamp
}
var file_ark_v1_rpc_proto_depIdxs = []int32{
	4, // 0: ark.v1.Metadata.created_at:type_name -> google.protobuf.Timestamp
	0, // 1: ark.v1.UploadFileRequest.metadata:type_name -> ark.v1.Metadata
	1, // 2: ark.v1.UploadFileRequest.chunk:type_name -> ark.v1.Chunk
	2, // 3: ark.v1.ArkApi.UploadFile:input_type -> ark.v1.UploadFileRequest
	3, // 4: ark.v1.ArkApi.UploadFile:output_type -> ark.v1.UploadFileResponse
	4, // [4:5] is the sub-list for method output_type
	3, // [3:4] is the sub-list for method input_type
	3, // [3:3] is the sub-list for extension type_name
	3, // [3:3] is the sub-list for extension extendee
	0, // [0:3] is the sub-list for field type_name
}

func init() { file_ark_v1_rpc_proto_init() }
func file_ark_v1_rpc_proto_init() {
	if File_ark_v1_rpc_proto != nil {
		return
	}
	file_ark_v1_rpc_proto_msgTypes[2].OneofWrappers = []any{
		(*UploadFileRequest_Metadata)(nil),
		(*UploadFileRequest_Chunk)(nil),
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_ark_v1_rpc_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_ark_v1_rpc_proto_goTypes,
		DependencyIndexes: file_ark_v1_rpc_proto_depIdxs,
		MessageInfos:      file_ark_v1_rpc_proto_msgTypes,
	}.Build()
	File_ark_v1_rpc_proto = out.File
	file_ark_v1_rpc_proto_rawDesc = nil
	file_ark_v1_rpc_proto_goTypes = nil
	file_ark_v1_rpc_proto_depIdxs = nil
}
