// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v5.28.0--rc2
// source: pkg/pb/chat.proto

package pb

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.64.0 or later.
const _ = grpc.SupportPackageIsVersion9

const (
	ChatNCallService_GetOneToOneChats_FullMethodName        = "/chatNcall_proto.ChatNCallService/GetOneToOneChats"
	ChatNCallService_GetRecentChatProfiles_FullMethodName   = "/chatNcall_proto.ChatNCallService/GetRecentChatProfiles"
	ChatNCallService_GetGroupMembersInfo_FullMethodName     = "/chatNcall_proto.ChatNCallService/GetGroupMembersInfo"
	ChatNCallService_CreateNewGroup_FullMethodName          = "/chatNcall_proto.ChatNCallService/CreateNewGroup"
	ChatNCallService_GetUserGroupChatSummary_FullMethodName = "/chatNcall_proto.ChatNCallService/GetUserGroupChatSummary"
	ChatNCallService_GetOneToManyChats_FullMethodName       = "/chatNcall_proto.ChatNCallService/GetOneToManyChats"
	ChatNCallService_AddMembersToGroup_FullMethodName       = "/chatNcall_proto.ChatNCallService/AddMembersToGroup"
	ChatNCallService_RemoveMemberFromGroup_FullMethodName   = "/chatNcall_proto.ChatNCallService/RemoveMemberFromGroup"
	ChatNCallService_CreateRoom_FullMethodName              = "/chatNcall_proto.ChatNCallService/CreateRoom"
	ChatNCallService_GetRoom_FullMethodName                 = "/chatNcall_proto.ChatNCallService/GetRoom"
	ChatNCallService_InsertIntoRoom_FullMethodName          = "/chatNcall_proto.ChatNCallService/InsertIntoRoom"
	ChatNCallService_DeleteRoom_FullMethodName              = "/chatNcall_proto.ChatNCallService/DeleteRoom"
)

// ChatNCallServiceClient is the client API for ChatNCallService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ChatNCallServiceClient interface {
	GetOneToOneChats(ctx context.Context, in *RequestUserOneToOneChat, opts ...grpc.CallOption) (*ResponseUserOneToOneChat, error)
	GetRecentChatProfiles(ctx context.Context, in *RequestRecentChatProfiles, opts ...grpc.CallOption) (*ResponseRecentChatProfiles, error)
	GetGroupMembersInfo(ctx context.Context, in *RequestGroupMembersInfo, opts ...grpc.CallOption) (*ResponseGroupMembersInfo, error)
	CreateNewGroup(ctx context.Context, in *RequestNewGroup, opts ...grpc.CallOption) (*ResponseNewGroup, error)
	GetUserGroupChatSummary(ctx context.Context, in *RequestGroupChatSummary, opts ...grpc.CallOption) (*ResponseGroupChatSummary, error)
	GetOneToManyChats(ctx context.Context, in *RequestGetOneToManyChats, opts ...grpc.CallOption) (*ResponseGetOneToManyChats, error)
	AddMembersToGroup(ctx context.Context, in *RequestAddGroupMembers, opts ...grpc.CallOption) (*ResponseAddGroupMembers, error)
	RemoveMemberFromGroup(ctx context.Context, in *RequestRemoveGroupMember, opts ...grpc.CallOption) (*ResponseRemoveGroupMember, error)
	CreateRoom(ctx context.Context, in *CreateRoomRequest, opts ...grpc.CallOption) (*CreateRoomResponse, error)
	GetRoom(ctx context.Context, in *GetRoomRequest, opts ...grpc.CallOption) (*GetRoomResponse, error)
	InsertIntoRoom(ctx context.Context, in *InsertIntoRoomRequest, opts ...grpc.CallOption) (*InsertIntoRoomResponse, error)
	DeleteRoom(ctx context.Context, in *DeleteRoomRequest, opts ...grpc.CallOption) (*DeleteRoomResponse, error)
}

type chatNCallServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewChatNCallServiceClient(cc grpc.ClientConnInterface) ChatNCallServiceClient {
	return &chatNCallServiceClient{cc}
}

func (c *chatNCallServiceClient) GetOneToOneChats(ctx context.Context, in *RequestUserOneToOneChat, opts ...grpc.CallOption) (*ResponseUserOneToOneChat, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ResponseUserOneToOneChat)
	err := c.cc.Invoke(ctx, ChatNCallService_GetOneToOneChats_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *chatNCallServiceClient) GetRecentChatProfiles(ctx context.Context, in *RequestRecentChatProfiles, opts ...grpc.CallOption) (*ResponseRecentChatProfiles, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ResponseRecentChatProfiles)
	err := c.cc.Invoke(ctx, ChatNCallService_GetRecentChatProfiles_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *chatNCallServiceClient) GetGroupMembersInfo(ctx context.Context, in *RequestGroupMembersInfo, opts ...grpc.CallOption) (*ResponseGroupMembersInfo, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ResponseGroupMembersInfo)
	err := c.cc.Invoke(ctx, ChatNCallService_GetGroupMembersInfo_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *chatNCallServiceClient) CreateNewGroup(ctx context.Context, in *RequestNewGroup, opts ...grpc.CallOption) (*ResponseNewGroup, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ResponseNewGroup)
	err := c.cc.Invoke(ctx, ChatNCallService_CreateNewGroup_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *chatNCallServiceClient) GetUserGroupChatSummary(ctx context.Context, in *RequestGroupChatSummary, opts ...grpc.CallOption) (*ResponseGroupChatSummary, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ResponseGroupChatSummary)
	err := c.cc.Invoke(ctx, ChatNCallService_GetUserGroupChatSummary_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *chatNCallServiceClient) GetOneToManyChats(ctx context.Context, in *RequestGetOneToManyChats, opts ...grpc.CallOption) (*ResponseGetOneToManyChats, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ResponseGetOneToManyChats)
	err := c.cc.Invoke(ctx, ChatNCallService_GetOneToManyChats_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *chatNCallServiceClient) AddMembersToGroup(ctx context.Context, in *RequestAddGroupMembers, opts ...grpc.CallOption) (*ResponseAddGroupMembers, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ResponseAddGroupMembers)
	err := c.cc.Invoke(ctx, ChatNCallService_AddMembersToGroup_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *chatNCallServiceClient) RemoveMemberFromGroup(ctx context.Context, in *RequestRemoveGroupMember, opts ...grpc.CallOption) (*ResponseRemoveGroupMember, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ResponseRemoveGroupMember)
	err := c.cc.Invoke(ctx, ChatNCallService_RemoveMemberFromGroup_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *chatNCallServiceClient) CreateRoom(ctx context.Context, in *CreateRoomRequest, opts ...grpc.CallOption) (*CreateRoomResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(CreateRoomResponse)
	err := c.cc.Invoke(ctx, ChatNCallService_CreateRoom_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *chatNCallServiceClient) GetRoom(ctx context.Context, in *GetRoomRequest, opts ...grpc.CallOption) (*GetRoomResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetRoomResponse)
	err := c.cc.Invoke(ctx, ChatNCallService_GetRoom_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *chatNCallServiceClient) InsertIntoRoom(ctx context.Context, in *InsertIntoRoomRequest, opts ...grpc.CallOption) (*InsertIntoRoomResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(InsertIntoRoomResponse)
	err := c.cc.Invoke(ctx, ChatNCallService_InsertIntoRoom_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *chatNCallServiceClient) DeleteRoom(ctx context.Context, in *DeleteRoomRequest, opts ...grpc.CallOption) (*DeleteRoomResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(DeleteRoomResponse)
	err := c.cc.Invoke(ctx, ChatNCallService_DeleteRoom_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ChatNCallServiceServer is the server API for ChatNCallService service.
// All implementations must embed UnimplementedChatNCallServiceServer
// for forward compatibility.
type ChatNCallServiceServer interface {
	GetOneToOneChats(context.Context, *RequestUserOneToOneChat) (*ResponseUserOneToOneChat, error)
	GetRecentChatProfiles(context.Context, *RequestRecentChatProfiles) (*ResponseRecentChatProfiles, error)
	GetGroupMembersInfo(context.Context, *RequestGroupMembersInfo) (*ResponseGroupMembersInfo, error)
	CreateNewGroup(context.Context, *RequestNewGroup) (*ResponseNewGroup, error)
	GetUserGroupChatSummary(context.Context, *RequestGroupChatSummary) (*ResponseGroupChatSummary, error)
	GetOneToManyChats(context.Context, *RequestGetOneToManyChats) (*ResponseGetOneToManyChats, error)
	AddMembersToGroup(context.Context, *RequestAddGroupMembers) (*ResponseAddGroupMembers, error)
	RemoveMemberFromGroup(context.Context, *RequestRemoveGroupMember) (*ResponseRemoveGroupMember, error)
	CreateRoom(context.Context, *CreateRoomRequest) (*CreateRoomResponse, error)
	GetRoom(context.Context, *GetRoomRequest) (*GetRoomResponse, error)
	InsertIntoRoom(context.Context, *InsertIntoRoomRequest) (*InsertIntoRoomResponse, error)
	DeleteRoom(context.Context, *DeleteRoomRequest) (*DeleteRoomResponse, error)
	mustEmbedUnimplementedChatNCallServiceServer()
}

// UnimplementedChatNCallServiceServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedChatNCallServiceServer struct{}

func (UnimplementedChatNCallServiceServer) GetOneToOneChats(context.Context, *RequestUserOneToOneChat) (*ResponseUserOneToOneChat, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetOneToOneChats not implemented")
}
func (UnimplementedChatNCallServiceServer) GetRecentChatProfiles(context.Context, *RequestRecentChatProfiles) (*ResponseRecentChatProfiles, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetRecentChatProfiles not implemented")
}
func (UnimplementedChatNCallServiceServer) GetGroupMembersInfo(context.Context, *RequestGroupMembersInfo) (*ResponseGroupMembersInfo, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetGroupMembersInfo not implemented")
}
func (UnimplementedChatNCallServiceServer) CreateNewGroup(context.Context, *RequestNewGroup) (*ResponseNewGroup, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateNewGroup not implemented")
}
func (UnimplementedChatNCallServiceServer) GetUserGroupChatSummary(context.Context, *RequestGroupChatSummary) (*ResponseGroupChatSummary, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetUserGroupChatSummary not implemented")
}
func (UnimplementedChatNCallServiceServer) GetOneToManyChats(context.Context, *RequestGetOneToManyChats) (*ResponseGetOneToManyChats, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetOneToManyChats not implemented")
}
func (UnimplementedChatNCallServiceServer) AddMembersToGroup(context.Context, *RequestAddGroupMembers) (*ResponseAddGroupMembers, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AddMembersToGroup not implemented")
}
func (UnimplementedChatNCallServiceServer) RemoveMemberFromGroup(context.Context, *RequestRemoveGroupMember) (*ResponseRemoveGroupMember, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RemoveMemberFromGroup not implemented")
}
func (UnimplementedChatNCallServiceServer) CreateRoom(context.Context, *CreateRoomRequest) (*CreateRoomResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateRoom not implemented")
}
func (UnimplementedChatNCallServiceServer) GetRoom(context.Context, *GetRoomRequest) (*GetRoomResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetRoom not implemented")
}
func (UnimplementedChatNCallServiceServer) InsertIntoRoom(context.Context, *InsertIntoRoomRequest) (*InsertIntoRoomResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method InsertIntoRoom not implemented")
}
func (UnimplementedChatNCallServiceServer) DeleteRoom(context.Context, *DeleteRoomRequest) (*DeleteRoomResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteRoom not implemented")
}
func (UnimplementedChatNCallServiceServer) mustEmbedUnimplementedChatNCallServiceServer() {}
func (UnimplementedChatNCallServiceServer) testEmbeddedByValue()                          {}

// UnsafeChatNCallServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ChatNCallServiceServer will
// result in compilation errors.
type UnsafeChatNCallServiceServer interface {
	mustEmbedUnimplementedChatNCallServiceServer()
}

func RegisterChatNCallServiceServer(s grpc.ServiceRegistrar, srv ChatNCallServiceServer) {
	// If the following call pancis, it indicates UnimplementedChatNCallServiceServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&ChatNCallService_ServiceDesc, srv)
}

func _ChatNCallService_GetOneToOneChats_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RequestUserOneToOneChat)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ChatNCallServiceServer).GetOneToOneChats(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ChatNCallService_GetOneToOneChats_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ChatNCallServiceServer).GetOneToOneChats(ctx, req.(*RequestUserOneToOneChat))
	}
	return interceptor(ctx, in, info, handler)
}

func _ChatNCallService_GetRecentChatProfiles_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RequestRecentChatProfiles)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ChatNCallServiceServer).GetRecentChatProfiles(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ChatNCallService_GetRecentChatProfiles_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ChatNCallServiceServer).GetRecentChatProfiles(ctx, req.(*RequestRecentChatProfiles))
	}
	return interceptor(ctx, in, info, handler)
}

func _ChatNCallService_GetGroupMembersInfo_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RequestGroupMembersInfo)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ChatNCallServiceServer).GetGroupMembersInfo(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ChatNCallService_GetGroupMembersInfo_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ChatNCallServiceServer).GetGroupMembersInfo(ctx, req.(*RequestGroupMembersInfo))
	}
	return interceptor(ctx, in, info, handler)
}

func _ChatNCallService_CreateNewGroup_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RequestNewGroup)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ChatNCallServiceServer).CreateNewGroup(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ChatNCallService_CreateNewGroup_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ChatNCallServiceServer).CreateNewGroup(ctx, req.(*RequestNewGroup))
	}
	return interceptor(ctx, in, info, handler)
}

func _ChatNCallService_GetUserGroupChatSummary_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RequestGroupChatSummary)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ChatNCallServiceServer).GetUserGroupChatSummary(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ChatNCallService_GetUserGroupChatSummary_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ChatNCallServiceServer).GetUserGroupChatSummary(ctx, req.(*RequestGroupChatSummary))
	}
	return interceptor(ctx, in, info, handler)
}

func _ChatNCallService_GetOneToManyChats_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RequestGetOneToManyChats)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ChatNCallServiceServer).GetOneToManyChats(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ChatNCallService_GetOneToManyChats_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ChatNCallServiceServer).GetOneToManyChats(ctx, req.(*RequestGetOneToManyChats))
	}
	return interceptor(ctx, in, info, handler)
}

func _ChatNCallService_AddMembersToGroup_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RequestAddGroupMembers)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ChatNCallServiceServer).AddMembersToGroup(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ChatNCallService_AddMembersToGroup_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ChatNCallServiceServer).AddMembersToGroup(ctx, req.(*RequestAddGroupMembers))
	}
	return interceptor(ctx, in, info, handler)
}

func _ChatNCallService_RemoveMemberFromGroup_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RequestRemoveGroupMember)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ChatNCallServiceServer).RemoveMemberFromGroup(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ChatNCallService_RemoveMemberFromGroup_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ChatNCallServiceServer).RemoveMemberFromGroup(ctx, req.(*RequestRemoveGroupMember))
	}
	return interceptor(ctx, in, info, handler)
}

func _ChatNCallService_CreateRoom_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateRoomRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ChatNCallServiceServer).CreateRoom(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ChatNCallService_CreateRoom_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ChatNCallServiceServer).CreateRoom(ctx, req.(*CreateRoomRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ChatNCallService_GetRoom_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetRoomRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ChatNCallServiceServer).GetRoom(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ChatNCallService_GetRoom_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ChatNCallServiceServer).GetRoom(ctx, req.(*GetRoomRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ChatNCallService_InsertIntoRoom_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(InsertIntoRoomRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ChatNCallServiceServer).InsertIntoRoom(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ChatNCallService_InsertIntoRoom_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ChatNCallServiceServer).InsertIntoRoom(ctx, req.(*InsertIntoRoomRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ChatNCallService_DeleteRoom_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteRoomRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ChatNCallServiceServer).DeleteRoom(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ChatNCallService_DeleteRoom_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ChatNCallServiceServer).DeleteRoom(ctx, req.(*DeleteRoomRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// ChatNCallService_ServiceDesc is the grpc.ServiceDesc for ChatNCallService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var ChatNCallService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "chatNcall_proto.ChatNCallService",
	HandlerType: (*ChatNCallServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetOneToOneChats",
			Handler:    _ChatNCallService_GetOneToOneChats_Handler,
		},
		{
			MethodName: "GetRecentChatProfiles",
			Handler:    _ChatNCallService_GetRecentChatProfiles_Handler,
		},
		{
			MethodName: "GetGroupMembersInfo",
			Handler:    _ChatNCallService_GetGroupMembersInfo_Handler,
		},
		{
			MethodName: "CreateNewGroup",
			Handler:    _ChatNCallService_CreateNewGroup_Handler,
		},
		{
			MethodName: "GetUserGroupChatSummary",
			Handler:    _ChatNCallService_GetUserGroupChatSummary_Handler,
		},
		{
			MethodName: "GetOneToManyChats",
			Handler:    _ChatNCallService_GetOneToManyChats_Handler,
		},
		{
			MethodName: "AddMembersToGroup",
			Handler:    _ChatNCallService_AddMembersToGroup_Handler,
		},
		{
			MethodName: "RemoveMemberFromGroup",
			Handler:    _ChatNCallService_RemoveMemberFromGroup_Handler,
		},
		{
			MethodName: "CreateRoom",
			Handler:    _ChatNCallService_CreateRoom_Handler,
		},
		{
			MethodName: "GetRoom",
			Handler:    _ChatNCallService_GetRoom_Handler,
		},
		{
			MethodName: "InsertIntoRoom",
			Handler:    _ChatNCallService_InsertIntoRoom_Handler,
		},
		{
			MethodName: "DeleteRoom",
			Handler:    _ChatNCallService_DeleteRoom_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "pkg/pb/chat.proto",
}