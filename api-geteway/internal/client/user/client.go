package user

import (
	"context"
	"fmt"
	"time"

	"github.com/Lama189/soundwave-platform/api-geteway/internal/delivery/http"
	userpb "github.com/Lama189/soundwave-platform/gen/go/user/v1"
	"github.com/google/uuid"
	"google.golang.org/grpc"
)

type client struct {
	grpcClient userpb.UserServiceClient
	timeout    time.Duration
}

func NewClient(conn grpc.ClientConnInterface, defaultTimeout time.Duration) Client {
	if defaultTimeout <= 0 {
		defaultTimeout = 5 * time.Second
	}

	return &client{
		grpcClient: userpb.NewUserServiceClient(conn),
		timeout:    defaultTimeout,
	}
}

func (c *client) Register(ctx context.Context, req http.RegisterRequest) (*http.UserResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	pdReq := &userpb.RegisterRequest{
		Phone:    req.Phone,
		Password: req.Password,
	}

	pbRes, err := c.grpcClient.Register(ctx, pdReq)
	if err != nil {
		return nil, err
	}

	userID, err := uuid.Parse(pbRes.GetUser().GetId())
	if err != nil {
		return nil, fmt.Errorf("invalid user id returned from gRPC: %w", err)
	}

	return &http.UserResponse{
		ID:        userID,
		Phone:     pbRes.GetUser().GetPhone(),
		CreatedAt: pbRes.GetUser().CreatedAt.AsTime(),
	}, nil
}

func (c *client) Login(ctx context.Context, req http.LoginRequest) (*http.AuthResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	pbReq := &userpb.LoginRequest{
		Phone:    req.Phone,
		Password: req.Password,
	}

	pbRes, err := c.grpcClient.Login(ctx, pbReq)
	if err != nil {
		return nil, err
	}

	userID, err := uuid.Parse(pbRes.GetUser().GetId())
	if err != nil {
		return nil, fmt.Errorf("invalid user id returned from gRPC: %w", err)
	}

	return &http.AuthResponse{
		User: http.UserResponse{
			ID:        userID,
			Phone:     pbRes.GetUser().GetPhone(),
			CreatedAt: pbRes.GetUser().CreatedAt.AsTime(),
		},
		AccessToken:  pbRes.GetAccessToken(),
		RefreshToken: pbRes.GetRefreshToken(),
	}, nil
}

func (c *client) Refresh(ctx context.Context, req http.RefreshRequest) (*http.RefreshResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	pbReq := &userpb.RefreshRequest{
		RefreshToken: req.RefreshToken,
	}

	pbRes, err := c.grpcClient.Refresh(ctx, pbReq)
	if err != nil {
		return nil, err
	}

	return &http.RefreshResponse{
		AccessToken:  pbRes.GetAccessToken(),
		RefreshToken: pbRes.GetRefreshToken(),
	}, nil
}

func (c *client) GetMe(ctx context.Context, userID uuid.UUID) (*http.UserResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	pbReq := &userpb.GetUserRequest{
		Id: userID.String(),
	}

	pbRes, err := c.grpcClient.GetUser(ctx, pbReq)
	if err != nil {
		return nil, err
	}

	parsedID, err := uuid.Parse(pbRes.GetUser().GetId())
	if err != nil {
		return nil, fmt.Errorf("invalid user id returned from gRPC: %w", err)
	}

	return &http.UserResponse{
		ID:        parsedID,
		Phone:     pbRes.GetUser().GetPhone(),
		CreatedAt: pbRes.GetUser().GetCreatedAt().AsTime(),
	}, nil
}
