package v1

import (
	"context"

	"google.golang.org/grpc"
)

const (
	UserServiceName          = "cdek.user.v1.UserService"
	GetCurrentUserMethodName = "/cdek.user.v1.UserService/GetCurrentUser"
	BatchGetUsersMethodName  = "/cdek.user.v1.UserService/BatchGetUsers"
)

type UserSummary struct {
	Id        string `json:"id"`
	Name      string `json:"name"`
	Title     string `json:"title"`
	Company   string `json:"company"`
	Level     int32  `json:"level"`
	LevelText string `json:"levelText"`
	JoinedAt  string `json:"joinedAt"`
	Location  string `json:"location"`
	Team      string `json:"team"`
}

type GetCurrentUserRequest struct {
	UserId string `json:"userId"`
}

type GetCurrentUserResponse struct {
	User *UserSummary `json:"user"`
}

type BatchGetUsersRequest struct {
	UserIds []string `json:"userIds"`
}

type BatchGetUsersResponse struct {
	Users []*UserSummary `json:"users"`
}

type UserServiceClient interface {
	GetCurrentUser(ctx context.Context, in *GetCurrentUserRequest, opts ...grpc.CallOption) (*GetCurrentUserResponse, error)
	BatchGetUsers(ctx context.Context, in *BatchGetUsersRequest, opts ...grpc.CallOption) (*BatchGetUsersResponse, error)
}

type userServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewUserServiceClient(cc grpc.ClientConnInterface) UserServiceClient {
	return &userServiceClient{cc: cc}
}

func (c *userServiceClient) GetCurrentUser(ctx context.Context, in *GetCurrentUserRequest, opts ...grpc.CallOption) (*GetCurrentUserResponse, error) {
	out := new(GetCurrentUserResponse)
	if err := c.cc.Invoke(ctx, GetCurrentUserMethodName, in, out, opts...); err != nil {
		return nil, err
	}

	return out, nil
}

func (c *userServiceClient) BatchGetUsers(ctx context.Context, in *BatchGetUsersRequest, opts ...grpc.CallOption) (*BatchGetUsersResponse, error) {
	out := new(BatchGetUsersResponse)
	if err := c.cc.Invoke(ctx, BatchGetUsersMethodName, in, out, opts...); err != nil {
		return nil, err
	}

	return out, nil
}

type UserServiceServer interface {
	GetCurrentUser(context.Context, *GetCurrentUserRequest) (*GetCurrentUserResponse, error)
	BatchGetUsers(context.Context, *BatchGetUsersRequest) (*BatchGetUsersResponse, error)
}

func RegisterUserServiceServer(s grpc.ServiceRegistrar, srv UserServiceServer) {
	s.RegisterService(&grpc.ServiceDesc{
		ServiceName: UserServiceName,
		HandlerType: (*UserServiceServer)(nil),
		Methods: []grpc.MethodDesc{
			{
				MethodName: "GetCurrentUser",
				Handler:    getCurrentUserHandler(srv),
			},
			{
				MethodName: "BatchGetUsers",
				Handler:    batchGetUsersHandler(srv),
			},
		},
		Streams:  []grpc.StreamDesc{},
		Metadata: "proto/user/v1/user.proto",
	}, srv)
}

func getCurrentUserHandler(srv UserServiceServer) grpc.MethodHandler {
	return func(service any, ctx context.Context, dec func(any) error, interceptor grpc.UnaryServerInterceptor) (any, error) {
		in := new(GetCurrentUserRequest)
		if err := dec(in); err != nil {
			return nil, err
		}

		if interceptor == nil {
			return srv.GetCurrentUser(ctx, in)
		}

		info := &grpc.UnaryServerInfo{
			Server:     service,
			FullMethod: GetCurrentUserMethodName,
		}

		handler := func(ctx context.Context, req any) (any, error) {
			return srv.GetCurrentUser(ctx, req.(*GetCurrentUserRequest))
		}

		return interceptor(ctx, in, info, handler)
	}
}

func batchGetUsersHandler(srv UserServiceServer) grpc.MethodHandler {
	return func(service any, ctx context.Context, dec func(any) error, interceptor grpc.UnaryServerInterceptor) (any, error) {
		in := new(BatchGetUsersRequest)
		if err := dec(in); err != nil {
			return nil, err
		}

		if interceptor == nil {
			return srv.BatchGetUsers(ctx, in)
		}

		info := &grpc.UnaryServerInfo{
			Server:     service,
			FullMethod: BatchGetUsersMethodName,
		}

		handler := func(ctx context.Context, req any) (any, error) {
			return srv.BatchGetUsers(ctx, req.(*BatchGetUsersRequest))
		}

		return interceptor(ctx, in, info, handler)
	}
}
