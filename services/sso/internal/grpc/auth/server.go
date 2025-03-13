package auth

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/kordyd/remember_me-golang/protos/gens/go/sso"
	"github.com/kordyd/remember_me-golang/services/sso/internal/services/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Auth interface {
	Login(ctx context.Context, email string, password string, appId string) (string, error)
	Register(ctx context.Context, email string, password string) (string, error)
	// TODO Logout()
}

type serverAPI struct {
	sso.UnimplementedAuthServer
	auth Auth
}

func Register(gRPCServer *grpc.Server, auth Auth) {
	sso.RegisterAuthServer(gRPCServer, &serverAPI{auth: auth})
}

func (s serverAPI) Login(ctx context.Context, req *sso.LoginRequest) (*sso.LoginResponse, error) {
	if req.GetEmail() == "" {
		return nil, status.Error(codes.InvalidArgument, "missing email")
	}

	if req.GetPassword() == "" {
		return nil, status.Error(codes.InvalidArgument, "missing password")
	}

	if req.GetAppId() == "" {
		return nil, status.Error(codes.InvalidArgument, "missing app_id")
	}
	if err := uuid.Validate(req.GetAppId()); err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid app_id, need to be uuid")
	}
	token, err := s.auth.Login(ctx, req.GetEmail(), req.GetPassword(), req.GetAppId())
	if err != nil {
		if errors.Is(err, auth.ErrInvalidCredentials) {
			return nil, status.Error(codes.InvalidArgument, "invalid credentials")
		}
		return nil, status.Error(codes.Internal, "something went wrong")
	}
	return &sso.LoginResponse{Token: token}, nil
}

func (s serverAPI) Register(ctx context.Context, req *sso.RegisterRequest) (*sso.RegisterResponse, error) {
	if req.GetEmail() == "" {
		return nil, status.Error(codes.InvalidArgument, "missing email")
	}
	if req.GetPassword() == "" {
		return nil, status.Error(codes.InvalidArgument, "missing password")
	}
	userId, err := s.auth.Register(ctx, req.GetEmail(), req.GetPassword())
	if err != nil {
		if errors.Is(err, auth.ErrUserAlreadyExists) {
			return nil, status.Error(codes.InvalidArgument, "user already exists")
		}
		return nil, status.Error(codes.Internal, "something went wrong")
	}
	return &sso.RegisterResponse{UserId: userId}, nil
}
