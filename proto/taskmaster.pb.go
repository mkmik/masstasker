// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.33.0
// 	protoc        v4.25.3
// source: taskmaster.proto

package taskmaster

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	anypb "google.golang.org/protobuf/types/known/anypb"
	durationpb "google.golang.org/protobuf/types/known/durationpb"
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

type UpdateRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Created    []*Task    `protobuf:"bytes,1,rep,name=created,proto3" json:"created,omitempty"`
	Deleted    []*TaskRef `protobuf:"bytes,2,rep,name=deleted,proto3" json:"deleted,omitempty"`
	Predicates []uint64   `protobuf:"varint,3,rep,packed,name=predicates,proto3" json:"predicates,omitempty"`
}

func (x *UpdateRequest) Reset() {
	*x = UpdateRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_taskmaster_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UpdateRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UpdateRequest) ProtoMessage() {}

func (x *UpdateRequest) ProtoReflect() protoreflect.Message {
	mi := &file_taskmaster_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UpdateRequest.ProtoReflect.Descriptor instead.
func (*UpdateRequest) Descriptor() ([]byte, []int) {
	return file_taskmaster_proto_rawDescGZIP(), []int{0}
}

func (x *UpdateRequest) GetCreated() []*Task {
	if x != nil {
		return x.Created
	}
	return nil
}

func (x *UpdateRequest) GetDeleted() []*TaskRef {
	if x != nil {
		return x.Deleted
	}
	return nil
}

func (x *UpdateRequest) GetPredicates() []uint64 {
	if x != nil {
		return x.Predicates
	}
	return nil
}

type UpdateResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	CreatedIds []uint64 `protobuf:"varint,1,rep,packed,name=created_ids,json=createdIds,proto3" json:"created_ids,omitempty"`
}

func (x *UpdateResponse) Reset() {
	*x = UpdateResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_taskmaster_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UpdateResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UpdateResponse) ProtoMessage() {}

func (x *UpdateResponse) ProtoReflect() protoreflect.Message {
	mi := &file_taskmaster_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UpdateResponse.ProtoReflect.Descriptor instead.
func (*UpdateResponse) Descriptor() ([]byte, []int) {
	return file_taskmaster_proto_rawDescGZIP(), []int{1}
}

func (x *UpdateResponse) GetCreatedIds() []uint64 {
	if x != nil {
		return x.CreatedIds
	}
	return nil
}

type Task struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id        uint64                 `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	Group     string                 `protobuf:"bytes,2,opt,name=group,proto3" json:"group,omitempty"`
	Data      *anypb.Any             `protobuf:"bytes,3,opt,name=data,proto3" json:"data,omitempty"`
	Labels    map[string]string      `protobuf:"bytes,4,rep,name=labels,proto3" json:"labels,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	NotBefore *timestamppb.Timestamp `protobuf:"bytes,5,opt,name=not_before,json=notBefore,proto3" json:"not_before,omitempty"`
}

func (x *Task) Reset() {
	*x = Task{}
	if protoimpl.UnsafeEnabled {
		mi := &file_taskmaster_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Task) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Task) ProtoMessage() {}

func (x *Task) ProtoReflect() protoreflect.Message {
	mi := &file_taskmaster_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Task.ProtoReflect.Descriptor instead.
func (*Task) Descriptor() ([]byte, []int) {
	return file_taskmaster_proto_rawDescGZIP(), []int{2}
}

func (x *Task) GetId() uint64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *Task) GetGroup() string {
	if x != nil {
		return x.Group
	}
	return ""
}

func (x *Task) GetData() *anypb.Any {
	if x != nil {
		return x.Data
	}
	return nil
}

func (x *Task) GetLabels() map[string]string {
	if x != nil {
		return x.Labels
	}
	return nil
}

func (x *Task) GetNotBefore() *timestamppb.Timestamp {
	if x != nil {
		return x.NotBefore
	}
	return nil
}

type TaskRef struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Types that are assignable to Sel:
	//
	//	*TaskRef_Id
	//	*TaskRef_Selector
	Sel isTaskRef_Sel `protobuf_oneof:"sel"`
}

func (x *TaskRef) Reset() {
	*x = TaskRef{}
	if protoimpl.UnsafeEnabled {
		mi := &file_taskmaster_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TaskRef) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TaskRef) ProtoMessage() {}

func (x *TaskRef) ProtoReflect() protoreflect.Message {
	mi := &file_taskmaster_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TaskRef.ProtoReflect.Descriptor instead.
func (*TaskRef) Descriptor() ([]byte, []int) {
	return file_taskmaster_proto_rawDescGZIP(), []int{3}
}

func (m *TaskRef) GetSel() isTaskRef_Sel {
	if m != nil {
		return m.Sel
	}
	return nil
}

func (x *TaskRef) GetId() uint64 {
	if x, ok := x.GetSel().(*TaskRef_Id); ok {
		return x.Id
	}
	return 0
}

func (x *TaskRef) GetSelector() *LabelSelector {
	if x, ok := x.GetSel().(*TaskRef_Selector); ok {
		return x.Selector
	}
	return nil
}

type isTaskRef_Sel interface {
	isTaskRef_Sel()
}

type TaskRef_Id struct {
	Id uint64 `protobuf:"varint,1,opt,name=id,proto3,oneof"`
}

type TaskRef_Selector struct {
	Selector *LabelSelector `protobuf:"bytes,2,opt,name=selector,proto3,oneof"`
}

func (*TaskRef_Id) isTaskRef_Sel() {}

func (*TaskRef_Selector) isTaskRef_Sel() {}

type LabelSelector struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Labels map[string]string `protobuf:"bytes,1,rep,name=labels,proto3" json:"labels,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
}

func (x *LabelSelector) Reset() {
	*x = LabelSelector{}
	if protoimpl.UnsafeEnabled {
		mi := &file_taskmaster_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *LabelSelector) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LabelSelector) ProtoMessage() {}

func (x *LabelSelector) ProtoReflect() protoreflect.Message {
	mi := &file_taskmaster_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LabelSelector.ProtoReflect.Descriptor instead.
func (*LabelSelector) Descriptor() ([]byte, []int) {
	return file_taskmaster_proto_rawDescGZIP(), []int{4}
}

func (x *LabelSelector) GetLabels() map[string]string {
	if x != nil {
		return x.Labels
	}
	return nil
}

type QueryRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Group  string                 `protobuf:"bytes,1,opt,name=group,proto3" json:"group,omitempty"`
	Now    *timestamppb.Timestamp `protobuf:"bytes,4,opt,name=now,proto3" json:"now,omitempty"`
	OwnFor *durationpb.Duration   `protobuf:"bytes,3,opt,name=own_for,json=ownFor,proto3" json:"own_for,omitempty"`
	Wait   bool                   `protobuf:"varint,5,opt,name=wait,proto3" json:"wait,omitempty"`
}

func (x *QueryRequest) Reset() {
	*x = QueryRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_taskmaster_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *QueryRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*QueryRequest) ProtoMessage() {}

func (x *QueryRequest) ProtoReflect() protoreflect.Message {
	mi := &file_taskmaster_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use QueryRequest.ProtoReflect.Descriptor instead.
func (*QueryRequest) Descriptor() ([]byte, []int) {
	return file_taskmaster_proto_rawDescGZIP(), []int{5}
}

func (x *QueryRequest) GetGroup() string {
	if x != nil {
		return x.Group
	}
	return ""
}

func (x *QueryRequest) GetNow() *timestamppb.Timestamp {
	if x != nil {
		return x.Now
	}
	return nil
}

func (x *QueryRequest) GetOwnFor() *durationpb.Duration {
	if x != nil {
		return x.OwnFor
	}
	return nil
}

func (x *QueryRequest) GetWait() bool {
	if x != nil {
		return x.Wait
	}
	return false
}

type QueryResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Task *Task `protobuf:"bytes,1,opt,name=task,proto3" json:"task,omitempty"`
}

func (x *QueryResponse) Reset() {
	*x = QueryResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_taskmaster_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *QueryResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*QueryResponse) ProtoMessage() {}

func (x *QueryResponse) ProtoReflect() protoreflect.Message {
	mi := &file_taskmaster_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use QueryResponse.ProtoReflect.Descriptor instead.
func (*QueryResponse) Descriptor() ([]byte, []int) {
	return file_taskmaster_proto_rawDescGZIP(), []int{6}
}

func (x *QueryResponse) GetTask() *Task {
	if x != nil {
		return x.Task
	}
	return nil
}

type DebugRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *DebugRequest) Reset() {
	*x = DebugRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_taskmaster_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DebugRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DebugRequest) ProtoMessage() {}

func (x *DebugRequest) ProtoReflect() protoreflect.Message {
	mi := &file_taskmaster_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DebugRequest.ProtoReflect.Descriptor instead.
func (*DebugRequest) Descriptor() ([]byte, []int) {
	return file_taskmaster_proto_rawDescGZIP(), []int{7}
}

type DebugResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *DebugResponse) Reset() {
	*x = DebugResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_taskmaster_proto_msgTypes[8]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DebugResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DebugResponse) ProtoMessage() {}

func (x *DebugResponse) ProtoReflect() protoreflect.Message {
	mi := &file_taskmaster_proto_msgTypes[8]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DebugResponse.ProtoReflect.Descriptor instead.
func (*DebugResponse) Descriptor() ([]byte, []int) {
	return file_taskmaster_proto_rawDescGZIP(), []int{8}
}

type Test struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Foo string `protobuf:"bytes,1,opt,name=foo,proto3" json:"foo,omitempty"`
}

func (x *Test) Reset() {
	*x = Test{}
	if protoimpl.UnsafeEnabled {
		mi := &file_taskmaster_proto_msgTypes[9]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Test) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Test) ProtoMessage() {}

func (x *Test) ProtoReflect() protoreflect.Message {
	mi := &file_taskmaster_proto_msgTypes[9]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Test.ProtoReflect.Descriptor instead.
func (*Test) Descriptor() ([]byte, []int) {
	return file_taskmaster_proto_rawDescGZIP(), []int{9}
}

func (x *Test) GetFoo() string {
	if x != nil {
		return x.Foo
	}
	return ""
}

var File_taskmaster_proto protoreflect.FileDescriptor

var file_taskmaster_proto_rawDesc = []byte{
	0x0a, 0x10, 0x74, 0x61, 0x73, 0x6b, 0x6d, 0x61, 0x73, 0x74, 0x65, 0x72, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x12, 0x0a, 0x74, 0x61, 0x73, 0x6b, 0x6d, 0x61, 0x73, 0x74, 0x65, 0x72, 0x1a, 0x1e,
	0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f,
	0x64, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1f,
	0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f,
	0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a,
	0x19, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66,
	0x2f, 0x61, 0x6e, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x8a, 0x01, 0x0a, 0x0d, 0x55,
	0x70, 0x64, 0x61, 0x74, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x2a, 0x0a, 0x07,
	0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x10, 0x2e,
	0x74, 0x61, 0x73, 0x6b, 0x6d, 0x61, 0x73, 0x74, 0x65, 0x72, 0x2e, 0x54, 0x61, 0x73, 0x6b, 0x52,
	0x07, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x12, 0x2d, 0x0a, 0x07, 0x64, 0x65, 0x6c, 0x65,
	0x74, 0x65, 0x64, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x13, 0x2e, 0x74, 0x61, 0x73, 0x6b,
	0x6d, 0x61, 0x73, 0x74, 0x65, 0x72, 0x2e, 0x54, 0x61, 0x73, 0x6b, 0x52, 0x65, 0x66, 0x52, 0x07,
	0x64, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x64, 0x12, 0x1e, 0x0a, 0x0a, 0x70, 0x72, 0x65, 0x64, 0x69,
	0x63, 0x61, 0x74, 0x65, 0x73, 0x18, 0x03, 0x20, 0x03, 0x28, 0x04, 0x52, 0x0a, 0x70, 0x72, 0x65,
	0x64, 0x69, 0x63, 0x61, 0x74, 0x65, 0x73, 0x22, 0x31, 0x0a, 0x0e, 0x55, 0x70, 0x64, 0x61, 0x74,
	0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x1f, 0x0a, 0x0b, 0x63, 0x72, 0x65,
	0x61, 0x74, 0x65, 0x64, 0x5f, 0x69, 0x64, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x04, 0x52, 0x0a,
	0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x49, 0x64, 0x73, 0x22, 0x82, 0x02, 0x0a, 0x04, 0x54,
	0x61, 0x73, 0x6b, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x04, 0x52,
	0x02, 0x69, 0x64, 0x12, 0x14, 0x0a, 0x05, 0x67, 0x72, 0x6f, 0x75, 0x70, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x05, 0x67, 0x72, 0x6f, 0x75, 0x70, 0x12, 0x28, 0x0a, 0x04, 0x64, 0x61, 0x74,
	0x61, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x14, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x41, 0x6e, 0x79, 0x52, 0x04, 0x64,
	0x61, 0x74, 0x61, 0x12, 0x34, 0x0a, 0x06, 0x6c, 0x61, 0x62, 0x65, 0x6c, 0x73, 0x18, 0x04, 0x20,
	0x03, 0x28, 0x0b, 0x32, 0x1c, 0x2e, 0x74, 0x61, 0x73, 0x6b, 0x6d, 0x61, 0x73, 0x74, 0x65, 0x72,
	0x2e, 0x54, 0x61, 0x73, 0x6b, 0x2e, 0x4c, 0x61, 0x62, 0x65, 0x6c, 0x73, 0x45, 0x6e, 0x74, 0x72,
	0x79, 0x52, 0x06, 0x6c, 0x61, 0x62, 0x65, 0x6c, 0x73, 0x12, 0x39, 0x0a, 0x0a, 0x6e, 0x6f, 0x74,
	0x5f, 0x62, 0x65, 0x66, 0x6f, 0x72, 0x65, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e,
	0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e,
	0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x09, 0x6e, 0x6f, 0x74, 0x42, 0x65,
	0x66, 0x6f, 0x72, 0x65, 0x1a, 0x39, 0x0a, 0x0b, 0x4c, 0x61, 0x62, 0x65, 0x6c, 0x73, 0x45, 0x6e,
	0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38, 0x01, 0x22,
	0x5b, 0x0a, 0x07, 0x54, 0x61, 0x73, 0x6b, 0x52, 0x65, 0x66, 0x12, 0x10, 0x0a, 0x02, 0x69, 0x64,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x04, 0x48, 0x00, 0x52, 0x02, 0x69, 0x64, 0x12, 0x37, 0x0a, 0x08,
	0x73, 0x65, 0x6c, 0x65, 0x63, 0x74, 0x6f, 0x72, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x19,
	0x2e, 0x74, 0x61, 0x73, 0x6b, 0x6d, 0x61, 0x73, 0x74, 0x65, 0x72, 0x2e, 0x4c, 0x61, 0x62, 0x65,
	0x6c, 0x53, 0x65, 0x6c, 0x65, 0x63, 0x74, 0x6f, 0x72, 0x48, 0x00, 0x52, 0x08, 0x73, 0x65, 0x6c,
	0x65, 0x63, 0x74, 0x6f, 0x72, 0x42, 0x05, 0x0a, 0x03, 0x73, 0x65, 0x6c, 0x22, 0x89, 0x01, 0x0a,
	0x0d, 0x4c, 0x61, 0x62, 0x65, 0x6c, 0x53, 0x65, 0x6c, 0x65, 0x63, 0x74, 0x6f, 0x72, 0x12, 0x3d,
	0x0a, 0x06, 0x6c, 0x61, 0x62, 0x65, 0x6c, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x25,
	0x2e, 0x74, 0x61, 0x73, 0x6b, 0x6d, 0x61, 0x73, 0x74, 0x65, 0x72, 0x2e, 0x4c, 0x61, 0x62, 0x65,
	0x6c, 0x53, 0x65, 0x6c, 0x65, 0x63, 0x74, 0x6f, 0x72, 0x2e, 0x4c, 0x61, 0x62, 0x65, 0x6c, 0x73,
	0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x06, 0x6c, 0x61, 0x62, 0x65, 0x6c, 0x73, 0x1a, 0x39, 0x0a,
	0x0b, 0x4c, 0x61, 0x62, 0x65, 0x6c, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03,
	0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x14,
	0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x76,
	0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38, 0x01, 0x22, 0x9a, 0x01, 0x0a, 0x0c, 0x51, 0x75, 0x65,
	0x72, 0x79, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x14, 0x0a, 0x05, 0x67, 0x72, 0x6f,
	0x75, 0x70, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x67, 0x72, 0x6f, 0x75, 0x70, 0x12,
	0x2c, 0x0a, 0x03, 0x6e, 0x6f, 0x77, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67,
	0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54,
	0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x03, 0x6e, 0x6f, 0x77, 0x12, 0x32, 0x0a,
	0x07, 0x6f, 0x77, 0x6e, 0x5f, 0x66, 0x6f, 0x72, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x19,
	0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66,
	0x2e, 0x44, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x06, 0x6f, 0x77, 0x6e, 0x46, 0x6f,
	0x72, 0x12, 0x12, 0x0a, 0x04, 0x77, 0x61, 0x69, 0x74, 0x18, 0x05, 0x20, 0x01, 0x28, 0x08, 0x52,
	0x04, 0x77, 0x61, 0x69, 0x74, 0x22, 0x35, 0x0a, 0x0d, 0x51, 0x75, 0x65, 0x72, 0x79, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x24, 0x0a, 0x04, 0x74, 0x61, 0x73, 0x6b, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x10, 0x2e, 0x74, 0x61, 0x73, 0x6b, 0x6d, 0x61, 0x73, 0x74, 0x65,
	0x72, 0x2e, 0x54, 0x61, 0x73, 0x6b, 0x52, 0x04, 0x74, 0x61, 0x73, 0x6b, 0x22, 0x0e, 0x0a, 0x0c,
	0x44, 0x65, 0x62, 0x75, 0x67, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x22, 0x0f, 0x0a, 0x0d,
	0x44, 0x65, 0x62, 0x75, 0x67, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x18, 0x0a,
	0x04, 0x54, 0x65, 0x73, 0x74, 0x12, 0x10, 0x0a, 0x03, 0x66, 0x6f, 0x6f, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x03, 0x66, 0x6f, 0x6f, 0x32, 0xc9, 0x01, 0x0a, 0x0a, 0x54, 0x61, 0x73, 0x6b,
	0x6d, 0x61, 0x73, 0x74, 0x65, 0x72, 0x12, 0x3f, 0x0a, 0x06, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65,
	0x12, 0x19, 0x2e, 0x74, 0x61, 0x73, 0x6b, 0x6d, 0x61, 0x73, 0x74, 0x65, 0x72, 0x2e, 0x55, 0x70,
	0x64, 0x61, 0x74, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1a, 0x2e, 0x74, 0x61,
	0x73, 0x6b, 0x6d, 0x61, 0x73, 0x74, 0x65, 0x72, 0x2e, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x3c, 0x0a, 0x05, 0x51, 0x75, 0x65, 0x72, 0x79,
	0x12, 0x18, 0x2e, 0x74, 0x61, 0x73, 0x6b, 0x6d, 0x61, 0x73, 0x74, 0x65, 0x72, 0x2e, 0x51, 0x75,
	0x65, 0x72, 0x79, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x19, 0x2e, 0x74, 0x61, 0x73,
	0x6b, 0x6d, 0x61, 0x73, 0x74, 0x65, 0x72, 0x2e, 0x51, 0x75, 0x65, 0x72, 0x79, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x3c, 0x0a, 0x05, 0x44, 0x65, 0x62, 0x75, 0x67, 0x12, 0x18,
	0x2e, 0x74, 0x61, 0x73, 0x6b, 0x6d, 0x61, 0x73, 0x74, 0x65, 0x72, 0x2e, 0x44, 0x65, 0x62, 0x75,
	0x67, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x19, 0x2e, 0x74, 0x61, 0x73, 0x6b, 0x6d,
	0x61, 0x73, 0x74, 0x65, 0x72, 0x2e, 0x44, 0x65, 0x62, 0x75, 0x67, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x42, 0x23, 0x5a, 0x21, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f,
	0x6d, 0x2f, 0x6d, 0x6b, 0x6d, 0x69, 0x6b, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x3b, 0x74, 0x61,
	0x73, 0x6b, 0x6d, 0x61, 0x73, 0x74, 0x65, 0x72, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_taskmaster_proto_rawDescOnce sync.Once
	file_taskmaster_proto_rawDescData = file_taskmaster_proto_rawDesc
)

func file_taskmaster_proto_rawDescGZIP() []byte {
	file_taskmaster_proto_rawDescOnce.Do(func() {
		file_taskmaster_proto_rawDescData = protoimpl.X.CompressGZIP(file_taskmaster_proto_rawDescData)
	})
	return file_taskmaster_proto_rawDescData
}

var file_taskmaster_proto_msgTypes = make([]protoimpl.MessageInfo, 12)
var file_taskmaster_proto_goTypes = []interface{}{
	(*UpdateRequest)(nil),         // 0: taskmaster.UpdateRequest
	(*UpdateResponse)(nil),        // 1: taskmaster.UpdateResponse
	(*Task)(nil),                  // 2: taskmaster.Task
	(*TaskRef)(nil),               // 3: taskmaster.TaskRef
	(*LabelSelector)(nil),         // 4: taskmaster.LabelSelector
	(*QueryRequest)(nil),          // 5: taskmaster.QueryRequest
	(*QueryResponse)(nil),         // 6: taskmaster.QueryResponse
	(*DebugRequest)(nil),          // 7: taskmaster.DebugRequest
	(*DebugResponse)(nil),         // 8: taskmaster.DebugResponse
	(*Test)(nil),                  // 9: taskmaster.Test
	nil,                           // 10: taskmaster.Task.LabelsEntry
	nil,                           // 11: taskmaster.LabelSelector.LabelsEntry
	(*anypb.Any)(nil),             // 12: google.protobuf.Any
	(*timestamppb.Timestamp)(nil), // 13: google.protobuf.Timestamp
	(*durationpb.Duration)(nil),   // 14: google.protobuf.Duration
}
var file_taskmaster_proto_depIdxs = []int32{
	2,  // 0: taskmaster.UpdateRequest.created:type_name -> taskmaster.Task
	3,  // 1: taskmaster.UpdateRequest.deleted:type_name -> taskmaster.TaskRef
	12, // 2: taskmaster.Task.data:type_name -> google.protobuf.Any
	10, // 3: taskmaster.Task.labels:type_name -> taskmaster.Task.LabelsEntry
	13, // 4: taskmaster.Task.not_before:type_name -> google.protobuf.Timestamp
	4,  // 5: taskmaster.TaskRef.selector:type_name -> taskmaster.LabelSelector
	11, // 6: taskmaster.LabelSelector.labels:type_name -> taskmaster.LabelSelector.LabelsEntry
	13, // 7: taskmaster.QueryRequest.now:type_name -> google.protobuf.Timestamp
	14, // 8: taskmaster.QueryRequest.own_for:type_name -> google.protobuf.Duration
	2,  // 9: taskmaster.QueryResponse.task:type_name -> taskmaster.Task
	0,  // 10: taskmaster.Taskmaster.Update:input_type -> taskmaster.UpdateRequest
	5,  // 11: taskmaster.Taskmaster.Query:input_type -> taskmaster.QueryRequest
	7,  // 12: taskmaster.Taskmaster.Debug:input_type -> taskmaster.DebugRequest
	1,  // 13: taskmaster.Taskmaster.Update:output_type -> taskmaster.UpdateResponse
	6,  // 14: taskmaster.Taskmaster.Query:output_type -> taskmaster.QueryResponse
	8,  // 15: taskmaster.Taskmaster.Debug:output_type -> taskmaster.DebugResponse
	13, // [13:16] is the sub-list for method output_type
	10, // [10:13] is the sub-list for method input_type
	10, // [10:10] is the sub-list for extension type_name
	10, // [10:10] is the sub-list for extension extendee
	0,  // [0:10] is the sub-list for field type_name
}

func init() { file_taskmaster_proto_init() }
func file_taskmaster_proto_init() {
	if File_taskmaster_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_taskmaster_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UpdateRequest); i {
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
		file_taskmaster_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UpdateResponse); i {
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
		file_taskmaster_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Task); i {
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
		file_taskmaster_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*TaskRef); i {
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
		file_taskmaster_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*LabelSelector); i {
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
		file_taskmaster_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*QueryRequest); i {
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
		file_taskmaster_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*QueryResponse); i {
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
		file_taskmaster_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DebugRequest); i {
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
		file_taskmaster_proto_msgTypes[8].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DebugResponse); i {
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
		file_taskmaster_proto_msgTypes[9].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Test); i {
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
	file_taskmaster_proto_msgTypes[3].OneofWrappers = []interface{}{
		(*TaskRef_Id)(nil),
		(*TaskRef_Selector)(nil),
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_taskmaster_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   12,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_taskmaster_proto_goTypes,
		DependencyIndexes: file_taskmaster_proto_depIdxs,
		MessageInfos:      file_taskmaster_proto_msgTypes,
	}.Build()
	File_taskmaster_proto = out.File
	file_taskmaster_proto_rawDesc = nil
	file_taskmaster_proto_goTypes = nil
	file_taskmaster_proto_depIdxs = nil
}
