// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: project_service.proto

package service

import (
	context "context"
	fmt "fmt"
	_ "github.com/gogo/protobuf/gogoproto"
	proto "github.com/gogo/protobuf/proto"
	project "github.com/uibricks/studio-engine/internal/pkg/proto/project"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion3 // please upgrade the proto package

func init() { proto.RegisterFile("project_service.proto", fileDescriptor_5169dba687284f3c) }

var fileDescriptor_5169dba687284f3c = []byte{
	// 366 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x8c, 0x93, 0xcd, 0x4a, 0xc3, 0x40,
	0x14, 0x85, 0x13, 0x10, 0x85, 0x29, 0xad, 0x10, 0x10, 0xb5, 0xca, 0xd4, 0x17, 0x68, 0x07, 0x74,
	0xeb, 0xaa, 0x15, 0x74, 0x51, 0x44, 0x5a, 0x75, 0xe1, 0x46, 0x9a, 0xf4, 0x1a, 0xa7, 0x3f, 0x33,
	0x71, 0x66, 0xd2, 0xe7, 0xf0, 0x31, 0x7c, 0x04, 0x97, 0x2e, 0xbb, 0xec, 0xd2, 0xa5, 0x4d, 0x5f,
	0xc0, 0xa5, 0x4b, 0x21, 0x33, 0x99, 0xfe, 0xa4, 0x68, 0x57, 0xc9, 0x3d, 0xe7, 0x9e, 0x8f, 0x13,
	0xc8, 0x45, 0x7b, 0x91, 0xe0, 0x3d, 0x08, 0xd4, 0xa3, 0x04, 0x31, 0xa2, 0x01, 0xd4, 0x22, 0xc1,
	0x15, 0xf7, 0x76, 0xcc, 0x58, 0xae, 0x86, 0x54, 0x3d, 0xc7, 0x7e, 0x2d, 0xe0, 0x43, 0x12, 0xf2,
	0x90, 0x93, 0xd4, 0xf7, 0xe3, 0xa7, 0x74, 0x4a, 0x87, 0xf4, 0x4d, 0xe7, 0xca, 0x19, 0x8e, 0x98,
	0xa7, 0x96, 0x4f, 0xbf, 0xb7, 0x50, 0xe9, 0x46, 0x2b, 0x6d, 0x0d, 0xf6, 0xea, 0xa8, 0x78, 0x17,
	0x75, 0x3b, 0x0a, 0x8c, 0xee, 0xed, 0xd7, 0xb2, 0x8c, 0x51, 0x5a, 0xf0, 0x12, 0x83, 0x54, 0xe5,
	0x83, 0xbc, 0x21, 0x23, 0xce, 0x24, 0x78, 0x57, 0xa8, 0xd0, 0xee, 0x8c, 0x2c, 0xe1, 0xc8, 0x2e,
	0x2e, 0xa8, 0x19, 0xe5, 0x78, 0xbd, 0x69, 0x48, 0x0d, 0x84, 0x2e, 0x41, 0xe5, 0x41, 0x46, 0x69,
	0xc6, 0xb4, 0xfb, 0x7f, 0x9d, 0x3a, 0x2a, 0xcc, 0x21, 0xd2, 0xcb, 0x2d, 0xca, 0x0c, 0x71, 0xb8,
	0xc6, 0x31, 0x8c, 0x26, 0x2a, 0x5e, 0xc0, 0x00, 0x14, 0x6c, 0xd4, 0x05, 0x5b, 0x73, 0x29, 0x64,
	0x69, 0x6d, 0xe4, 0xcd, 0x1b, 0xdd, 0x83, 0x90, 0x94, 0x33, 0xf9, 0x37, 0xf2, 0x64, 0xd5, 0xcc,
	0x62, 0x16, 0x7a, 0x8d, 0x4a, 0x2d, 0x90, 0x8a, 0x8b, 0xcd, 0x3a, 0x56, 0xac, 0xb9, 0x9c, 0xb2,
	0xbc, 0x5b, 0xb4, 0xab, 0xdb, 0x37, 0xf8, 0x30, 0xe2, 0x0c, 0x98, 0xf2, 0x2a, 0x2b, 0xdf, 0x65,
	0x9d, 0x7c, 0xcb, 0xdc, 0x82, 0xa6, 0xd6, 0x7b, 0xe3, 0x29, 0x76, 0x7e, 0xa6, 0xd8, 0x7d, 0x4b,
	0xb0, 0xfb, 0x9e, 0x60, 0xf7, 0x23, 0xc1, 0xee, 0x38, 0xc1, 0xee, 0x24, 0xc1, 0xee, 0x57, 0x82,
	0xdd, 0xd7, 0x19, 0x76, 0x26, 0x33, 0xec, 0x7c, 0xce, 0xb0, 0xf3, 0x70, 0xbe, 0xf0, 0x9b, 0xc7,
	0xd4, 0x17, 0x34, 0xe8, 0x4b, 0x22, 0x55, 0xdc, 0xa5, 0xbc, 0x0a, 0x2c, 0xa4, 0x0c, 0x08, 0x65,
	0x0a, 0x04, 0xeb, 0x0c, 0x48, 0xd4, 0x0f, 0xf5, 0x15, 0x10, 0x73, 0x24, 0xfe, 0x76, 0x3a, 0x9e,
	0xfd, 0x06, 0x00, 0x00, 0xff, 0xff, 0x99, 0x25, 0x62, 0xaf, 0x4d, 0x03, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// ProjectServiceClient is the client API for ProjectService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type ProjectServiceClient interface {
	// rpc unary
	// updates project for the specified project in the cache, if the project doesn't exist then create a new cache
	UpdateProject(ctx context.Context, in *project.ProjectRequest, opts ...grpc.CallOption) (*project.ProjectResponse, error)
	// rpc unary
	// persists project details to DB and clears it from cache
	SaveProject(ctx context.Context, in *project.SaveProjectRequest, opts ...grpc.CallOption) (*project.SaveProjectResponse, error)
	// rpc unary
	// gets project from cache/db by its luid
	GetProject(ctx context.Context, in *project.ProjectLuidRequest, opts ...grpc.CallOption) (*project.ProjectResponse, error)
	// rpc unary
	// gets all projects by its type and parentId
	GetProjects(ctx context.Context, in *project.ProjectsRequest, opts ...grpc.CallOption) (*project.ProjectsResponse, error)
	// rpc unary
	// deletes project from db and cache
	// permanently deletes project if marked deleted already, by moving it to history table
	DeleteProject(ctx context.Context, in *project.ProjectLuidRequest, opts ...grpc.CallOption) (*project.DeleteProjectResponse, error)
	// rpc unary
	GetProjectVersions(ctx context.Context, in *project.ProjectLuidRequest, opts ...grpc.CallOption) (*project.ProjectVersionsResponse, error)
	// rpc unary
	// restores a project by clearing deleted_at
	RestoreProject(ctx context.Context, in *project.ProjectLuidRequest, opts ...grpc.CallOption) (*project.RestoreProjectResponse, error)
	DeleteComponent(ctx context.Context, in *project.DeleteComponentRequest, opts ...grpc.CallOption) (*project.DeleteComponentResponse, error)
}

type projectServiceClient struct {
	cc *grpc.ClientConn
}

func NewProjectServiceClient(cc *grpc.ClientConn) ProjectServiceClient {
	return &projectServiceClient{cc}
}

func (c *projectServiceClient) UpdateProject(ctx context.Context, in *project.ProjectRequest, opts ...grpc.CallOption) (*project.ProjectResponse, error) {
	out := new(project.ProjectResponse)
	err := c.cc.Invoke(ctx, "/service.ProjectService/UpdateProject", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *projectServiceClient) SaveProject(ctx context.Context, in *project.SaveProjectRequest, opts ...grpc.CallOption) (*project.SaveProjectResponse, error) {
	out := new(project.SaveProjectResponse)
	err := c.cc.Invoke(ctx, "/service.ProjectService/SaveProject", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *projectServiceClient) GetProject(ctx context.Context, in *project.ProjectLuidRequest, opts ...grpc.CallOption) (*project.ProjectResponse, error) {
	out := new(project.ProjectResponse)
	err := c.cc.Invoke(ctx, "/service.ProjectService/GetProject", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *projectServiceClient) GetProjects(ctx context.Context, in *project.ProjectsRequest, opts ...grpc.CallOption) (*project.ProjectsResponse, error) {
	out := new(project.ProjectsResponse)
	err := c.cc.Invoke(ctx, "/service.ProjectService/GetProjects", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *projectServiceClient) DeleteProject(ctx context.Context, in *project.ProjectLuidRequest, opts ...grpc.CallOption) (*project.DeleteProjectResponse, error) {
	out := new(project.DeleteProjectResponse)
	err := c.cc.Invoke(ctx, "/service.ProjectService/DeleteProject", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *projectServiceClient) GetProjectVersions(ctx context.Context, in *project.ProjectLuidRequest, opts ...grpc.CallOption) (*project.ProjectVersionsResponse, error) {
	out := new(project.ProjectVersionsResponse)
	err := c.cc.Invoke(ctx, "/service.ProjectService/GetProjectVersions", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *projectServiceClient) RestoreProject(ctx context.Context, in *project.ProjectLuidRequest, opts ...grpc.CallOption) (*project.RestoreProjectResponse, error) {
	out := new(project.RestoreProjectResponse)
	err := c.cc.Invoke(ctx, "/service.ProjectService/RestoreProject", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *projectServiceClient) DeleteComponent(ctx context.Context, in *project.DeleteComponentRequest, opts ...grpc.CallOption) (*project.DeleteComponentResponse, error) {
	out := new(project.DeleteComponentResponse)
	err := c.cc.Invoke(ctx, "/service.ProjectService/DeleteComponent", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ProjectServiceServer is the server API for ProjectService service.
type ProjectServiceServer interface {
	// rpc unary
	// updates project for the specified project in the cache, if the project doesn't exist then create a new cache
	UpdateProject(context.Context, *project.ProjectRequest) (*project.ProjectResponse, error)
	// rpc unary
	// persists project details to DB and clears it from cache
	SaveProject(context.Context, *project.SaveProjectRequest) (*project.SaveProjectResponse, error)
	// rpc unary
	// gets project from cache/db by its luid
	GetProject(context.Context, *project.ProjectLuidRequest) (*project.ProjectResponse, error)
	// rpc unary
	// gets all projects by its type and parentId
	GetProjects(context.Context, *project.ProjectsRequest) (*project.ProjectsResponse, error)
	// rpc unary
	// deletes project from db and cache
	// permanently deletes project if marked deleted already, by moving it to history table
	DeleteProject(context.Context, *project.ProjectLuidRequest) (*project.DeleteProjectResponse, error)
	// rpc unary
	GetProjectVersions(context.Context, *project.ProjectLuidRequest) (*project.ProjectVersionsResponse, error)
	// rpc unary
	// restores a project by clearing deleted_at
	RestoreProject(context.Context, *project.ProjectLuidRequest) (*project.RestoreProjectResponse, error)
	DeleteComponent(context.Context, *project.DeleteComponentRequest) (*project.DeleteComponentResponse, error)
}

// UnimplementedProjectServiceServer can be embedded to have forward compatible implementations.
type UnimplementedProjectServiceServer struct {
}

func (*UnimplementedProjectServiceServer) UpdateProject(ctx context.Context, req *project.ProjectRequest) (*project.ProjectResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateProject not implemented")
}
func (*UnimplementedProjectServiceServer) SaveProject(ctx context.Context, req *project.SaveProjectRequest) (*project.SaveProjectResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SaveProject not implemented")
}
func (*UnimplementedProjectServiceServer) GetProject(ctx context.Context, req *project.ProjectLuidRequest) (*project.ProjectResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetProject not implemented")
}
func (*UnimplementedProjectServiceServer) GetProjects(ctx context.Context, req *project.ProjectsRequest) (*project.ProjectsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetProjects not implemented")
}
func (*UnimplementedProjectServiceServer) DeleteProject(ctx context.Context, req *project.ProjectLuidRequest) (*project.DeleteProjectResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteProject not implemented")
}
func (*UnimplementedProjectServiceServer) GetProjectVersions(ctx context.Context, req *project.ProjectLuidRequest) (*project.ProjectVersionsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetProjectVersions not implemented")
}
func (*UnimplementedProjectServiceServer) RestoreProject(ctx context.Context, req *project.ProjectLuidRequest) (*project.RestoreProjectResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RestoreProject not implemented")
}
func (*UnimplementedProjectServiceServer) DeleteComponent(ctx context.Context, req *project.DeleteComponentRequest) (*project.DeleteComponentResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteComponent not implemented")
}

func RegisterProjectServiceServer(s *grpc.Server, srv ProjectServiceServer) {
	s.RegisterService(&_ProjectService_serviceDesc, srv)
}

func _ProjectService_UpdateProject_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(project.ProjectRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ProjectServiceServer).UpdateProject(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/service.ProjectService/UpdateProject",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ProjectServiceServer).UpdateProject(ctx, req.(*project.ProjectRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ProjectService_SaveProject_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(project.SaveProjectRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ProjectServiceServer).SaveProject(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/service.ProjectService/SaveProject",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ProjectServiceServer).SaveProject(ctx, req.(*project.SaveProjectRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ProjectService_GetProject_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(project.ProjectLuidRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ProjectServiceServer).GetProject(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/service.ProjectService/GetProject",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ProjectServiceServer).GetProject(ctx, req.(*project.ProjectLuidRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ProjectService_GetProjects_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(project.ProjectsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ProjectServiceServer).GetProjects(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/service.ProjectService/GetProjects",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ProjectServiceServer).GetProjects(ctx, req.(*project.ProjectsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ProjectService_DeleteProject_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(project.ProjectLuidRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ProjectServiceServer).DeleteProject(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/service.ProjectService/DeleteProject",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ProjectServiceServer).DeleteProject(ctx, req.(*project.ProjectLuidRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ProjectService_GetProjectVersions_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(project.ProjectLuidRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ProjectServiceServer).GetProjectVersions(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/service.ProjectService/GetProjectVersions",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ProjectServiceServer).GetProjectVersions(ctx, req.(*project.ProjectLuidRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ProjectService_RestoreProject_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(project.ProjectLuidRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ProjectServiceServer).RestoreProject(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/service.ProjectService/RestoreProject",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ProjectServiceServer).RestoreProject(ctx, req.(*project.ProjectLuidRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ProjectService_DeleteComponent_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(project.DeleteComponentRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ProjectServiceServer).DeleteComponent(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/service.ProjectService/DeleteComponent",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ProjectServiceServer).DeleteComponent(ctx, req.(*project.DeleteComponentRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _ProjectService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "service.ProjectService",
	HandlerType: (*ProjectServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "UpdateProject",
			Handler:    _ProjectService_UpdateProject_Handler,
		},
		{
			MethodName: "SaveProject",
			Handler:    _ProjectService_SaveProject_Handler,
		},
		{
			MethodName: "GetProject",
			Handler:    _ProjectService_GetProject_Handler,
		},
		{
			MethodName: "GetProjects",
			Handler:    _ProjectService_GetProjects_Handler,
		},
		{
			MethodName: "DeleteProject",
			Handler:    _ProjectService_DeleteProject_Handler,
		},
		{
			MethodName: "GetProjectVersions",
			Handler:    _ProjectService_GetProjectVersions_Handler,
		},
		{
			MethodName: "RestoreProject",
			Handler:    _ProjectService_RestoreProject_Handler,
		},
		{
			MethodName: "DeleteComponent",
			Handler:    _ProjectService_DeleteComponent_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "project_service.proto",
}
