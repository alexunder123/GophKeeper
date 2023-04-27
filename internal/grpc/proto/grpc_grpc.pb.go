// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v4.22.2
// source: proto/grpc.proto

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

const (
	GophKeeper_NewSessionID_FullMethodName = "/grpc.GophKeeper/NewSessionID"
	GophKeeper_NewUser_FullMethodName      = "/grpc.GophKeeper/NewUser"
	GophKeeper_LoginUser_FullMethodName    = "/grpc.GophKeeper/LoginUser"
	GophKeeper_UserData_FullMethodName     = "/grpc.GophKeeper/UserData"
	GophKeeper_TimeStamp_FullMethodName    = "/grpc.GophKeeper/TimeStamp"
	GophKeeper_UpdateData_FullMethodName   = "/grpc.GophKeeper/UpdateData"
	GophKeeper_LogOut_FullMethodName       = "/grpc.GophKeeper/LogOut"
)

// GophKeeperClient is the client API for GophKeeper service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type GophKeeperClient interface {
	NewSessionID(ctx context.Context, in *NewSessionIDRequest, opts ...grpc.CallOption) (*NewSessionIDResponce, error)
	NewUser(ctx context.Context, in *NewUserRequest, opts ...grpc.CallOption) (*NewUserResponce, error)
	LoginUser(ctx context.Context, in *LoginUserRequest, opts ...grpc.CallOption) (*LoginUserResponce, error)
	UserData(ctx context.Context, in *UserDataRequest, opts ...grpc.CallOption) (*UserDataResponce, error)
	TimeStamp(ctx context.Context, in *TimeStampRequest, opts ...grpc.CallOption) (*TimeStampResponce, error)
	UpdateData(ctx context.Context, in *UpdateDataRequest, opts ...grpc.CallOption) (*UpdateDataResponce, error)
	LogOut(ctx context.Context, in *LogOutRequest, opts ...grpc.CallOption) (*LogOutResponce, error)
}

type gophKeeperClient struct {
	cc grpc.ClientConnInterface
}

func NewGophKeeperClient(cc grpc.ClientConnInterface) GophKeeperClient {
	return &gophKeeperClient{cc}
}

func (c *gophKeeperClient) NewSessionID(ctx context.Context, in *NewSessionIDRequest, opts ...grpc.CallOption) (*NewSessionIDResponce, error) {
	out := new(NewSessionIDResponce)
	err := c.cc.Invoke(ctx, GophKeeper_NewSessionID_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *gophKeeperClient) NewUser(ctx context.Context, in *NewUserRequest, opts ...grpc.CallOption) (*NewUserResponce, error) {
	out := new(NewUserResponce)
	err := c.cc.Invoke(ctx, GophKeeper_NewUser_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *gophKeeperClient) LoginUser(ctx context.Context, in *LoginUserRequest, opts ...grpc.CallOption) (*LoginUserResponce, error) {
	out := new(LoginUserResponce)
	err := c.cc.Invoke(ctx, GophKeeper_LoginUser_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *gophKeeperClient) UserData(ctx context.Context, in *UserDataRequest, opts ...grpc.CallOption) (*UserDataResponce, error) {
	out := new(UserDataResponce)
	err := c.cc.Invoke(ctx, GophKeeper_UserData_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *gophKeeperClient) TimeStamp(ctx context.Context, in *TimeStampRequest, opts ...grpc.CallOption) (*TimeStampResponce, error) {
	out := new(TimeStampResponce)
	err := c.cc.Invoke(ctx, GophKeeper_TimeStamp_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *gophKeeperClient) UpdateData(ctx context.Context, in *UpdateDataRequest, opts ...grpc.CallOption) (*UpdateDataResponce, error) {
	out := new(UpdateDataResponce)
	err := c.cc.Invoke(ctx, GophKeeper_UpdateData_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *gophKeeperClient) LogOut(ctx context.Context, in *LogOutRequest, opts ...grpc.CallOption) (*LogOutResponce, error) {
	out := new(LogOutResponce)
	err := c.cc.Invoke(ctx, GophKeeper_LogOut_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// GophKeeperServer is the server API for GophKeeper service.
// All implementations must embed UnimplementedGophKeeperServer
// for forward compatibility
type GophKeeperServer interface {
	NewSessionID(context.Context, *NewSessionIDRequest) (*NewSessionIDResponce, error)
	NewUser(context.Context, *NewUserRequest) (*NewUserResponce, error)
	LoginUser(context.Context, *LoginUserRequest) (*LoginUserResponce, error)
	UserData(context.Context, *UserDataRequest) (*UserDataResponce, error)
	TimeStamp(context.Context, *TimeStampRequest) (*TimeStampResponce, error)
	UpdateData(context.Context, *UpdateDataRequest) (*UpdateDataResponce, error)
	LogOut(context.Context, *LogOutRequest) (*LogOutResponce, error)
	mustEmbedUnimplementedGophKeeperServer()
}

// UnimplementedGophKeeperServer must be embedded to have forward compatible implementations.
type UnimplementedGophKeeperServer struct {
}

func (UnimplementedGophKeeperServer) NewSessionID(context.Context, *NewSessionIDRequest) (*NewSessionIDResponce, error) {
	return nil, status.Errorf(codes.Unimplemented, "method NewSessionID not implemented")
}
func (UnimplementedGophKeeperServer) NewUser(context.Context, *NewUserRequest) (*NewUserResponce, error) {
	return nil, status.Errorf(codes.Unimplemented, "method NewUser not implemented")
}
func (UnimplementedGophKeeperServer) LoginUser(context.Context, *LoginUserRequest) (*LoginUserResponce, error) {
	return nil, status.Errorf(codes.Unimplemented, "method LoginUser not implemented")
}
func (UnimplementedGophKeeperServer) UserData(context.Context, *UserDataRequest) (*UserDataResponce, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UserData not implemented")
}
func (UnimplementedGophKeeperServer) TimeStamp(context.Context, *TimeStampRequest) (*TimeStampResponce, error) {
	return nil, status.Errorf(codes.Unimplemented, "method TimeStamp not implemented")
}
func (UnimplementedGophKeeperServer) UpdateData(context.Context, *UpdateDataRequest) (*UpdateDataResponce, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateData not implemented")
}
func (UnimplementedGophKeeperServer) LogOut(context.Context, *LogOutRequest) (*LogOutResponce, error) {
	return nil, status.Errorf(codes.Unimplemented, "method LogOut not implemented")
}
func (UnimplementedGophKeeperServer) mustEmbedUnimplementedGophKeeperServer() {}

// UnsafeGophKeeperServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to GophKeeperServer will
// result in compilation errors.
type UnsafeGophKeeperServer interface {
	mustEmbedUnimplementedGophKeeperServer()
}

func RegisterGophKeeperServer(s grpc.ServiceRegistrar, srv GophKeeperServer) {
	s.RegisterService(&GophKeeper_ServiceDesc, srv)
}

func _GophKeeper_NewSessionID_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(NewSessionIDRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GophKeeperServer).NewSessionID(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: GophKeeper_NewSessionID_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GophKeeperServer).NewSessionID(ctx, req.(*NewSessionIDRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _GophKeeper_NewUser_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(NewUserRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GophKeeperServer).NewUser(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: GophKeeper_NewUser_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GophKeeperServer).NewUser(ctx, req.(*NewUserRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _GophKeeper_LoginUser_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(LoginUserRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GophKeeperServer).LoginUser(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: GophKeeper_LoginUser_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GophKeeperServer).LoginUser(ctx, req.(*LoginUserRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _GophKeeper_UserData_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UserDataRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GophKeeperServer).UserData(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: GophKeeper_UserData_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GophKeeperServer).UserData(ctx, req.(*UserDataRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _GophKeeper_TimeStamp_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(TimeStampRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GophKeeperServer).TimeStamp(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: GophKeeper_TimeStamp_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GophKeeperServer).TimeStamp(ctx, req.(*TimeStampRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _GophKeeper_UpdateData_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateDataRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GophKeeperServer).UpdateData(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: GophKeeper_UpdateData_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GophKeeperServer).UpdateData(ctx, req.(*UpdateDataRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _GophKeeper_LogOut_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(LogOutRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GophKeeperServer).LogOut(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: GophKeeper_LogOut_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GophKeeperServer).LogOut(ctx, req.(*LogOutRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// GophKeeper_ServiceDesc is the grpc.ServiceDesc for GophKeeper service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var GophKeeper_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "grpc.GophKeeper",
	HandlerType: (*GophKeeperServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "NewSessionID",
			Handler:    _GophKeeper_NewSessionID_Handler,
		},
		{
			MethodName: "NewUser",
			Handler:    _GophKeeper_NewUser_Handler,
		},
		{
			MethodName: "LoginUser",
			Handler:    _GophKeeper_LoginUser_Handler,
		},
		{
			MethodName: "UserData",
			Handler:    _GophKeeper_UserData_Handler,
		},
		{
			MethodName: "TimeStamp",
			Handler:    _GophKeeper_TimeStamp_Handler,
		},
		{
			MethodName: "UpdateData",
			Handler:    _GophKeeper_UpdateData_Handler,
		},
		{
			MethodName: "LogOut",
			Handler:    _GophKeeper_LogOut_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/grpc.proto",
}
