// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.25.0
// 	protoc        v3.6.1
// source: newsservice.proto

package news

import (
	context "context"
	proto "github.com/golang/protobuf/proto"
	_struct "github.com/golang/protobuf/ptypes/struct"
	_ "google.golang.org/genproto/googleapis/api/annotations"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
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

// This is a compile-time assertion that a sufficiently up-to-date version
// of the legacy proto package is being used.
const _ = proto.ProtoPackageIsVersion4

type LastNews struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	News []*_struct.Struct `protobuf:"bytes,1,rep,name=news,proto3" json:"news,omitempty"`
}

func (x *LastNews) Reset() {
	*x = LastNews{}
	if protoimpl.UnsafeEnabled {
		mi := &file_newsservice_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *LastNews) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LastNews) ProtoMessage() {}

func (x *LastNews) ProtoReflect() protoreflect.Message {
	mi := &file_newsservice_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LastNews.ProtoReflect.Descriptor instead.
func (*LastNews) Descriptor() ([]byte, []int) {
	return file_newsservice_proto_rawDescGZIP(), []int{0}
}

func (x *LastNews) GetNews() []*_struct.Struct {
	if x != nil {
		return x.News
	}
	return nil
}

type MessageParams struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Sources  string `protobuf:"bytes,1,opt,name=sources,proto3" json:"sources,omitempty"`
	Keywords string `protobuf:"bytes,2,opt,name=keywords,proto3" json:"keywords,omitempty"`
}

func (x *MessageParams) Reset() {
	*x = MessageParams{}
	if protoimpl.UnsafeEnabled {
		mi := &file_newsservice_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MessageParams) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MessageParams) ProtoMessage() {}

func (x *MessageParams) ProtoReflect() protoreflect.Message {
	mi := &file_newsservice_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MessageParams.ProtoReflect.Descriptor instead.
func (*MessageParams) Descriptor() ([]byte, []int) {
	return file_newsservice_proto_rawDescGZIP(), []int{1}
}

func (x *MessageParams) GetSources() string {
	if x != nil {
		return x.Sources
	}
	return ""
}

func (x *MessageParams) GetKeywords() string {
	if x != nil {
		return x.Keywords
	}
	return ""
}

var File_newsservice_proto protoreflect.FileDescriptor

var file_newsservice_proto_rawDesc = []byte{
	0x0a, 0x11, 0x6e, 0x65, 0x77, 0x73, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x12, 0x04, 0x6e, 0x65, 0x77, 0x73, 0x1a, 0x1c, 0x67, 0x6f, 0x6f, 0x67, 0x6c,
	0x65, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x61, 0x6e, 0x6e, 0x6f, 0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e,
	0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1c, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x73, 0x74, 0x72, 0x75, 0x63, 0x74, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x37, 0x0a, 0x08, 0x4c, 0x61, 0x73, 0x74, 0x4e, 0x65, 0x77,
	0x73, 0x12, 0x2b, 0x0a, 0x04, 0x6e, 0x65, 0x77, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32,
	0x17, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75,
	0x66, 0x2e, 0x53, 0x74, 0x72, 0x75, 0x63, 0x74, 0x52, 0x04, 0x6e, 0x65, 0x77, 0x73, 0x22, 0x45,
	0x0a, 0x0d, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x50, 0x61, 0x72, 0x61, 0x6d, 0x73, 0x12,
	0x18, 0x0a, 0x07, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x07, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x12, 0x1a, 0x0a, 0x08, 0x6b, 0x65, 0x79,
	0x77, 0x6f, 0x72, 0x64, 0x73, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x6b, 0x65, 0x79,
	0x77, 0x6f, 0x72, 0x64, 0x73, 0x32, 0x57, 0x0a, 0x0b, 0x4e, 0x65, 0x77, 0x73, 0x53, 0x65, 0x72,
	0x76, 0x69, 0x63, 0x65, 0x12, 0x48, 0x0a, 0x07, 0x47, 0x65, 0x74, 0x4e, 0x65, 0x77, 0x73, 0x12,
	0x13, 0x2e, 0x6e, 0x65, 0x77, 0x73, 0x2e, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x50, 0x61,
	0x72, 0x61, 0x6d, 0x73, 0x1a, 0x0e, 0x2e, 0x6e, 0x65, 0x77, 0x73, 0x2e, 0x4c, 0x61, 0x73, 0x74,
	0x4e, 0x65, 0x77, 0x73, 0x22, 0x18, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x12, 0x22, 0x0d, 0x2f, 0x76,
	0x31, 0x2f, 0x6e, 0x65, 0x77, 0x73, 0x2f, 0x6c, 0x61, 0x73, 0x74, 0x3a, 0x01, 0x2a, 0x62, 0x06,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_newsservice_proto_rawDescOnce sync.Once
	file_newsservice_proto_rawDescData = file_newsservice_proto_rawDesc
)

func file_newsservice_proto_rawDescGZIP() []byte {
	file_newsservice_proto_rawDescOnce.Do(func() {
		file_newsservice_proto_rawDescData = protoimpl.X.CompressGZIP(file_newsservice_proto_rawDescData)
	})
	return file_newsservice_proto_rawDescData
}

var file_newsservice_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_newsservice_proto_goTypes = []interface{}{
	(*LastNews)(nil),       // 0: news.LastNews
	(*MessageParams)(nil),  // 1: news.MessageParams
	(*_struct.Struct)(nil), // 2: google.protobuf.Struct
}
var file_newsservice_proto_depIdxs = []int32{
	2, // 0: news.LastNews.news:type_name -> google.protobuf.Struct
	1, // 1: news.NewsService.GetNews:input_type -> news.MessageParams
	0, // 2: news.NewsService.GetNews:output_type -> news.LastNews
	2, // [2:3] is the sub-list for method output_type
	1, // [1:2] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_newsservice_proto_init() }
func file_newsservice_proto_init() {
	if File_newsservice_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_newsservice_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*LastNews); i {
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
		file_newsservice_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*MessageParams); i {
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
			RawDescriptor: file_newsservice_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_newsservice_proto_goTypes,
		DependencyIndexes: file_newsservice_proto_depIdxs,
		MessageInfos:      file_newsservice_proto_msgTypes,
	}.Build()
	File_newsservice_proto = out.File
	file_newsservice_proto_rawDesc = nil
	file_newsservice_proto_goTypes = nil
	file_newsservice_proto_depIdxs = nil
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConnInterface

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion6

// NewsServiceClient is the client API for NewsService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type NewsServiceClient interface {
	GetNews(ctx context.Context, in *MessageParams, opts ...grpc.CallOption) (*LastNews, error)
}

type newsServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewNewsServiceClient(cc grpc.ClientConnInterface) NewsServiceClient {
	return &newsServiceClient{cc}
}

func (c *newsServiceClient) GetNews(ctx context.Context, in *MessageParams, opts ...grpc.CallOption) (*LastNews, error) {
	out := new(LastNews)
	err := c.cc.Invoke(ctx, "/news.NewsService/GetNews", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// NewsServiceServer is the server API for NewsService service.
type NewsServiceServer interface {
	GetNews(context.Context, *MessageParams) (*LastNews, error)
}

// UnimplementedNewsServiceServer can be embedded to have forward compatible implementations.
type UnimplementedNewsServiceServer struct {
}

func (*UnimplementedNewsServiceServer) GetNews(context.Context, *MessageParams) (*LastNews, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetNews not implemented")
}

func RegisterNewsServiceServer(s *grpc.Server, srv NewsServiceServer) {
	s.RegisterService(&_NewsService_serviceDesc, srv)
}

func _NewsService_GetNews_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MessageParams)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NewsServiceServer).GetNews(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/news.NewsService/GetNews",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NewsServiceServer).GetNews(ctx, req.(*MessageParams))
	}
	return interceptor(ctx, in, info, handler)
}

var _NewsService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "news.NewsService",
	HandlerType: (*NewsServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetNews",
			Handler:    _NewsService_GetNews_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "newsservice.proto",
}
