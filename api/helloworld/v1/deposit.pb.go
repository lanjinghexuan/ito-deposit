// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.6
// 	protoc        v5.26.1
// source: helloworld/v1/deposit.proto

package v1

import (
	_ "google.golang.org/genproto/googleapis/api/annotations"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
	unsafe "unsafe"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type SendCodeByOrderReq struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	OrderNo       string                 `protobuf:"bytes,1,opt,name=order_no,json=orderNo,proto3" json:"order_no,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *SendCodeByOrderReq) Reset() {
	*x = SendCodeByOrderReq{}
	mi := &file_helloworld_v1_deposit_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *SendCodeByOrderReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SendCodeByOrderReq) ProtoMessage() {}

func (x *SendCodeByOrderReq) ProtoReflect() protoreflect.Message {
	mi := &file_helloworld_v1_deposit_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SendCodeByOrderReq.ProtoReflect.Descriptor instead.
func (*SendCodeByOrderReq) Descriptor() ([]byte, []int) {
	return file_helloworld_v1_deposit_proto_rawDescGZIP(), []int{0}
}

func (x *SendCodeByOrderReq) GetOrderNo() string {
	if x != nil {
		return x.OrderNo
	}
	return ""
}

type SendCodeByOrderRes struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Msg           string                 `protobuf:"bytes,1,opt,name=msg,proto3" json:"msg,omitempty"`
	Code          int32                  `protobuf:"varint,2,opt,name=code,proto3" json:"code,omitempty"`
	Data          string                 `protobuf:"bytes,3,opt,name=data,proto3" json:"data,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *SendCodeByOrderRes) Reset() {
	*x = SendCodeByOrderRes{}
	mi := &file_helloworld_v1_deposit_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *SendCodeByOrderRes) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SendCodeByOrderRes) ProtoMessage() {}

func (x *SendCodeByOrderRes) ProtoReflect() protoreflect.Message {
	mi := &file_helloworld_v1_deposit_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SendCodeByOrderRes.ProtoReflect.Descriptor instead.
func (*SendCodeByOrderRes) Descriptor() ([]byte, []int) {
	return file_helloworld_v1_deposit_proto_rawDescGZIP(), []int{1}
}

func (x *SendCodeByOrderRes) GetMsg() string {
	if x != nil {
		return x.Msg
	}
	return ""
}

func (x *SendCodeByOrderRes) GetCode() int32 {
	if x != nil {
		return x.Code
	}
	return 0
}

func (x *SendCodeByOrderRes) GetData() string {
	if x != nil {
		return x.Data
	}
	return ""
}

type UpdateDepositLockerIdReq struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	OrderId       string                 `protobuf:"bytes,1,opt,name=order_id,json=orderId,proto3" json:"order_id,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *UpdateDepositLockerIdReq) Reset() {
	*x = UpdateDepositLockerIdReq{}
	mi := &file_helloworld_v1_deposit_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *UpdateDepositLockerIdReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UpdateDepositLockerIdReq) ProtoMessage() {}

func (x *UpdateDepositLockerIdReq) ProtoReflect() protoreflect.Message {
	mi := &file_helloworld_v1_deposit_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UpdateDepositLockerIdReq.ProtoReflect.Descriptor instead.
func (*UpdateDepositLockerIdReq) Descriptor() ([]byte, []int) {
	return file_helloworld_v1_deposit_proto_rawDescGZIP(), []int{2}
}

func (x *UpdateDepositLockerIdReq) GetOrderId() string {
	if x != nil {
		return x.OrderId
	}
	return ""
}

type UpdateDepositLockerIdRes struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Code          int32                  `protobuf:"varint,2,opt,name=code,proto3" json:"code,omitempty"`
	Msg           string                 `protobuf:"bytes,3,opt,name=msg,proto3" json:"msg,omitempty"`
	LockerId      int32                  `protobuf:"varint,1,opt,name=locker_id,json=lockerId,proto3" json:"locker_id,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *UpdateDepositLockerIdRes) Reset() {
	*x = UpdateDepositLockerIdRes{}
	mi := &file_helloworld_v1_deposit_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *UpdateDepositLockerIdRes) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UpdateDepositLockerIdRes) ProtoMessage() {}

func (x *UpdateDepositLockerIdRes) ProtoReflect() protoreflect.Message {
	mi := &file_helloworld_v1_deposit_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UpdateDepositLockerIdRes.ProtoReflect.Descriptor instead.
func (*UpdateDepositLockerIdRes) Descriptor() ([]byte, []int) {
	return file_helloworld_v1_deposit_proto_rawDescGZIP(), []int{3}
}

func (x *UpdateDepositLockerIdRes) GetCode() int32 {
	if x != nil {
		return x.Code
	}
	return 0
}

func (x *UpdateDepositLockerIdRes) GetMsg() string {
	if x != nil {
		return x.Msg
	}
	return ""
}

func (x *UpdateDepositLockerIdRes) GetLockerId() int32 {
	if x != nil {
		return x.LockerId
	}
	return 0
}

type GetDepositLockerReq struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	LockerId      int32                  `protobuf:"varint,1,opt,name=locker_id,json=lockerId,proto3" json:"locker_id,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *GetDepositLockerReq) Reset() {
	*x = GetDepositLockerReq{}
	mi := &file_helloworld_v1_deposit_proto_msgTypes[4]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetDepositLockerReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetDepositLockerReq) ProtoMessage() {}

func (x *GetDepositLockerReq) ProtoReflect() protoreflect.Message {
	mi := &file_helloworld_v1_deposit_proto_msgTypes[4]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetDepositLockerReq.ProtoReflect.Descriptor instead.
func (*GetDepositLockerReq) Descriptor() ([]byte, []int) {
	return file_helloworld_v1_deposit_proto_rawDescGZIP(), []int{4}
}

func (x *GetDepositLockerReq) GetLockerId() int32 {
	if x != nil {
		return x.LockerId
	}
	return 0
}

type GetDepositLockerRes struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Address       string                 `protobuf:"bytes,1,opt,name=address,proto3" json:"address,omitempty"`
	Name          string                 `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	Longitude     float32                `protobuf:"fixed32,3,opt,name=longitude,proto3" json:"longitude,omitempty"`
	Latitude      float32                `protobuf:"fixed32,4,opt,name=latitude,proto3" json:"latitude,omitempty"`
	Locker        []*Locker              `protobuf:"bytes,5,rep,name=locker,proto3" json:"locker,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *GetDepositLockerRes) Reset() {
	*x = GetDepositLockerRes{}
	mi := &file_helloworld_v1_deposit_proto_msgTypes[5]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetDepositLockerRes) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetDepositLockerRes) ProtoMessage() {}

func (x *GetDepositLockerRes) ProtoReflect() protoreflect.Message {
	mi := &file_helloworld_v1_deposit_proto_msgTypes[5]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetDepositLockerRes.ProtoReflect.Descriptor instead.
func (*GetDepositLockerRes) Descriptor() ([]byte, []int) {
	return file_helloworld_v1_deposit_proto_rawDescGZIP(), []int{5}
}

func (x *GetDepositLockerRes) GetAddress() string {
	if x != nil {
		return x.Address
	}
	return ""
}

func (x *GetDepositLockerRes) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *GetDepositLockerRes) GetLongitude() float32 {
	if x != nil {
		return x.Longitude
	}
	return 0
}

func (x *GetDepositLockerRes) GetLatitude() float32 {
	if x != nil {
		return x.Latitude
	}
	return 0
}

func (x *GetDepositLockerRes) GetLocker() []*Locker {
	if x != nil {
		return x.Locker
	}
	return nil
}

type Locker struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Name          string                 `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Description   string                 `protobuf:"bytes,2,opt,name=description,proto3" json:"description,omitempty"`
	Size          string                 `protobuf:"bytes,3,opt,name=size,proto3" json:"size,omitempty"`
	Num           int32                  `protobuf:"varint,4,opt,name=num,proto3" json:"num,omitempty"`
	HourlyRate    float32                `protobuf:"fixed32,5,opt,name=hourly_rate,json=hourlyRate,proto3" json:"hourly_rate,omitempty"`
	LockerType    int32                  `protobuf:"varint,6,opt,name=locker_type,json=lockerType,proto3" json:"locker_type,omitempty"`
	FreeDuration  float32                `protobuf:"fixed32,7,opt,name=free_duration,json=freeDuration,proto3" json:"free_duration,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Locker) Reset() {
	*x = Locker{}
	mi := &file_helloworld_v1_deposit_proto_msgTypes[6]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Locker) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Locker) ProtoMessage() {}

func (x *Locker) ProtoReflect() protoreflect.Message {
	mi := &file_helloworld_v1_deposit_proto_msgTypes[6]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Locker.ProtoReflect.Descriptor instead.
func (*Locker) Descriptor() ([]byte, []int) {
	return file_helloworld_v1_deposit_proto_rawDescGZIP(), []int{6}
}

func (x *Locker) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *Locker) GetDescription() string {
	if x != nil {
		return x.Description
	}
	return ""
}

func (x *Locker) GetSize() string {
	if x != nil {
		return x.Size
	}
	return ""
}

func (x *Locker) GetNum() int32 {
	if x != nil {
		return x.Num
	}
	return 0
}

func (x *Locker) GetHourlyRate() float32 {
	if x != nil {
		return x.HourlyRate
	}
	return 0
}

func (x *Locker) GetLockerType() int32 {
	if x != nil {
		return x.LockerType
	}
	return 0
}

func (x *Locker) GetFreeDuration() float32 {
	if x != nil {
		return x.FreeDuration
	}
	return 0
}

type CreateDepositRequest struct {
	state             protoimpl.MessageState `protogen:"open.v1"`
	ScheduledDuration int32                  `protobuf:"varint,1,opt,name=scheduled_duration,json=scheduledDuration,proto3" json:"scheduled_duration,omitempty"` //预计存储时间
	LockerType        int32                  `protobuf:"varint,2,opt,name=locker_type,json=lockerType,proto3" json:"locker_type,omitempty"`                      //柜子类型
	CabinetId         int32                  `protobuf:"varint,3,opt,name=cabinet_id,json=cabinetId,proto3" json:"cabinet_id,omitempty"`                         //网点id
	unknownFields     protoimpl.UnknownFields
	sizeCache         protoimpl.SizeCache
}

func (x *CreateDepositRequest) Reset() {
	*x = CreateDepositRequest{}
	mi := &file_helloworld_v1_deposit_proto_msgTypes[7]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *CreateDepositRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateDepositRequest) ProtoMessage() {}

func (x *CreateDepositRequest) ProtoReflect() protoreflect.Message {
	mi := &file_helloworld_v1_deposit_proto_msgTypes[7]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateDepositRequest.ProtoReflect.Descriptor instead.
func (*CreateDepositRequest) Descriptor() ([]byte, []int) {
	return file_helloworld_v1_deposit_proto_rawDescGZIP(), []int{7}
}

func (x *CreateDepositRequest) GetScheduledDuration() int32 {
	if x != nil {
		return x.ScheduledDuration
	}
	return 0
}

func (x *CreateDepositRequest) GetLockerType() int32 {
	if x != nil {
		return x.LockerType
	}
	return 0
}

func (x *CreateDepositRequest) GetCabinetId() int32 {
	if x != nil {
		return x.CabinetId
	}
	return 0
}

type CreateDepositReply struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Code          int32                  `protobuf:"varint,1,opt,name=code,proto3" json:"code,omitempty"`
	Msg           string                 `protobuf:"bytes,2,opt,name=msg,proto3" json:"msg,omitempty"`
	Data          *DepositReplyData      `protobuf:"bytes,3,opt,name=data,proto3" json:"data,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *CreateDepositReply) Reset() {
	*x = CreateDepositReply{}
	mi := &file_helloworld_v1_deposit_proto_msgTypes[8]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *CreateDepositReply) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateDepositReply) ProtoMessage() {}

func (x *CreateDepositReply) ProtoReflect() protoreflect.Message {
	mi := &file_helloworld_v1_deposit_proto_msgTypes[8]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateDepositReply.ProtoReflect.Descriptor instead.
func (*CreateDepositReply) Descriptor() ([]byte, []int) {
	return file_helloworld_v1_deposit_proto_rawDescGZIP(), []int{8}
}

func (x *CreateDepositReply) GetCode() int32 {
	if x != nil {
		return x.Code
	}
	return 0
}

func (x *CreateDepositReply) GetMsg() string {
	if x != nil {
		return x.Msg
	}
	return ""
}

func (x *CreateDepositReply) GetData() *DepositReplyData {
	if x != nil {
		return x.Data
	}
	return nil
}

type DepositReplyData struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	OrderNo       string                 `protobuf:"bytes,1,opt,name=order_no,json=orderNo,proto3" json:"order_no,omitempty"`
	LockerId      int32                  `protobuf:"varint,2,opt,name=locker_id,json=lockerId,proto3" json:"locker_id,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *DepositReplyData) Reset() {
	*x = DepositReplyData{}
	mi := &file_helloworld_v1_deposit_proto_msgTypes[9]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *DepositReplyData) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DepositReplyData) ProtoMessage() {}

func (x *DepositReplyData) ProtoReflect() protoreflect.Message {
	mi := &file_helloworld_v1_deposit_proto_msgTypes[9]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DepositReplyData.ProtoReflect.Descriptor instead.
func (*DepositReplyData) Descriptor() ([]byte, []int) {
	return file_helloworld_v1_deposit_proto_rawDescGZIP(), []int{9}
}

func (x *DepositReplyData) GetOrderNo() string {
	if x != nil {
		return x.OrderNo
	}
	return ""
}

func (x *DepositReplyData) GetLockerId() int32 {
	if x != nil {
		return x.LockerId
	}
	return 0
}

type UpdateDepositRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *UpdateDepositRequest) Reset() {
	*x = UpdateDepositRequest{}
	mi := &file_helloworld_v1_deposit_proto_msgTypes[10]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *UpdateDepositRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UpdateDepositRequest) ProtoMessage() {}

func (x *UpdateDepositRequest) ProtoReflect() protoreflect.Message {
	mi := &file_helloworld_v1_deposit_proto_msgTypes[10]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UpdateDepositRequest.ProtoReflect.Descriptor instead.
func (*UpdateDepositRequest) Descriptor() ([]byte, []int) {
	return file_helloworld_v1_deposit_proto_rawDescGZIP(), []int{10}
}

type UpdateDepositReply struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *UpdateDepositReply) Reset() {
	*x = UpdateDepositReply{}
	mi := &file_helloworld_v1_deposit_proto_msgTypes[11]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *UpdateDepositReply) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UpdateDepositReply) ProtoMessage() {}

func (x *UpdateDepositReply) ProtoReflect() protoreflect.Message {
	mi := &file_helloworld_v1_deposit_proto_msgTypes[11]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UpdateDepositReply.ProtoReflect.Descriptor instead.
func (*UpdateDepositReply) Descriptor() ([]byte, []int) {
	return file_helloworld_v1_deposit_proto_rawDescGZIP(), []int{11}
}

type DeleteDepositRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *DeleteDepositRequest) Reset() {
	*x = DeleteDepositRequest{}
	mi := &file_helloworld_v1_deposit_proto_msgTypes[12]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *DeleteDepositRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeleteDepositRequest) ProtoMessage() {}

func (x *DeleteDepositRequest) ProtoReflect() protoreflect.Message {
	mi := &file_helloworld_v1_deposit_proto_msgTypes[12]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeleteDepositRequest.ProtoReflect.Descriptor instead.
func (*DeleteDepositRequest) Descriptor() ([]byte, []int) {
	return file_helloworld_v1_deposit_proto_rawDescGZIP(), []int{12}
}

type DeleteDepositReply struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *DeleteDepositReply) Reset() {
	*x = DeleteDepositReply{}
	mi := &file_helloworld_v1_deposit_proto_msgTypes[13]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *DeleteDepositReply) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeleteDepositReply) ProtoMessage() {}

func (x *DeleteDepositReply) ProtoReflect() protoreflect.Message {
	mi := &file_helloworld_v1_deposit_proto_msgTypes[13]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeleteDepositReply.ProtoReflect.Descriptor instead.
func (*DeleteDepositReply) Descriptor() ([]byte, []int) {
	return file_helloworld_v1_deposit_proto_rawDescGZIP(), []int{13}
}

type GetDepositRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *GetDepositRequest) Reset() {
	*x = GetDepositRequest{}
	mi := &file_helloworld_v1_deposit_proto_msgTypes[14]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetDepositRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetDepositRequest) ProtoMessage() {}

func (x *GetDepositRequest) ProtoReflect() protoreflect.Message {
	mi := &file_helloworld_v1_deposit_proto_msgTypes[14]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetDepositRequest.ProtoReflect.Descriptor instead.
func (*GetDepositRequest) Descriptor() ([]byte, []int) {
	return file_helloworld_v1_deposit_proto_rawDescGZIP(), []int{14}
}

type GetDepositReply struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *GetDepositReply) Reset() {
	*x = GetDepositReply{}
	mi := &file_helloworld_v1_deposit_proto_msgTypes[15]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetDepositReply) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetDepositReply) ProtoMessage() {}

func (x *GetDepositReply) ProtoReflect() protoreflect.Message {
	mi := &file_helloworld_v1_deposit_proto_msgTypes[15]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetDepositReply.ProtoReflect.Descriptor instead.
func (*GetDepositReply) Descriptor() ([]byte, []int) {
	return file_helloworld_v1_deposit_proto_rawDescGZIP(), []int{15}
}

type ListDepositRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ListDepositRequest) Reset() {
	*x = ListDepositRequest{}
	mi := &file_helloworld_v1_deposit_proto_msgTypes[16]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ListDepositRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListDepositRequest) ProtoMessage() {}

func (x *ListDepositRequest) ProtoReflect() protoreflect.Message {
	mi := &file_helloworld_v1_deposit_proto_msgTypes[16]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListDepositRequest.ProtoReflect.Descriptor instead.
func (*ListDepositRequest) Descriptor() ([]byte, []int) {
	return file_helloworld_v1_deposit_proto_rawDescGZIP(), []int{16}
}

type ListDepositReply struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ListDepositReply) Reset() {
	*x = ListDepositReply{}
	mi := &file_helloworld_v1_deposit_proto_msgTypes[17]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ListDepositReply) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListDepositReply) ProtoMessage() {}

func (x *ListDepositReply) ProtoReflect() protoreflect.Message {
	mi := &file_helloworld_v1_deposit_proto_msgTypes[17]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListDepositReply.ProtoReflect.Descriptor instead.
func (*ListDepositReply) Descriptor() ([]byte, []int) {
	return file_helloworld_v1_deposit_proto_rawDescGZIP(), []int{17}
}

type ReturnTokenReq struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ReturnTokenReq) Reset() {
	*x = ReturnTokenReq{}
	mi := &file_helloworld_v1_deposit_proto_msgTypes[18]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ReturnTokenReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ReturnTokenReq) ProtoMessage() {}

func (x *ReturnTokenReq) ProtoReflect() protoreflect.Message {
	mi := &file_helloworld_v1_deposit_proto_msgTypes[18]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ReturnTokenReq.ProtoReflect.Descriptor instead.
func (*ReturnTokenReq) Descriptor() ([]byte, []int) {
	return file_helloworld_v1_deposit_proto_rawDescGZIP(), []int{18}
}

type ReturnTokenRes struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Token         string                 `protobuf:"bytes,1,opt,name=token,proto3" json:"token,omitempty"`
	Coe           int32                  `protobuf:"varint,2,opt,name=coe,proto3" json:"coe,omitempty"`
	Msg           string                 `protobuf:"bytes,3,opt,name=msg,proto3" json:"msg,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ReturnTokenRes) Reset() {
	*x = ReturnTokenRes{}
	mi := &file_helloworld_v1_deposit_proto_msgTypes[19]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ReturnTokenRes) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ReturnTokenRes) ProtoMessage() {}

func (x *ReturnTokenRes) ProtoReflect() protoreflect.Message {
	mi := &file_helloworld_v1_deposit_proto_msgTypes[19]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ReturnTokenRes.ProtoReflect.Descriptor instead.
func (*ReturnTokenRes) Descriptor() ([]byte, []int) {
	return file_helloworld_v1_deposit_proto_rawDescGZIP(), []int{19}
}

func (x *ReturnTokenRes) GetToken() string {
	if x != nil {
		return x.Token
	}
	return ""
}

func (x *ReturnTokenRes) GetCoe() int32 {
	if x != nil {
		return x.Coe
	}
	return 0
}

func (x *ReturnTokenRes) GetMsg() string {
	if x != nil {
		return x.Msg
	}
	return ""
}

type DecodeTokenRes struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Data          string                 `protobuf:"bytes,1,opt,name=data,proto3" json:"data,omitempty"`
	Coe           int32                  `protobuf:"varint,2,opt,name=coe,proto3" json:"coe,omitempty"`
	Msg           string                 `protobuf:"bytes,3,opt,name=msg,proto3" json:"msg,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *DecodeTokenRes) Reset() {
	*x = DecodeTokenRes{}
	mi := &file_helloworld_v1_deposit_proto_msgTypes[20]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *DecodeTokenRes) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DecodeTokenRes) ProtoMessage() {}

func (x *DecodeTokenRes) ProtoReflect() protoreflect.Message {
	mi := &file_helloworld_v1_deposit_proto_msgTypes[20]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DecodeTokenRes.ProtoReflect.Descriptor instead.
func (*DecodeTokenRes) Descriptor() ([]byte, []int) {
	return file_helloworld_v1_deposit_proto_rawDescGZIP(), []int{20}
}

func (x *DecodeTokenRes) GetData() string {
	if x != nil {
		return x.Data
	}
	return ""
}

func (x *DecodeTokenRes) GetCoe() int32 {
	if x != nil {
		return x.Coe
	}
	return 0
}

func (x *DecodeTokenRes) GetMsg() string {
	if x != nil {
		return x.Msg
	}
	return ""
}

var File_helloworld_v1_deposit_proto protoreflect.FileDescriptor

const file_helloworld_v1_deposit_proto_rawDesc = "" +
	"\n" +
	"\x1bhelloworld/v1/deposit.proto\x12\x11api.helloworld.v1\x1a\x1cgoogle/api/annotations.proto\"/\n" +
	"\x12SendCodeByOrderReq\x12\x19\n" +
	"\border_no\x18\x01 \x01(\tR\aorderNo\"N\n" +
	"\x12SendCodeByOrderRes\x12\x10\n" +
	"\x03msg\x18\x01 \x01(\tR\x03msg\x12\x12\n" +
	"\x04code\x18\x02 \x01(\x05R\x04code\x12\x12\n" +
	"\x04data\x18\x03 \x01(\tR\x04data\"5\n" +
	"\x18UpdateDepositLockerIdReq\x12\x19\n" +
	"\border_id\x18\x01 \x01(\tR\aorderId\"]\n" +
	"\x18UpdateDepositLockerIdRes\x12\x12\n" +
	"\x04code\x18\x02 \x01(\x05R\x04code\x12\x10\n" +
	"\x03msg\x18\x03 \x01(\tR\x03msg\x12\x1b\n" +
	"\tlocker_id\x18\x01 \x01(\x05R\blockerId\"2\n" +
	"\x13GetDepositLockerReq\x12\x1b\n" +
	"\tlocker_id\x18\x01 \x01(\x05R\blockerId\"\xb0\x01\n" +
	"\x13GetDepositLockerRes\x12\x18\n" +
	"\aaddress\x18\x01 \x01(\tR\aaddress\x12\x12\n" +
	"\x04name\x18\x02 \x01(\tR\x04name\x12\x1c\n" +
	"\tlongitude\x18\x03 \x01(\x02R\tlongitude\x12\x1a\n" +
	"\blatitude\x18\x04 \x01(\x02R\blatitude\x121\n" +
	"\x06locker\x18\x05 \x03(\v2\x19.api.helloworld.v1.LockerR\x06locker\"\xcb\x01\n" +
	"\x06Locker\x12\x12\n" +
	"\x04name\x18\x01 \x01(\tR\x04name\x12 \n" +
	"\vdescription\x18\x02 \x01(\tR\vdescription\x12\x12\n" +
	"\x04size\x18\x03 \x01(\tR\x04size\x12\x10\n" +
	"\x03num\x18\x04 \x01(\x05R\x03num\x12\x1f\n" +
	"\vhourly_rate\x18\x05 \x01(\x02R\n" +
	"hourlyRate\x12\x1f\n" +
	"\vlocker_type\x18\x06 \x01(\x05R\n" +
	"lockerType\x12#\n" +
	"\rfree_duration\x18\a \x01(\x02R\ffreeDuration\"\x85\x01\n" +
	"\x14CreateDepositRequest\x12-\n" +
	"\x12scheduled_duration\x18\x01 \x01(\x05R\x11scheduledDuration\x12\x1f\n" +
	"\vlocker_type\x18\x02 \x01(\x05R\n" +
	"lockerType\x12\x1d\n" +
	"\n" +
	"cabinet_id\x18\x03 \x01(\x05R\tcabinetId\"s\n" +
	"\x12CreateDepositReply\x12\x12\n" +
	"\x04code\x18\x01 \x01(\x05R\x04code\x12\x10\n" +
	"\x03msg\x18\x02 \x01(\tR\x03msg\x127\n" +
	"\x04data\x18\x03 \x01(\v2#.api.helloworld.v1.DepositReplyDataR\x04data\"J\n" +
	"\x10DepositReplyData\x12\x19\n" +
	"\border_no\x18\x01 \x01(\tR\aorderNo\x12\x1b\n" +
	"\tlocker_id\x18\x02 \x01(\x05R\blockerId\"\x16\n" +
	"\x14UpdateDepositRequest\"\x14\n" +
	"\x12UpdateDepositReply\"\x16\n" +
	"\x14DeleteDepositRequest\"\x14\n" +
	"\x12DeleteDepositReply\"\x13\n" +
	"\x11GetDepositRequest\"\x11\n" +
	"\x0fGetDepositReply\"\x14\n" +
	"\x12ListDepositRequest\"\x12\n" +
	"\x10ListDepositReply\"\x10\n" +
	"\x0eReturnTokenReq\"J\n" +
	"\x0eReturnTokenRes\x12\x14\n" +
	"\x05token\x18\x01 \x01(\tR\x05token\x12\x10\n" +
	"\x03coe\x18\x02 \x01(\x05R\x03coe\x12\x10\n" +
	"\x03msg\x18\x03 \x01(\tR\x03msg\"H\n" +
	"\x0edecodeTokenRes\x12\x12\n" +
	"\x04data\x18\x01 \x01(\tR\x04data\x12\x10\n" +
	"\x03coe\x18\x02 \x01(\x05R\x03coe\x12\x10\n" +
	"\x03msg\x18\x03 \x01(\tR\x03msg2\x90\t\n" +
	"\aDeposit\x12\x82\x01\n" +
	"\rCreateDeposit\x12'.api.helloworld.v1.CreateDepositRequest\x1a%.api.helloworld.v1.CreateDepositReply\"!\x82\xd3\xe4\x93\x02\x1b:\x01*\"\x16/deposit/createDeposit\x12_\n" +
	"\rUpdateDeposit\x12'.api.helloworld.v1.UpdateDepositRequest\x1a%.api.helloworld.v1.UpdateDepositReply\x12_\n" +
	"\rDeleteDeposit\x12'.api.helloworld.v1.DeleteDepositRequest\x1a%.api.helloworld.v1.DeleteDepositReply\x12V\n" +
	"\n" +
	"GetDeposit\x12$.api.helloworld.v1.GetDepositRequest\x1a\".api.helloworld.v1.GetDepositReply\x12k\n" +
	"\vListDeposit\x12%.api.helloworld.v1.ListDepositRequest\x1a#.api.helloworld.v1.ListDepositReply\"\x10\x82\xd3\xe4\x93\x02\n" +
	"\x12\b/deposit\x12i\n" +
	"\vReturnToken\x12!.api.helloworld.v1.ReturnTokenReq\x1a!.api.helloworld.v1.ReturnTokenRes\"\x14\x82\xd3\xe4\x93\x02\x0e\x12\f/returntoken\x12i\n" +
	"\vDecodeToken\x12!.api.helloworld.v1.ReturnTokenReq\x1a!.api.helloworld.v1.ReturnTokenRes\"\x14\x82\xd3\xe4\x93\x02\x0e\x12\f/decodetoken\x12}\n" +
	"\x10GetDepositLocker\x12&.api.helloworld.v1.GetDepositLockerReq\x1a&.api.helloworld.v1.GetDepositLockerRes\"\x19\x82\xd3\xe4\x93\x02\x13\x12\x11/getDepositLocker\x12\x9c\x01\n" +
	"\x15UpdateDepositLockerId\x12+.api.helloworld.v1.UpdateDepositLockerIdReq\x1a+.api.helloworld.v1.UpdateDepositLockerIdRes\")\x82\xd3\xe4\x93\x02#:\x01*\"\x1e/deposit/updateDepositLockerId\x12\x84\x01\n" +
	"\x0fSendCodeByOrder\x12%.api.helloworld.v1.SendCodeByOrderReq\x1a%.api.helloworld.v1.SendCodeByOrderRes\"#\x82\xd3\xe4\x93\x02\x1d:\x01*\"\x18/deposit/sendCodeByOrderB7\n" +
	"\x11api.helloworld.v1P\x01Z ito-deposit/api/helloworld/v1;v1b\x06proto3"

var (
	file_helloworld_v1_deposit_proto_rawDescOnce sync.Once
	file_helloworld_v1_deposit_proto_rawDescData []byte
)

func file_helloworld_v1_deposit_proto_rawDescGZIP() []byte {
	file_helloworld_v1_deposit_proto_rawDescOnce.Do(func() {
		file_helloworld_v1_deposit_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_helloworld_v1_deposit_proto_rawDesc), len(file_helloworld_v1_deposit_proto_rawDesc)))
	})
	return file_helloworld_v1_deposit_proto_rawDescData
}

var file_helloworld_v1_deposit_proto_msgTypes = make([]protoimpl.MessageInfo, 21)
var file_helloworld_v1_deposit_proto_goTypes = []any{
	(*SendCodeByOrderReq)(nil),       // 0: api.helloworld.v1.SendCodeByOrderReq
	(*SendCodeByOrderRes)(nil),       // 1: api.helloworld.v1.SendCodeByOrderRes
	(*UpdateDepositLockerIdReq)(nil), // 2: api.helloworld.v1.UpdateDepositLockerIdReq
	(*UpdateDepositLockerIdRes)(nil), // 3: api.helloworld.v1.UpdateDepositLockerIdRes
	(*GetDepositLockerReq)(nil),      // 4: api.helloworld.v1.GetDepositLockerReq
	(*GetDepositLockerRes)(nil),      // 5: api.helloworld.v1.GetDepositLockerRes
	(*Locker)(nil),                   // 6: api.helloworld.v1.Locker
	(*CreateDepositRequest)(nil),     // 7: api.helloworld.v1.CreateDepositRequest
	(*CreateDepositReply)(nil),       // 8: api.helloworld.v1.CreateDepositReply
	(*DepositReplyData)(nil),         // 9: api.helloworld.v1.DepositReplyData
	(*UpdateDepositRequest)(nil),     // 10: api.helloworld.v1.UpdateDepositRequest
	(*UpdateDepositReply)(nil),       // 11: api.helloworld.v1.UpdateDepositReply
	(*DeleteDepositRequest)(nil),     // 12: api.helloworld.v1.DeleteDepositRequest
	(*DeleteDepositReply)(nil),       // 13: api.helloworld.v1.DeleteDepositReply
	(*GetDepositRequest)(nil),        // 14: api.helloworld.v1.GetDepositRequest
	(*GetDepositReply)(nil),          // 15: api.helloworld.v1.GetDepositReply
	(*ListDepositRequest)(nil),       // 16: api.helloworld.v1.ListDepositRequest
	(*ListDepositReply)(nil),         // 17: api.helloworld.v1.ListDepositReply
	(*ReturnTokenReq)(nil),           // 18: api.helloworld.v1.ReturnTokenReq
	(*ReturnTokenRes)(nil),           // 19: api.helloworld.v1.ReturnTokenRes
	(*DecodeTokenRes)(nil),           // 20: api.helloworld.v1.decodeTokenRes
}
var file_helloworld_v1_deposit_proto_depIdxs = []int32{
	6,  // 0: api.helloworld.v1.GetDepositLockerRes.locker:type_name -> api.helloworld.v1.Locker
	9,  // 1: api.helloworld.v1.CreateDepositReply.data:type_name -> api.helloworld.v1.DepositReplyData
	7,  // 2: api.helloworld.v1.Deposit.CreateDeposit:input_type -> api.helloworld.v1.CreateDepositRequest
	10, // 3: api.helloworld.v1.Deposit.UpdateDeposit:input_type -> api.helloworld.v1.UpdateDepositRequest
	12, // 4: api.helloworld.v1.Deposit.DeleteDeposit:input_type -> api.helloworld.v1.DeleteDepositRequest
	14, // 5: api.helloworld.v1.Deposit.GetDeposit:input_type -> api.helloworld.v1.GetDepositRequest
	16, // 6: api.helloworld.v1.Deposit.ListDeposit:input_type -> api.helloworld.v1.ListDepositRequest
	18, // 7: api.helloworld.v1.Deposit.ReturnToken:input_type -> api.helloworld.v1.ReturnTokenReq
	18, // 8: api.helloworld.v1.Deposit.DecodeToken:input_type -> api.helloworld.v1.ReturnTokenReq
	4,  // 9: api.helloworld.v1.Deposit.GetDepositLocker:input_type -> api.helloworld.v1.GetDepositLockerReq
	2,  // 10: api.helloworld.v1.Deposit.UpdateDepositLockerId:input_type -> api.helloworld.v1.UpdateDepositLockerIdReq
	0,  // 11: api.helloworld.v1.Deposit.SendCodeByOrder:input_type -> api.helloworld.v1.SendCodeByOrderReq
	8,  // 12: api.helloworld.v1.Deposit.CreateDeposit:output_type -> api.helloworld.v1.CreateDepositReply
	11, // 13: api.helloworld.v1.Deposit.UpdateDeposit:output_type -> api.helloworld.v1.UpdateDepositReply
	13, // 14: api.helloworld.v1.Deposit.DeleteDeposit:output_type -> api.helloworld.v1.DeleteDepositReply
	15, // 15: api.helloworld.v1.Deposit.GetDeposit:output_type -> api.helloworld.v1.GetDepositReply
	17, // 16: api.helloworld.v1.Deposit.ListDeposit:output_type -> api.helloworld.v1.ListDepositReply
	19, // 17: api.helloworld.v1.Deposit.ReturnToken:output_type -> api.helloworld.v1.ReturnTokenRes
	19, // 18: api.helloworld.v1.Deposit.DecodeToken:output_type -> api.helloworld.v1.ReturnTokenRes
	5,  // 19: api.helloworld.v1.Deposit.GetDepositLocker:output_type -> api.helloworld.v1.GetDepositLockerRes
	3,  // 20: api.helloworld.v1.Deposit.UpdateDepositLockerId:output_type -> api.helloworld.v1.UpdateDepositLockerIdRes
	1,  // 21: api.helloworld.v1.Deposit.SendCodeByOrder:output_type -> api.helloworld.v1.SendCodeByOrderRes
	12, // [12:22] is the sub-list for method output_type
	2,  // [2:12] is the sub-list for method input_type
	2,  // [2:2] is the sub-list for extension type_name
	2,  // [2:2] is the sub-list for extension extendee
	0,  // [0:2] is the sub-list for field type_name
}

func init() { file_helloworld_v1_deposit_proto_init() }
func file_helloworld_v1_deposit_proto_init() {
	if File_helloworld_v1_deposit_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_helloworld_v1_deposit_proto_rawDesc), len(file_helloworld_v1_deposit_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   21,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_helloworld_v1_deposit_proto_goTypes,
		DependencyIndexes: file_helloworld_v1_deposit_proto_depIdxs,
		MessageInfos:      file_helloworld_v1_deposit_proto_msgTypes,
	}.Build()
	File_helloworld_v1_deposit_proto = out.File
	file_helloworld_v1_deposit_proto_goTypes = nil
	file_helloworld_v1_deposit_proto_depIdxs = nil
}
