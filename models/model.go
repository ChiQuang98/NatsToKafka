package models

import "github.com/gogo/protobuf/proto"

type Channel struct {
	Id string `json:"id"`
	Name string `json:"name"`
}
type MessageNats struct {
	IdChannel string
	MessageData string
}
type Message struct {
	Channel              string   `protobuf:"bytes,1,opt,name=channel,proto3" json:"channel,omitempty"`
	Subtopic             string   `protobuf:"bytes,2,opt,name=subtopic,proto3" json:"subtopic,omitempty"`
	Publisher            string   `protobuf:"bytes,3,opt,name=publisher,proto3" json:"publisher,omitempty"`
	Protocol             string   `protobuf:"bytes,4,opt,name=protocol,proto3" json:"protocol,omitempty"`
	Payload              []byte   `protobuf:"bytes,5,opt,name=payload,proto3" json:"payload,omitempty"`
	Created              int64    `protobuf:"varint,6,opt,name=created,proto3" json:"created,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}
func (m *Message) Reset()         { *m = Message{} }
func (m *Message) String() string { return proto.CompactTextString(m) }
func (*Message) ProtoMessage()    {}







