package v1

import (
	"context"
	"errors"

	userpb "github.com/Lama189/ecommerce-core/gen/go/user/v1"
	"github.com/Lama189/ecommerce-core/user-service/internal/domain"
	"github.com/Lama189/ecommerce-core/user-service/internal/service/user"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type UserGRPCServer struct {
	userpb.UnimplementedUserServiceServer
	service UserService
}

func NewUserGRPCServer(service UserService) *UserGRPCServer {
	return &UserGRPCServer{
		service: service,
	}
}

func (s *UserGRPCServer) Register(ctx context.Context, req *userpb.RegisterRequest) (*userpb.RegisterResponse, error) {
	if req.GetPhone() == "" || req.GetPassword() == "" {
		return nil, status.Error(codes.InvalidArgument, "phone and password are required")
	}

	u, err := s.service.Create(ctx, user.CreateUserDTO{
		Phone:    req.GetPhone(),
		Password: req.GetPassword(),
	})
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrConflict):
			return nil, status.Error(codes.AlreadyExists, "user with this phone already exists")
		case errors.Is(err, domain.ErrInvalidInput):
			return nil, status.Error(codes.InvalidArgument, err.Error())
		default:
			return nil, status.Error(codes.Internal, "internal server error")
		}
	}

	return &userpb.RegisterResponse{
		User: &userpb.User{
			Id:        u.ID.String(),
			Phone:     u.Phone,
			CreatedAt: timestamppb.New(u.CreatedAt),
		},
	}, nil
}

func (s *UserGRPCServer) Login(ctx context.Context, req *userpb.LoginRequest) (*userpb.LoginResponse, error) {
	if req.GetPhone() == "" || req.GetPassword() == "" {
		return nil, status.Error(codes.InvalidArgument, "phone and password are required")
	}

	res, err := s.service.Login(ctx, req.GetPhone(), req.GetPassword())
	if err != nil {
		if errors.Is(err, user.ErrInvalidCredentials) {
			return nil, status.Error(codes.Unauthenticated, "invalid phone or password")
		}
		return nil, status.Error(codes.Internal, "internal server error")
	}

	return &userpb.LoginResponse{
		User: &userpb.User{
			Id:        res.User.ID.String(),
			Phone:     res.User.Phone,
			CreatedAt: timestamppb.New(res.User.CreatedAt),
		},
		AccessToken:  res.Tokens.AccessToken,
		RefreshToken: res.Tokens.RefreshToken,
	}, nil
}

func (s *UserGRPCServer) GetUser(ctx context.Context, req *userpb.GetUserRequest) (*userpb.GetUserResponse, error) {
	userID, err := uuid.Parse(req.GetId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid user id format")
	}

	u, err := s.service.GetByID(ctx, userID)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return nil, status.Error(codes.NotFound, "user not found")
		}
		return nil, status.Error(codes.Internal, "internal server error")
	}

	return &userpb.GetUserResponse{
		User: &userpb.User{
			Id:        u.ID.String(),
			Phone:     u.Phone,
			CreatedAt: timestamppb.New(u.CreatedAt),
		},
	}, nil
}
