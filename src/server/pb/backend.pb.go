// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v3.6.1
// source: backend.proto

package __

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

// Category types
type Category int32

const (
	Category_CATEGORY_PASS Category = 0
	Category_CATEGORY_TEXT Category = 1
	Category_CATEGORY_CARD Category = 2
	Category_CATEGORY_FILE Category = 3
)

// Enum value maps for Category.
var (
	Category_name = map[int32]string{
		0: "CATEGORY_PASS",
		1: "CATEGORY_TEXT",
		2: "CATEGORY_CARD",
		3: "CATEGORY_FILE",
	}
	Category_value = map[string]int32{
		"CATEGORY_PASS": 0,
		"CATEGORY_TEXT": 1,
		"CATEGORY_CARD": 2,
		"CATEGORY_FILE": 3,
	}
)

func (x Category) Enum() *Category {
	p := new(Category)
	*p = x
	return p
}

func (x Category) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (Category) Descriptor() protoreflect.EnumDescriptor {
	return file_backend_proto_enumTypes[0].Descriptor()
}

func (Category) Type() protoreflect.EnumType {
	return &file_backend_proto_enumTypes[0]
}

func (x Category) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use Category.Descriptor instead.
func (Category) EnumDescriptor() ([]byte, []int) {
	return file_backend_proto_rawDescGZIP(), []int{0}
}

type Empty_DTO struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *Empty_DTO) Reset() {
	*x = Empty_DTO{}
	if protoimpl.UnsafeEnabled {
		mi := &file_backend_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Empty_DTO) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Empty_DTO) ProtoMessage() {}

func (x *Empty_DTO) ProtoReflect() protoreflect.Message {
	mi := &file_backend_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Empty_DTO.ProtoReflect.Descriptor instead.
func (*Empty_DTO) Descriptor() ([]byte, []int) {
	return file_backend_proto_rawDescGZIP(), []int{0}
}

type LoginPassPair struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Login    string `protobuf:"bytes,1,opt,name=Login,proto3" json:"Login,omitempty"`
	Password string `protobuf:"bytes,2,opt,name=Password,proto3" json:"Password,omitempty"`
}

func (x *LoginPassPair) Reset() {
	*x = LoginPassPair{}
	if protoimpl.UnsafeEnabled {
		mi := &file_backend_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *LoginPassPair) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LoginPassPair) ProtoMessage() {}

func (x *LoginPassPair) ProtoReflect() protoreflect.Message {
	mi := &file_backend_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LoginPassPair.ProtoReflect.Descriptor instead.
func (*LoginPassPair) Descriptor() ([]byte, []int) {
	return file_backend_proto_rawDescGZIP(), []int{1}
}

func (x *LoginPassPair) GetLogin() string {
	if x != nil {
		return x.Login
	}
	return ""
}

func (x *LoginPassPair) GetPassword() string {
	if x != nil {
		return x.Password
	}
	return ""
}

type SessionID_DTO struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	SID string `protobuf:"bytes,1,opt,name=SID,proto3" json:"SID,omitempty"`
}

func (x *SessionID_DTO) Reset() {
	*x = SessionID_DTO{}
	if protoimpl.UnsafeEnabled {
		mi := &file_backend_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SessionID_DTO) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SessionID_DTO) ProtoMessage() {}

func (x *SessionID_DTO) ProtoReflect() protoreflect.Message {
	mi := &file_backend_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SessionID_DTO.ProtoReflect.Descriptor instead.
func (*SessionID_DTO) Descriptor() ([]byte, []int) {
	return file_backend_proto_rawDescGZIP(), []int{2}
}

func (x *SessionID_DTO) GetSID() string {
	if x != nil {
		return x.SID
	}
	return ""
}

type GetCategoryInfo_DTO struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Category Category `protobuf:"varint,1,opt,name=Category,proto3,enum=internal.server.pb.Category" json:"Category,omitempty"`
}

func (x *GetCategoryInfo_DTO) Reset() {
	*x = GetCategoryInfo_DTO{}
	if protoimpl.UnsafeEnabled {
		mi := &file_backend_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetCategoryInfo_DTO) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetCategoryInfo_DTO) ProtoMessage() {}

func (x *GetCategoryInfo_DTO) ProtoReflect() protoreflect.Message {
	mi := &file_backend_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetCategoryInfo_DTO.ProtoReflect.Descriptor instead.
func (*GetCategoryInfo_DTO) Descriptor() ([]byte, []int) {
	return file_backend_proto_rawDescGZIP(), []int{3}
}

func (x *GetCategoryInfo_DTO) GetCategory() Category {
	if x != nil {
		return x.Category
	}
	return Category_CATEGORY_PASS
}

type DataInfo struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	DataID   string `protobuf:"bytes,1,opt,name=DataID,proto3" json:"DataID,omitempty"`
	Metadata string `protobuf:"bytes,2,opt,name=Metadata,proto3" json:"Metadata,omitempty"`
}

func (x *DataInfo) Reset() {
	*x = DataInfo{}
	if protoimpl.UnsafeEnabled {
		mi := &file_backend_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DataInfo) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DataInfo) ProtoMessage() {}

func (x *DataInfo) ProtoReflect() protoreflect.Message {
	mi := &file_backend_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DataInfo.ProtoReflect.Descriptor instead.
func (*DataInfo) Descriptor() ([]byte, []int) {
	return file_backend_proto_rawDescGZIP(), []int{4}
}

func (x *DataInfo) GetDataID() string {
	if x != nil {
		return x.DataID
	}
	return ""
}

func (x *DataInfo) GetMetadata() string {
	if x != nil {
		return x.Metadata
	}
	return ""
}

type CategoryHead_DTO struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Info []*DataInfo `protobuf:"bytes,1,rep,name=info,proto3" json:"info,omitempty"`
}

func (x *CategoryHead_DTO) Reset() {
	*x = CategoryHead_DTO{}
	if protoimpl.UnsafeEnabled {
		mi := &file_backend_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CategoryHead_DTO) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CategoryHead_DTO) ProtoMessage() {}

func (x *CategoryHead_DTO) ProtoReflect() protoreflect.Message {
	mi := &file_backend_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CategoryHead_DTO.ProtoReflect.Descriptor instead.
func (*CategoryHead_DTO) Descriptor() ([]byte, []int) {
	return file_backend_proto_rawDescGZIP(), []int{5}
}

func (x *CategoryHead_DTO) GetInfo() []*DataInfo {
	if x != nil {
		return x.Info
	}
	return nil
}

// data types to transfer
type DataID_DTO struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ID string `protobuf:"bytes,1,opt,name=ID,proto3" json:"ID,omitempty"`
}

func (x *DataID_DTO) Reset() {
	*x = DataID_DTO{}
	if protoimpl.UnsafeEnabled {
		mi := &file_backend_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DataID_DTO) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DataID_DTO) ProtoMessage() {}

func (x *DataID_DTO) ProtoReflect() protoreflect.Message {
	mi := &file_backend_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DataID_DTO.ProtoReflect.Descriptor instead.
func (*DataID_DTO) Descriptor() ([]byte, []int) {
	return file_backend_proto_rawDescGZIP(), []int{6}
}

func (x *DataID_DTO) GetID() string {
	if x != nil {
		return x.ID
	}
	return ""
}

type SecureData_DTO struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Data     []byte `protobuf:"bytes,1,opt,name=Data,proto3" json:"Data,omitempty"`
	Metadata string `protobuf:"bytes,2,opt,name=Metadata,proto3" json:"Metadata,omitempty"`
}

func (x *SecureData_DTO) Reset() {
	*x = SecureData_DTO{}
	if protoimpl.UnsafeEnabled {
		mi := &file_backend_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SecureData_DTO) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SecureData_DTO) ProtoMessage() {}

func (x *SecureData_DTO) ProtoReflect() protoreflect.Message {
	mi := &file_backend_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SecureData_DTO.ProtoReflect.Descriptor instead.
func (*SecureData_DTO) Descriptor() ([]byte, []int) {
	return file_backend_proto_rawDescGZIP(), []int{7}
}

func (x *SecureData_DTO) GetData() []byte {
	if x != nil {
		return x.Data
	}
	return nil
}

func (x *SecureData_DTO) GetMetadata() string {
	if x != nil {
		return x.Metadata
	}
	return ""
}

var File_backend_proto protoreflect.FileDescriptor

var file_backend_proto_rawDesc = []byte{
	0x0a, 0x0d, 0x62, 0x61, 0x63, 0x6b, 0x65, 0x6e, 0x64, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12,
	0x12, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x2e, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72,
	0x2e, 0x70, 0x62, 0x22, 0x0b, 0x0a, 0x09, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x5f, 0x44, 0x54, 0x4f,
	0x22, 0x41, 0x0a, 0x0d, 0x4c, 0x6f, 0x67, 0x69, 0x6e, 0x50, 0x61, 0x73, 0x73, 0x50, 0x61, 0x69,
	0x72, 0x12, 0x14, 0x0a, 0x05, 0x4c, 0x6f, 0x67, 0x69, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x05, 0x4c, 0x6f, 0x67, 0x69, 0x6e, 0x12, 0x1a, 0x0a, 0x08, 0x50, 0x61, 0x73, 0x73, 0x77,
	0x6f, 0x72, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x50, 0x61, 0x73, 0x73, 0x77,
	0x6f, 0x72, 0x64, 0x22, 0x21, 0x0a, 0x0d, 0x53, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x49, 0x44,
	0x5f, 0x44, 0x54, 0x4f, 0x12, 0x10, 0x0a, 0x03, 0x53, 0x49, 0x44, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x03, 0x53, 0x49, 0x44, 0x22, 0x4f, 0x0a, 0x13, 0x47, 0x65, 0x74, 0x43, 0x61, 0x74,
	0x65, 0x67, 0x6f, 0x72, 0x79, 0x49, 0x6e, 0x66, 0x6f, 0x5f, 0x44, 0x54, 0x4f, 0x12, 0x38, 0x0a,
	0x08, 0x43, 0x61, 0x74, 0x65, 0x67, 0x6f, 0x72, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0e, 0x32,
	0x1c, 0x2e, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x2e, 0x73, 0x65, 0x72, 0x76, 0x65,
	0x72, 0x2e, 0x70, 0x62, 0x2e, 0x43, 0x61, 0x74, 0x65, 0x67, 0x6f, 0x72, 0x79, 0x52, 0x08, 0x43,
	0x61, 0x74, 0x65, 0x67, 0x6f, 0x72, 0x79, 0x22, 0x3e, 0x0a, 0x08, 0x44, 0x61, 0x74, 0x61, 0x49,
	0x6e, 0x66, 0x6f, 0x12, 0x16, 0x0a, 0x06, 0x44, 0x61, 0x74, 0x61, 0x49, 0x44, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x06, 0x44, 0x61, 0x74, 0x61, 0x49, 0x44, 0x12, 0x1a, 0x0a, 0x08, 0x4d,
	0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x4d,
	0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x22, 0x44, 0x0a, 0x10, 0x43, 0x61, 0x74, 0x65, 0x67,
	0x6f, 0x72, 0x79, 0x48, 0x65, 0x61, 0x64, 0x5f, 0x44, 0x54, 0x4f, 0x12, 0x30, 0x0a, 0x04, 0x69,
	0x6e, 0x66, 0x6f, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x1c, 0x2e, 0x69, 0x6e, 0x74, 0x65,
	0x72, 0x6e, 0x61, 0x6c, 0x2e, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x2e, 0x70, 0x62, 0x2e, 0x44,
	0x61, 0x74, 0x61, 0x49, 0x6e, 0x66, 0x6f, 0x52, 0x04, 0x69, 0x6e, 0x66, 0x6f, 0x22, 0x1c, 0x0a,
	0x0a, 0x44, 0x61, 0x74, 0x61, 0x49, 0x44, 0x5f, 0x44, 0x54, 0x4f, 0x12, 0x0e, 0x0a, 0x02, 0x49,
	0x44, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x49, 0x44, 0x22, 0x40, 0x0a, 0x0e, 0x53,
	0x65, 0x63, 0x75, 0x72, 0x65, 0x44, 0x61, 0x74, 0x61, 0x5f, 0x44, 0x54, 0x4f, 0x12, 0x12, 0x0a,
	0x04, 0x44, 0x61, 0x74, 0x61, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x04, 0x44, 0x61, 0x74,
	0x61, 0x12, 0x1a, 0x0a, 0x08, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x08, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x2a, 0x56, 0x0a,
	0x08, 0x43, 0x61, 0x74, 0x65, 0x67, 0x6f, 0x72, 0x79, 0x12, 0x11, 0x0a, 0x0d, 0x43, 0x41, 0x54,
	0x45, 0x47, 0x4f, 0x52, 0x59, 0x5f, 0x50, 0x41, 0x53, 0x53, 0x10, 0x00, 0x12, 0x11, 0x0a, 0x0d,
	0x43, 0x41, 0x54, 0x45, 0x47, 0x4f, 0x52, 0x59, 0x5f, 0x54, 0x45, 0x58, 0x54, 0x10, 0x01, 0x12,
	0x11, 0x0a, 0x0d, 0x43, 0x41, 0x54, 0x45, 0x47, 0x4f, 0x52, 0x59, 0x5f, 0x43, 0x41, 0x52, 0x44,
	0x10, 0x02, 0x12, 0x11, 0x0a, 0x0d, 0x43, 0x41, 0x54, 0x45, 0x47, 0x4f, 0x52, 0x59, 0x5f, 0x46,
	0x49, 0x4c, 0x45, 0x10, 0x03, 0x32, 0x8c, 0x04, 0x0a, 0x0a, 0x47, 0x6f, 0x70, 0x68, 0x4b, 0x65,
	0x65, 0x70, 0x65, 0x72, 0x12, 0x46, 0x0a, 0x04, 0x50, 0x69, 0x6e, 0x67, 0x12, 0x1d, 0x2e, 0x69,
	0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x2e, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x2e, 0x70,
	0x62, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x5f, 0x44, 0x54, 0x4f, 0x1a, 0x1d, 0x2e, 0x69, 0x6e,
	0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x2e, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x2e, 0x70, 0x62,
	0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x5f, 0x44, 0x54, 0x4f, 0x22, 0x00, 0x12, 0x4f, 0x0a, 0x05,
	0x4c, 0x6f, 0x67, 0x69, 0x6e, 0x12, 0x21, 0x2e, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c,
	0x2e, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x2e, 0x70, 0x62, 0x2e, 0x4c, 0x6f, 0x67, 0x69, 0x6e,
	0x50, 0x61, 0x73, 0x73, 0x50, 0x61, 0x69, 0x72, 0x1a, 0x21, 0x2e, 0x69, 0x6e, 0x74, 0x65, 0x72,
	0x6e, 0x61, 0x6c, 0x2e, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x2e, 0x70, 0x62, 0x2e, 0x53, 0x65,
	0x73, 0x73, 0x69, 0x6f, 0x6e, 0x49, 0x44, 0x5f, 0x44, 0x54, 0x4f, 0x22, 0x00, 0x12, 0x5a, 0x0a,
	0x10, 0x4b, 0x69, 0x63, 0x6b, 0x4f, 0x74, 0x68, 0x65, 0x72, 0x53, 0x65, 0x73, 0x73, 0x69, 0x6f,
	0x6e, 0x12, 0x21, 0x2e, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x2e, 0x73, 0x65, 0x72,
	0x76, 0x65, 0x72, 0x2e, 0x70, 0x62, 0x2e, 0x4c, 0x6f, 0x67, 0x69, 0x6e, 0x50, 0x61, 0x73, 0x73,
	0x50, 0x61, 0x69, 0x72, 0x1a, 0x21, 0x2e, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x2e,
	0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x2e, 0x70, 0x62, 0x2e, 0x53, 0x65, 0x73, 0x73, 0x69, 0x6f,
	0x6e, 0x49, 0x44, 0x5f, 0x44, 0x54, 0x4f, 0x22, 0x00, 0x12, 0x52, 0x0a, 0x08, 0x52, 0x65, 0x67,
	0x69, 0x73, 0x74, 0x65, 0x72, 0x12, 0x21, 0x2e, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c,
	0x2e, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x2e, 0x70, 0x62, 0x2e, 0x4c, 0x6f, 0x67, 0x69, 0x6e,
	0x50, 0x61, 0x73, 0x73, 0x50, 0x61, 0x69, 0x72, 0x1a, 0x21, 0x2e, 0x69, 0x6e, 0x74, 0x65, 0x72,
	0x6e, 0x61, 0x6c, 0x2e, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x2e, 0x70, 0x62, 0x2e, 0x53, 0x65,
	0x73, 0x73, 0x69, 0x6f, 0x6e, 0x49, 0x44, 0x5f, 0x44, 0x54, 0x4f, 0x22, 0x00, 0x12, 0x5a, 0x0a,
	0x12, 0x53, 0x74, 0x6f, 0x72, 0x65, 0x4c, 0x6f, 0x67, 0x69, 0x6e, 0x50, 0x61, 0x73, 0x73, 0x77,
	0x6f, 0x72, 0x64, 0x12, 0x22, 0x2e, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x2e, 0x73,
	0x65, 0x72, 0x76, 0x65, 0x72, 0x2e, 0x70, 0x62, 0x2e, 0x53, 0x65, 0x63, 0x75, 0x72, 0x65, 0x44,
	0x61, 0x74, 0x61, 0x5f, 0x44, 0x54, 0x4f, 0x1a, 0x1e, 0x2e, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e,
	0x61, 0x6c, 0x2e, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x2e, 0x70, 0x62, 0x2e, 0x44, 0x61, 0x74,
	0x61, 0x49, 0x44, 0x5f, 0x44, 0x54, 0x4f, 0x22, 0x00, 0x12, 0x59, 0x0a, 0x11, 0x4c, 0x6f, 0x61,
	0x64, 0x4c, 0x6f, 0x67, 0x69, 0x6e, 0x50, 0x61, 0x73, 0x73, 0x77, 0x6f, 0x72, 0x64, 0x12, 0x1e,
	0x2e, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x2e, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72,
	0x2e, 0x70, 0x62, 0x2e, 0x44, 0x61, 0x74, 0x61, 0x49, 0x44, 0x5f, 0x44, 0x54, 0x4f, 0x1a, 0x22,
	0x2e, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x2e, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72,
	0x2e, 0x70, 0x62, 0x2e, 0x53, 0x65, 0x63, 0x75, 0x72, 0x65, 0x44, 0x61, 0x74, 0x61, 0x5f, 0x44,
	0x54, 0x4f, 0x22, 0x00, 0x42, 0x03, 0x5a, 0x01, 0x2e, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x33,
}

var (
	file_backend_proto_rawDescOnce sync.Once
	file_backend_proto_rawDescData = file_backend_proto_rawDesc
)

func file_backend_proto_rawDescGZIP() []byte {
	file_backend_proto_rawDescOnce.Do(func() {
		file_backend_proto_rawDescData = protoimpl.X.CompressGZIP(file_backend_proto_rawDescData)
	})
	return file_backend_proto_rawDescData
}

var file_backend_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_backend_proto_msgTypes = make([]protoimpl.MessageInfo, 8)
var file_backend_proto_goTypes = []interface{}{
	(Category)(0),               // 0: internal.server.pb.Category
	(*Empty_DTO)(nil),           // 1: internal.server.pb.Empty_DTO
	(*LoginPassPair)(nil),       // 2: internal.server.pb.LoginPassPair
	(*SessionID_DTO)(nil),       // 3: internal.server.pb.SessionID_DTO
	(*GetCategoryInfo_DTO)(nil), // 4: internal.server.pb.GetCategoryInfo_DTO
	(*DataInfo)(nil),            // 5: internal.server.pb.DataInfo
	(*CategoryHead_DTO)(nil),    // 6: internal.server.pb.CategoryHead_DTO
	(*DataID_DTO)(nil),          // 7: internal.server.pb.DataID_DTO
	(*SecureData_DTO)(nil),      // 8: internal.server.pb.SecureData_DTO
}
var file_backend_proto_depIdxs = []int32{
	0, // 0: internal.server.pb.GetCategoryInfo_DTO.Category:type_name -> internal.server.pb.Category
	5, // 1: internal.server.pb.CategoryHead_DTO.info:type_name -> internal.server.pb.DataInfo
	1, // 2: internal.server.pb.GophKeeper.Ping:input_type -> internal.server.pb.Empty_DTO
	2, // 3: internal.server.pb.GophKeeper.Login:input_type -> internal.server.pb.LoginPassPair
	2, // 4: internal.server.pb.GophKeeper.KickOtherSession:input_type -> internal.server.pb.LoginPassPair
	2, // 5: internal.server.pb.GophKeeper.Register:input_type -> internal.server.pb.LoginPassPair
	8, // 6: internal.server.pb.GophKeeper.StoreLoginPassword:input_type -> internal.server.pb.SecureData_DTO
	7, // 7: internal.server.pb.GophKeeper.LoadLoginPassword:input_type -> internal.server.pb.DataID_DTO
	1, // 8: internal.server.pb.GophKeeper.Ping:output_type -> internal.server.pb.Empty_DTO
	3, // 9: internal.server.pb.GophKeeper.Login:output_type -> internal.server.pb.SessionID_DTO
	3, // 10: internal.server.pb.GophKeeper.KickOtherSession:output_type -> internal.server.pb.SessionID_DTO
	3, // 11: internal.server.pb.GophKeeper.Register:output_type -> internal.server.pb.SessionID_DTO
	7, // 12: internal.server.pb.GophKeeper.StoreLoginPassword:output_type -> internal.server.pb.DataID_DTO
	8, // 13: internal.server.pb.GophKeeper.LoadLoginPassword:output_type -> internal.server.pb.SecureData_DTO
	8, // [8:14] is the sub-list for method output_type
	2, // [2:8] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_backend_proto_init() }
func file_backend_proto_init() {
	if File_backend_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_backend_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Empty_DTO); i {
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
		file_backend_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*LoginPassPair); i {
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
		file_backend_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SessionID_DTO); i {
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
		file_backend_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetCategoryInfo_DTO); i {
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
		file_backend_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DataInfo); i {
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
		file_backend_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CategoryHead_DTO); i {
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
		file_backend_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DataID_DTO); i {
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
		file_backend_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SecureData_DTO); i {
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
			RawDescriptor: file_backend_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   8,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_backend_proto_goTypes,
		DependencyIndexes: file_backend_proto_depIdxs,
		EnumInfos:         file_backend_proto_enumTypes,
		MessageInfos:      file_backend_proto_msgTypes,
	}.Build()
	File_backend_proto = out.File
	file_backend_proto_rawDesc = nil
	file_backend_proto_goTypes = nil
	file_backend_proto_depIdxs = nil
}