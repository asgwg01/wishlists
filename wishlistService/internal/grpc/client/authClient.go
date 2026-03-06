package client

import (
	"context"
	"fmt"
	authv1 "pkg/proto/auth/v1"
	"pkg/types/trace"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

type UserInfo struct {
	UserID string
	Email  string
	Name   string
}

type AuthClient struct {
	client authv1.AuthServiceClient
	conn   *grpc.ClientConn
}

func NewAuthClient(addr string, port string) (*AuthClient, error) {
	conn, err := grpc.NewClient(
		addr+":"+port,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
			ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
			defer cancel()
			return invoker(ctx, method, req, reply, cc, opts...)
		}),
	)

	if err != nil {
		return nil, fmt.Errorf("failed to connect to auth service: %w", err)
	}

	return &AuthClient{
		client: authv1.NewAuthServiceClient(conn),
		conn:   conn,
	}, nil
}

// GetUserInfo получает информацию о пользователе по ID
func (c *AuthClient) GetUserInfo(ctx context.Context, userID string) (*UserInfo, error) {
	ctx = trace.InjectIntoGRPC(ctx)

	resp, err := c.client.GetUserInfo(ctx, &authv1.UserInfoRequest{
		UserId: userID,
	})
	if err != nil {
		if st, ok := status.FromError(err); ok {
			switch st.Code() {
			case codes.NotFound:
				return nil, fmt.Errorf("user %s not found", userID)
			case codes.InvalidArgument:
				return nil, fmt.Errorf("invalid user id: %s", userID)
			}
		}
		return nil, fmt.Errorf("failed to get user info: %w", err)
	}

	return &UserInfo{
		UserID: resp.UserId,
		Email:  resp.Email,
		Name:   resp.Name,
	}, nil
}

func (c *AuthClient) Close() error {
	return c.conn.Close()
}
