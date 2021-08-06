// Code generated by protoc-gen-go. DO NOT EDIT.
// source: logutil.proto

package logutil

import (
	fmt "fmt"
	math "math"

	vttime "github.com/rock-go/rock-chameleon-go/vitess/go/vt/proto/vttime"
	proto "github.com/golang/protobuf/proto"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

// Level is the level of the log messages.
type Level int32

const (
	// The usual logging levels.
	// Should be logged using logging facility.
	Level_INFO    Level = 0
	Level_WARNING Level = 1
	Level_ERROR   Level = 2
	// For messages that may contains non-logging events.
	// Should be logged to console directly.
	Level_CONSOLE Level = 3
)

var Level_name = map[int32]string{
	0: "INFO",
	1: "WARNING",
	2: "ERROR",
	3: "CONSOLE",
}

var Level_value = map[string]int32{
	"INFO":    0,
	"WARNING": 1,
	"ERROR":   2,
	"CONSOLE": 3,
}

func (x Level) String() string {
	return proto.EnumName(Level_name, int32(x))
}

func (Level) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_31f5dd3702a8edf9, []int{0}
}

// Event is a single logging event
type Event struct {
	Time                 *vttime.Time `protobuf:"bytes,1,opt,name=time,proto3" json:"time,omitempty"`
	Level                Level        `protobuf:"varint,2,opt,name=level,proto3,enum=logutil.Level" json:"level,omitempty"`
	File                 string       `protobuf:"bytes,3,opt,name=file,proto3" json:"file,omitempty"`
	Line                 int64        `protobuf:"varint,4,opt,name=line,proto3" json:"line,omitempty"`
	Value                string       `protobuf:"bytes,5,opt,name=value,proto3" json:"value,omitempty"`
	XXX_NoUnkeyedLiteral struct{}     `json:"-"`
	XXX_unrecognized     []byte       `json:"-"`
	XXX_sizecache        int32        `json:"-"`
}

func (m *Event) Reset()         { *m = Event{} }
func (m *Event) String() string { return proto.CompactTextString(m) }
func (*Event) ProtoMessage()    {}
func (*Event) Descriptor() ([]byte, []int) {
	return fileDescriptor_31f5dd3702a8edf9, []int{0}
}

func (m *Event) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Event.Unmarshal(m, b)
}
func (m *Event) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Event.Marshal(b, m, deterministic)
}
func (m *Event) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Event.Merge(m, src)
}
func (m *Event) XXX_Size() int {
	return xxx_messageInfo_Event.Size(m)
}
func (m *Event) XXX_DiscardUnknown() {
	xxx_messageInfo_Event.DiscardUnknown(m)
}

var xxx_messageInfo_Event proto.InternalMessageInfo

func (m *Event) GetTime() *vttime.Time {
	if m != nil {
		return m.Time
	}
	return nil
}

func (m *Event) GetLevel() Level {
	if m != nil {
		return m.Level
	}
	return Level_INFO
}

func (m *Event) GetFile() string {
	if m != nil {
		return m.File
	}
	return ""
}

func (m *Event) GetLine() int64 {
	if m != nil {
		return m.Line
	}
	return 0
}

func (m *Event) GetValue() string {
	if m != nil {
		return m.Value
	}
	return ""
}

func init() {
	proto.RegisterEnum("logutil.Level", Level_name, Level_value)
	proto.RegisterType((*Event)(nil), "logutil.Event")
}

func init() { proto.RegisterFile("logutil.proto", fileDescriptor_31f5dd3702a8edf9) }

var fileDescriptor_31f5dd3702a8edf9 = []byte{
	// 235 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x34, 0x8f, 0x41, 0x4b, 0x03, 0x31,
	0x10, 0x85, 0x4d, 0x77, 0x63, 0xed, 0x54, 0xcb, 0x32, 0x78, 0x08, 0x9e, 0x82, 0x14, 0x59, 0x3c,
	0x6c, 0xa0, 0x82, 0x77, 0x95, 0x55, 0x0a, 0x65, 0x17, 0xa2, 0x20, 0x78, 0x53, 0x18, 0x4b, 0x20,
	0x6d, 0xc4, 0xa6, 0xf9, 0x17, 0xfe, 0x67, 0xd9, 0x49, 0x7b, 0x7b, 0xef, 0x7b, 0x8f, 0xc7, 0x0c,
	0x5c, 0xf8, 0xb0, 0xde, 0x47, 0xe7, 0x9b, 0x9f, 0xdf, 0x10, 0x03, 0x8e, 0x0f, 0xf6, 0x0a, 0xa2,
	0xdb, 0x50, 0x86, 0xd7, 0x7f, 0x02, 0x64, 0x9b, 0x68, 0x1b, 0x51, 0x43, 0x39, 0x70, 0x25, 0xb4,
	0xa8, 0xa7, 0x8b, 0xf3, 0x26, 0x45, 0xae, 0xbd, 0xb9, 0x0d, 0x59, 0x4e, 0x70, 0x0e, 0xd2, 0x53,
	0x22, 0xaf, 0x46, 0x5a, 0xd4, 0xb3, 0xc5, 0xac, 0x39, 0xee, 0xaf, 0x06, 0x6a, 0x73, 0x88, 0x08,
	0xe5, 0xb7, 0xf3, 0xa4, 0x0a, 0x2d, 0xea, 0x89, 0x65, 0x3d, 0x30, 0xef, 0xb6, 0xa4, 0x4a, 0x2d,
	0xea, 0xc2, 0xb2, 0xc6, 0x4b, 0x90, 0xe9, 0xd3, 0xef, 0x49, 0x49, 0x2e, 0x66, 0x73, 0x7b, 0x0f,
	0x92, 0xd7, 0xf0, 0x0c, 0xca, 0x65, 0xf7, 0xdc, 0x57, 0x27, 0x38, 0x85, 0xf1, 0xfb, 0x83, 0xed,
	0x96, 0xdd, 0x4b, 0x25, 0x70, 0x02, 0xb2, 0xb5, 0xb6, 0xb7, 0xd5, 0x68, 0xe0, 0x4f, 0x7d, 0xf7,
	0xda, 0xaf, 0xda, 0xaa, 0x78, 0xbc, 0xf9, 0x98, 0x27, 0x17, 0x69, 0xb7, 0x6b, 0x5c, 0x30, 0x59,
	0x99, 0x75, 0x30, 0x29, 0x1a, 0xfe, 0xd3, 0x1c, 0x4e, 0xfd, 0x3a, 0x65, 0x7b, 0xf7, 0x1f, 0x00,
	0x00, 0xff, 0xff, 0xdd, 0xfa, 0x9b, 0x9a, 0x1c, 0x01, 0x00, 0x00,
}
