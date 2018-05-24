// Code generated by protoc-gen-go. DO NOT EDIT.
// source: Sentiment.proto

package protos

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"
import timestamp "github.com/golang/protobuf/ptypes/timestamp"

import (
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

// Represents a single aggregated score for a single Thread
type ThreadSentiment struct {
	Score                float64              `protobuf:"fixed64,1,opt,name=score" json:"score,omitempty"`
	CreatedTime          *timestamp.Timestamp `protobuf:"bytes,2,opt,name=createdTime" json:"createdTime,omitempty"`
	Thread               *Thread              `protobuf:"bytes,3,opt,name=thread" json:"thread,omitempty"`
	XXX_NoUnkeyedLiteral struct{}             `json:"-"`
	XXX_unrecognized     []byte               `json:"-"`
	XXX_sizecache        int32                `json:"-"`
}

func (m *ThreadSentiment) Reset()         { *m = ThreadSentiment{} }
func (m *ThreadSentiment) String() string { return proto.CompactTextString(m) }
func (*ThreadSentiment) ProtoMessage()    {}
func (*ThreadSentiment) Descriptor() ([]byte, []int) {
	return fileDescriptor_Sentiment_66f68986a74aa617, []int{0}
}
func (m *ThreadSentiment) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ThreadSentiment.Unmarshal(m, b)
}
func (m *ThreadSentiment) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ThreadSentiment.Marshal(b, m, deterministic)
}
func (dst *ThreadSentiment) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ThreadSentiment.Merge(dst, src)
}
func (m *ThreadSentiment) XXX_Size() int {
	return xxx_messageInfo_ThreadSentiment.Size(m)
}
func (m *ThreadSentiment) XXX_DiscardUnknown() {
	xxx_messageInfo_ThreadSentiment.DiscardUnknown(m)
}

var xxx_messageInfo_ThreadSentiment proto.InternalMessageInfo

func (m *ThreadSentiment) GetScore() float64 {
	if m != nil {
		return m.Score
	}
	return 0
}

func (m *ThreadSentiment) GetCreatedTime() *timestamp.Timestamp {
	if m != nil {
		return m.CreatedTime
	}
	return nil
}

func (m *ThreadSentiment) GetThread() *Thread {
	if m != nil {
		return m.Thread
	}
	return nil
}

// Represents a collection of ThreadSentiment scores
type GetAllCurrentClientThreadsSentimentResponse struct {
	Sentiments           []*ThreadSentiment `protobuf:"bytes,1,rep,name=sentiments" json:"sentiments,omitempty"`
	XXX_NoUnkeyedLiteral struct{}           `json:"-"`
	XXX_unrecognized     []byte             `json:"-"`
	XXX_sizecache        int32              `json:"-"`
}

func (m *GetAllCurrentClientThreadsSentimentResponse) Reset() {
	*m = GetAllCurrentClientThreadsSentimentResponse{}
}
func (m *GetAllCurrentClientThreadsSentimentResponse) String() string {
	return proto.CompactTextString(m)
}
func (*GetAllCurrentClientThreadsSentimentResponse) ProtoMessage() {}
func (*GetAllCurrentClientThreadsSentimentResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_Sentiment_66f68986a74aa617, []int{1}
}
func (m *GetAllCurrentClientThreadsSentimentResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetAllCurrentClientThreadsSentimentResponse.Unmarshal(m, b)
}
func (m *GetAllCurrentClientThreadsSentimentResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetAllCurrentClientThreadsSentimentResponse.Marshal(b, m, deterministic)
}
func (dst *GetAllCurrentClientThreadsSentimentResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetAllCurrentClientThreadsSentimentResponse.Merge(dst, src)
}
func (m *GetAllCurrentClientThreadsSentimentResponse) XXX_Size() int {
	return xxx_messageInfo_GetAllCurrentClientThreadsSentimentResponse.Size(m)
}
func (m *GetAllCurrentClientThreadsSentimentResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_GetAllCurrentClientThreadsSentimentResponse.DiscardUnknown(m)
}

var xxx_messageInfo_GetAllCurrentClientThreadsSentimentResponse proto.InternalMessageInfo

func (m *GetAllCurrentClientThreadsSentimentResponse) GetSentiments() []*ThreadSentiment {
	if m != nil {
		return m.Sentiments
	}
	return nil
}

func init() {
	proto.RegisterType((*ThreadSentiment)(nil), "protos.ThreadSentiment")
	proto.RegisterType((*GetAllCurrentClientThreadsSentimentResponse)(nil), "protos.GetAllCurrentClientThreadsSentimentResponse")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// SentimentAnalysisClient is the client API for SentimentAnalysis service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type SentimentAnalysisClient interface {
	// Given a Thread, returns that thread's aggregated sentiment
	GetCurrentThreadSentiment(ctx context.Context, in *Thread, opts ...grpc.CallOption) (*ThreadSentiment, error)
	// Given a Cient, returns a collection of all current Thread sentiments
	GetAllCurrentClientThreadsSentiment(ctx context.Context, in *Client, opts ...grpc.CallOption) (*GetAllCurrentClientThreadsSentimentResponse, error)
}

type sentimentAnalysisClient struct {
	cc *grpc.ClientConn
}

func NewSentimentAnalysisClient(cc *grpc.ClientConn) SentimentAnalysisClient {
	return &sentimentAnalysisClient{cc}
}

func (c *sentimentAnalysisClient) GetCurrentThreadSentiment(ctx context.Context, in *Thread, opts ...grpc.CallOption) (*ThreadSentiment, error) {
	out := new(ThreadSentiment)
	err := c.cc.Invoke(ctx, "/protos.SentimentAnalysis/GetCurrentThreadSentiment", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *sentimentAnalysisClient) GetAllCurrentClientThreadsSentiment(ctx context.Context, in *Client, opts ...grpc.CallOption) (*GetAllCurrentClientThreadsSentimentResponse, error) {
	out := new(GetAllCurrentClientThreadsSentimentResponse)
	err := c.cc.Invoke(ctx, "/protos.SentimentAnalysis/GetAllCurrentClientThreadsSentiment", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// SentimentAnalysisServer is the server API for SentimentAnalysis service.
type SentimentAnalysisServer interface {
	// Given a Thread, returns that thread's aggregated sentiment
	GetCurrentThreadSentiment(context.Context, *Thread) (*ThreadSentiment, error)
	// Given a Cient, returns a collection of all current Thread sentiments
	GetAllCurrentClientThreadsSentiment(context.Context, *Client) (*GetAllCurrentClientThreadsSentimentResponse, error)
}

func RegisterSentimentAnalysisServer(s *grpc.Server, srv SentimentAnalysisServer) {
	s.RegisterService(&_SentimentAnalysis_serviceDesc, srv)
}

func _SentimentAnalysis_GetCurrentThreadSentiment_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Thread)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SentimentAnalysisServer).GetCurrentThreadSentiment(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/protos.SentimentAnalysis/GetCurrentThreadSentiment",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SentimentAnalysisServer).GetCurrentThreadSentiment(ctx, req.(*Thread))
	}
	return interceptor(ctx, in, info, handler)
}

func _SentimentAnalysis_GetAllCurrentClientThreadsSentiment_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Client)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SentimentAnalysisServer).GetAllCurrentClientThreadsSentiment(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/protos.SentimentAnalysis/GetAllCurrentClientThreadsSentiment",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SentimentAnalysisServer).GetAllCurrentClientThreadsSentiment(ctx, req.(*Client))
	}
	return interceptor(ctx, in, info, handler)
}

var _SentimentAnalysis_serviceDesc = grpc.ServiceDesc{
	ServiceName: "protos.SentimentAnalysis",
	HandlerType: (*SentimentAnalysisServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetCurrentThreadSentiment",
			Handler:    _SentimentAnalysis_GetCurrentThreadSentiment_Handler,
		},
		{
			MethodName: "GetAllCurrentClientThreadsSentiment",
			Handler:    _SentimentAnalysis_GetAllCurrentClientThreadsSentiment_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "Sentiment.proto",
}

func init() { proto.RegisterFile("Sentiment.proto", fileDescriptor_Sentiment_66f68986a74aa617) }

var fileDescriptor_Sentiment_66f68986a74aa617 = []byte{
	// 273 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x94, 0x91, 0xc1, 0x4a, 0xc3, 0x40,
	0x10, 0x86, 0x59, 0x8b, 0x05, 0x27, 0x60, 0x71, 0x11, 0x8c, 0xb9, 0x18, 0x2a, 0x48, 0x40, 0xd8,
	0x42, 0x7a, 0xf0, 0xe2, 0xa5, 0x56, 0xe8, 0x3d, 0xe6, 0x05, 0xd2, 0x76, 0x5a, 0x23, 0x9b, 0xdd,
	0xb0, 0x33, 0x3d, 0xf8, 0x10, 0x3e, 0x96, 0xef, 0x25, 0x66, 0x93, 0xb4, 0x16, 0x04, 0x3d, 0x85,
	0x99, 0xf9, 0xff, 0xff, 0x9b, 0xcc, 0xc2, 0xe8, 0x05, 0x0d, 0x97, 0x15, 0x1a, 0x56, 0xb5, 0xb3,
	0x6c, 0xe5, 0xb0, 0xf9, 0x50, 0x74, 0xf6, 0x64, 0xdb, 0x56, 0x74, 0xb3, 0xb5, 0x76, 0xab, 0x71,
	0xd2, 0x54, 0xcb, 0xdd, 0x66, 0xf2, 0x6d, 0x20, 0x2e, 0xaa, 0xda, 0x0b, 0xc6, 0x1f, 0x02, 0x46,
	0xf9, 0xab, 0xc3, 0x62, 0xdd, 0xa7, 0xc9, 0x4b, 0x38, 0xa5, 0x95, 0x75, 0x18, 0x8a, 0x58, 0x24,
	0x22, 0xf3, 0x85, 0x7c, 0x84, 0x60, 0xe5, 0xb0, 0x60, 0x5c, 0xe7, 0x65, 0x85, 0xe1, 0x49, 0x2c,
	0x92, 0x20, 0x8d, 0x94, 0x07, 0xa8, 0x0e, 0xa0, 0xf2, 0x0e, 0x90, 0x1d, 0xca, 0xe5, 0x1d, 0x0c,
	0xb9, 0xc1, 0x84, 0x83, 0xc6, 0x78, 0xee, 0x1d, 0xa4, 0x3c, 0x3c, 0x6b, 0xa7, 0xe3, 0x0d, 0xdc,
	0x2f, 0x90, 0x67, 0x5a, 0xcf, 0x77, 0xce, 0xa1, 0xe1, 0xb9, 0x2e, 0xd1, 0xb0, 0x17, 0x51, 0xbf,
	0x62, 0x86, 0x54, 0x5b, 0x43, 0x28, 0x1f, 0x00, 0xa8, 0x6b, 0x52, 0x28, 0xe2, 0x41, 0x12, 0xa4,
	0x57, 0x3f, 0xa3, 0xf7, 0xa6, 0x03, 0x69, 0xfa, 0x29, 0xe0, 0xa2, 0x9f, 0xcc, 0x4c, 0xa1, 0xdf,
	0xa9, 0x24, 0xf9, 0x0c, 0xd7, 0x0b, 0xe4, 0x16, 0x7d, 0x7c, 0x96, 0xa3, 0x95, 0xa3, 0xdf, 0x38,
	0xf2, 0x0d, 0x6e, 0xff, 0xf0, 0x0f, 0xfb, 0x3c, 0x3f, 0x8f, 0xa6, 0x5d, 0xfd, 0x8f, 0x03, 0x2c,
	0xfd, 0x9b, 0x4f, 0xbf, 0x02, 0x00, 0x00, 0xff, 0xff, 0x12, 0x9a, 0x02, 0x8e, 0x0d, 0x02, 0x00,
	0x00,
}
