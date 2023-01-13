// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.21.7
// source: api.proto

package proto

import (
	context "context"

	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// GophKeeperServerClient is the client API for GophKeeperServer service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type GophKeeperServerClient interface {
	Register(ctx context.Context, in *RegisterRequest, opts ...grpc.CallOption) (*RegisterResponse, error)
	SignIn(ctx context.Context, in *RegisterRequest, opts ...grpc.CallOption) (*RegisterResponse, error)
	AddCard(ctx context.Context, in *AddCardRequest, opts ...grpc.CallOption) (*AddCardResponse, error)
	AddLogin(ctx context.Context, in *AddLoginRequest, opts ...grpc.CallOption) (*AddLoginResponse, error)
	AddText(ctx context.Context, in *AddTextRequest, opts ...grpc.CallOption) (*AddTextResponse, error)
	AddBinary(ctx context.Context, in *AddBinRequest, opts ...grpc.CallOption) (*AddBinResponse, error)
	Card(ctx context.Context, in *CardRequest, opts ...grpc.CallOption) (*CardResponse, error)
	Login(ctx context.Context, in *LoginRequest, opts ...grpc.CallOption) (*LoginResponse, error)
	Text(ctx context.Context, in *TextRequest, opts ...grpc.CallOption) (*TextResponse, error)
	Binary(ctx context.Context, in *BinRequest, opts ...grpc.CallOption) (*BinResponse, error)
	DeleteCard(ctx context.Context, in *DeleteCardRequest, opts ...grpc.CallOption) (*DeleteCardResponse, error)
	DeleteLogin(ctx context.Context, in *DeleteLoginRequest, opts ...grpc.CallOption) (*DeleteLoginResponse, error)
	DeleteText(ctx context.Context, in *DeleteTextRequest, opts ...grpc.CallOption) (*DeleteTextResponse, error)
	DeleteBinary(ctx context.Context, in *DeleteBinRequest, opts ...grpc.CallOption) (*DeleteBinResponse, error)
}

type gophKeeperServerClient struct {
	cc grpc.ClientConnInterface
}

func NewGophKeeperServerClient(cc grpc.ClientConnInterface) GophKeeperServerClient {
	return &gophKeeperServerClient{cc}
}

func (c *gophKeeperServerClient) Register(ctx context.Context, in *RegisterRequest, opts ...grpc.CallOption) (*RegisterResponse, error) {
	out := new(RegisterResponse)
	err := c.cc.Invoke(ctx, "/proto.GophKeeperServer/Register", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *gophKeeperServerClient) SignIn(ctx context.Context, in *RegisterRequest, opts ...grpc.CallOption) (*RegisterResponse, error) {
	out := new(RegisterResponse)
	err := c.cc.Invoke(ctx, "/proto.GophKeeperServer/SignIn", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *gophKeeperServerClient) AddCard(ctx context.Context, in *AddCardRequest, opts ...grpc.CallOption) (*AddCardResponse, error) {
	out := new(AddCardResponse)
	err := c.cc.Invoke(ctx, "/proto.GophKeeperServer/AddCard", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *gophKeeperServerClient) AddLogin(ctx context.Context, in *AddLoginRequest, opts ...grpc.CallOption) (*AddLoginResponse, error) {
	out := new(AddLoginResponse)
	err := c.cc.Invoke(ctx, "/proto.GophKeeperServer/AddLogin", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *gophKeeperServerClient) AddText(ctx context.Context, in *AddTextRequest, opts ...grpc.CallOption) (*AddTextResponse, error) {
	out := new(AddTextResponse)
	err := c.cc.Invoke(ctx, "/proto.GophKeeperServer/AddText", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *gophKeeperServerClient) AddBinary(ctx context.Context, in *AddBinRequest, opts ...grpc.CallOption) (*AddBinResponse, error) {
	out := new(AddBinResponse)
	err := c.cc.Invoke(ctx, "/proto.GophKeeperServer/AddBinary", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *gophKeeperServerClient) Card(ctx context.Context, in *CardRequest, opts ...grpc.CallOption) (*CardResponse, error) {
	out := new(CardResponse)
	err := c.cc.Invoke(ctx, "/proto.GophKeeperServer/Card", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *gophKeeperServerClient) Login(ctx context.Context, in *LoginRequest, opts ...grpc.CallOption) (*LoginResponse, error) {
	out := new(LoginResponse)
	err := c.cc.Invoke(ctx, "/proto.GophKeeperServer/Login", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *gophKeeperServerClient) Text(ctx context.Context, in *TextRequest, opts ...grpc.CallOption) (*TextResponse, error) {
	out := new(TextResponse)
	err := c.cc.Invoke(ctx, "/proto.GophKeeperServer/Text", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *gophKeeperServerClient) Binary(ctx context.Context, in *BinRequest, opts ...grpc.CallOption) (*BinResponse, error) {
	out := new(BinResponse)
	err := c.cc.Invoke(ctx, "/proto.GophKeeperServer/Binary", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *gophKeeperServerClient) DeleteCard(ctx context.Context, in *DeleteCardRequest, opts ...grpc.CallOption) (*DeleteCardResponse, error) {
	out := new(DeleteCardResponse)
	err := c.cc.Invoke(ctx, "/proto.GophKeeperServer/DeleteCard", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *gophKeeperServerClient) DeleteLogin(ctx context.Context, in *DeleteLoginRequest, opts ...grpc.CallOption) (*DeleteLoginResponse, error) {
	out := new(DeleteLoginResponse)
	err := c.cc.Invoke(ctx, "/proto.GophKeeperServer/DeleteLogin", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *gophKeeperServerClient) DeleteText(ctx context.Context, in *DeleteTextRequest, opts ...grpc.CallOption) (*DeleteTextResponse, error) {
	out := new(DeleteTextResponse)
	err := c.cc.Invoke(ctx, "/proto.GophKeeperServer/DeleteText", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *gophKeeperServerClient) DeleteBinary(ctx context.Context, in *DeleteBinRequest, opts ...grpc.CallOption) (*DeleteBinResponse, error) {
	out := new(DeleteBinResponse)
	err := c.cc.Invoke(ctx, "/proto.GophKeeperServer/DeleteBinary", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// GophKeeperServerServer is the server API for GophKeeperServer service.
// All implementations must embed UnimplementedGophKeeperServerServer
// for forward compatibility
type GophKeeperServerServer interface {
	Register(context.Context, *RegisterRequest) (*RegisterResponse, error)
	SignIn(context.Context, *RegisterRequest) (*RegisterResponse, error)
	AddCard(context.Context, *AddCardRequest) (*AddCardResponse, error)
	AddLogin(context.Context, *AddLoginRequest) (*AddLoginResponse, error)
	AddText(context.Context, *AddTextRequest) (*AddTextResponse, error)
	AddBinary(context.Context, *AddBinRequest) (*AddBinResponse, error)
	Card(context.Context, *CardRequest) (*CardResponse, error)
	Login(context.Context, *LoginRequest) (*LoginResponse, error)
	Text(context.Context, *TextRequest) (*TextResponse, error)
	Binary(context.Context, *BinRequest) (*BinResponse, error)
	DeleteCard(context.Context, *DeleteCardRequest) (*DeleteCardResponse, error)
	DeleteLogin(context.Context, *DeleteLoginRequest) (*DeleteLoginResponse, error)
	DeleteText(context.Context, *DeleteTextRequest) (*DeleteTextResponse, error)
	DeleteBinary(context.Context, *DeleteBinRequest) (*DeleteBinResponse, error)
	mustEmbedUnimplementedGophKeeperServerServer()
}

// UnimplementedGophKeeperServerServer must be embedded to have forward compatible implementations.
type UnimplementedGophKeeperServerServer struct {
}

func (UnimplementedGophKeeperServerServer) Register(context.Context, *RegisterRequest) (*RegisterResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Register not implemented")
}
func (UnimplementedGophKeeperServerServer) SignIn(context.Context, *RegisterRequest) (*RegisterResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SignIn not implemented")
}
func (UnimplementedGophKeeperServerServer) AddCard(context.Context, *AddCardRequest) (*AddCardResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AddCard not implemented")
}
func (UnimplementedGophKeeperServerServer) AddLogin(context.Context, *AddLoginRequest) (*AddLoginResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AddLogin not implemented")
}
func (UnimplementedGophKeeperServerServer) AddText(context.Context, *AddTextRequest) (*AddTextResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AddText not implemented")
}
func (UnimplementedGophKeeperServerServer) AddBinary(context.Context, *AddBinRequest) (*AddBinResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AddBinary not implemented")
}
func (UnimplementedGophKeeperServerServer) Card(context.Context, *CardRequest) (*CardResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Card not implemented")
}
func (UnimplementedGophKeeperServerServer) Login(context.Context, *LoginRequest) (*LoginResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Login not implemented")
}
func (UnimplementedGophKeeperServerServer) Text(context.Context, *TextRequest) (*TextResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Text not implemented")
}
func (UnimplementedGophKeeperServerServer) Binary(context.Context, *BinRequest) (*BinResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Binary not implemented")
}
func (UnimplementedGophKeeperServerServer) DeleteCard(context.Context, *DeleteCardRequest) (*DeleteCardResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteCard not implemented")
}
func (UnimplementedGophKeeperServerServer) DeleteLogin(context.Context, *DeleteLoginRequest) (*DeleteLoginResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteLogin not implemented")
}
func (UnimplementedGophKeeperServerServer) DeleteText(context.Context, *DeleteTextRequest) (*DeleteTextResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteText not implemented")
}
func (UnimplementedGophKeeperServerServer) DeleteBinary(context.Context, *DeleteBinRequest) (*DeleteBinResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteBinary not implemented")
}
func (UnimplementedGophKeeperServerServer) mustEmbedUnimplementedGophKeeperServerServer() {}

// UnsafeGophKeeperServerServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to GophKeeperServerServer will
// result in compilation errors.
type UnsafeGophKeeperServerServer interface {
	mustEmbedUnimplementedGophKeeperServerServer()
}

func RegisterGophKeeperServerServer(s grpc.ServiceRegistrar, srv GophKeeperServerServer) {
	s.RegisterService(&GophKeeperServer_ServiceDesc, srv)
}

func _GophKeeperServer_Register_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RegisterRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GophKeeperServerServer).Register(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.GophKeeperServer/Register",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GophKeeperServerServer).Register(ctx, req.(*RegisterRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _GophKeeperServer_SignIn_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RegisterRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GophKeeperServerServer).SignIn(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.GophKeeperServer/SignIn",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GophKeeperServerServer).SignIn(ctx, req.(*RegisterRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _GophKeeperServer_AddCard_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AddCardRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GophKeeperServerServer).AddCard(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.GophKeeperServer/AddCard",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GophKeeperServerServer).AddCard(ctx, req.(*AddCardRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _GophKeeperServer_AddLogin_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AddLoginRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GophKeeperServerServer).AddLogin(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.GophKeeperServer/AddLogin",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GophKeeperServerServer).AddLogin(ctx, req.(*AddLoginRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _GophKeeperServer_AddText_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AddTextRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GophKeeperServerServer).AddText(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.GophKeeperServer/AddText",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GophKeeperServerServer).AddText(ctx, req.(*AddTextRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _GophKeeperServer_AddBinary_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AddBinRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GophKeeperServerServer).AddBinary(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.GophKeeperServer/AddBinary",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GophKeeperServerServer).AddBinary(ctx, req.(*AddBinRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _GophKeeperServer_Card_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CardRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GophKeeperServerServer).Card(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.GophKeeperServer/Card",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GophKeeperServerServer).Card(ctx, req.(*CardRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _GophKeeperServer_Login_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(LoginRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GophKeeperServerServer).Login(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.GophKeeperServer/Login",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GophKeeperServerServer).Login(ctx, req.(*LoginRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _GophKeeperServer_Text_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(TextRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GophKeeperServerServer).Text(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.GophKeeperServer/Text",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GophKeeperServerServer).Text(ctx, req.(*TextRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _GophKeeperServer_Binary_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(BinRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GophKeeperServerServer).Binary(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.GophKeeperServer/Binary",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GophKeeperServerServer).Binary(ctx, req.(*BinRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _GophKeeperServer_DeleteCard_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteCardRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GophKeeperServerServer).DeleteCard(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.GophKeeperServer/DeleteCard",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GophKeeperServerServer).DeleteCard(ctx, req.(*DeleteCardRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _GophKeeperServer_DeleteLogin_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteLoginRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GophKeeperServerServer).DeleteLogin(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.GophKeeperServer/DeleteLogin",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GophKeeperServerServer).DeleteLogin(ctx, req.(*DeleteLoginRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _GophKeeperServer_DeleteText_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteTextRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GophKeeperServerServer).DeleteText(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.GophKeeperServer/DeleteText",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GophKeeperServerServer).DeleteText(ctx, req.(*DeleteTextRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _GophKeeperServer_DeleteBinary_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteBinRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GophKeeperServerServer).DeleteBinary(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.GophKeeperServer/DeleteBinary",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GophKeeperServerServer).DeleteBinary(ctx, req.(*DeleteBinRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// GophKeeperServer_ServiceDesc is the grpc.ServiceDesc for GophKeeperServer service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var GophKeeperServer_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "proto.GophKeeperServer",
	HandlerType: (*GophKeeperServerServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Register",
			Handler:    _GophKeeperServer_Register_Handler,
		},
		{
			MethodName: "SignIn",
			Handler:    _GophKeeperServer_SignIn_Handler,
		},
		{
			MethodName: "AddCard",
			Handler:    _GophKeeperServer_AddCard_Handler,
		},
		{
			MethodName: "AddLogin",
			Handler:    _GophKeeperServer_AddLogin_Handler,
		},
		{
			MethodName: "AddText",
			Handler:    _GophKeeperServer_AddText_Handler,
		},
		{
			MethodName: "AddBinary",
			Handler:    _GophKeeperServer_AddBinary_Handler,
		},
		{
			MethodName: "Card",
			Handler:    _GophKeeperServer_Card_Handler,
		},
		{
			MethodName: "Login",
			Handler:    _GophKeeperServer_Login_Handler,
		},
		{
			MethodName: "Text",
			Handler:    _GophKeeperServer_Text_Handler,
		},
		{
			MethodName: "Binary",
			Handler:    _GophKeeperServer_Binary_Handler,
		},
		{
			MethodName: "DeleteCard",
			Handler:    _GophKeeperServer_DeleteCard_Handler,
		},
		{
			MethodName: "DeleteLogin",
			Handler:    _GophKeeperServer_DeleteLogin_Handler,
		},
		{
			MethodName: "DeleteText",
			Handler:    _GophKeeperServer_DeleteText_Handler,
		},
		{
			MethodName: "DeleteBinary",
			Handler:    _GophKeeperServer_DeleteBinary_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "api.proto",
}
