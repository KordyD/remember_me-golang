package auth

import (
	"context"
	"errors"
	"github.com/kordyd/remember_me-golang/services/sso/internal/services/auth"
	"testing"

	"github.com/google/uuid"
	"github.com/kordyd/remember_me-golang/protos/gens/go/sso"
	"github.com/kordyd/remember_me-golang/services/sso/internal/grpc/auth/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestLogin_Success(t *testing.T) {
	mockAuth := mocks.NewMockAuth(t)
	server := serverAPI{auth: mockAuth}

	email := "test@example.com"
	password := "password"
	appID := uuid.NewString()
	expectedToken := "token123"

	mockAuth.On("Login", mock.Anything, email, password, appID).
		Return(expectedToken, nil).Once()

	req := &sso.LoginRequest{
		Email:    email,
		Password: password,
		AppId:    appID,
	}

	resp, err := server.Login(context.Background(), req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, expectedToken, resp.Token)

}

func TestLogin_InvalidArguments(t *testing.T) {
	mockAuth := mocks.NewMockAuth(t)
	server := serverAPI{auth: mockAuth}

	tests := []struct {
		name    string
		request *sso.LoginRequest
		errCode codes.Code
	}{
		{"Missing email", &sso.LoginRequest{Password: "pass", AppId: uuid.NewString()}, codes.InvalidArgument},
		{"Missing password", &sso.LoginRequest{Email: "email@example.com", AppId: uuid.NewString()}, codes.InvalidArgument},
		{"Missing appId", &sso.LoginRequest{Email: "email@example.com", Password: "pass"}, codes.InvalidArgument},
		{"Invalid appId format", &sso.LoginRequest{Email: "email@example.com", Password: "pass", AppId: "invalid-uuid"}, codes.InvalidArgument},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			resp, err := server.Login(context.Background(), tc.request)
			assert.Nil(t, resp)
			assert.Error(t, err)
			assert.Equal(t, tc.errCode, status.Code(err))
		})
	}
}

func TestLogin_InvalidCredentials(t *testing.T) {
	mockAuth := mocks.NewMockAuth(t)
	server := serverAPI{auth: mockAuth}

	email := "test@example.com"
	password := "wrongpassword"
	appID := uuid.NewString()

	mockAuth.On("Login", mock.Anything, email, password, appID).
		Return("", auth.ErrInvalidCredentials).Once()

	req := &sso.LoginRequest{
		Email:    email,
		Password: password,
		AppId:    appID,
	}

	resp, err := server.Login(context.Background(), req)

	assert.Nil(t, resp)
	assert.Error(t, err)
	assert.Equal(t, codes.InvalidArgument, status.Code(err))
}

func TestLogin_InternalError(t *testing.T) {
	mockAuth := mocks.NewMockAuth(t)
	server := serverAPI{auth: mockAuth}

	email := "test@example.com"
	password := "password"
	appID := uuid.NewString()

	mockAuth.On("Login", mock.Anything, email, password, appID).
		Return("", errors.New("unexpected error")).Once()

	req := &sso.LoginRequest{
		Email:    email,
		Password: password,
		AppId:    appID,
	}

	resp, err := server.Login(context.Background(), req)

	assert.Nil(t, resp)
	assert.Error(t, err)
	assert.Equal(t, codes.Internal, status.Code(err))
}

func TestRegister_Success(t *testing.T) {
	mockAuth := mocks.NewMockAuth(t)
	server := serverAPI{auth: mockAuth}

	email := "test@example.com"
	password := "password"
	userID := uuid.NewString()

	mockAuth.On("Register", mock.Anything, email, password).
		Return(userID, nil).Once()

	req := &sso.RegisterRequest{
		Email:    email,
		Password: password,
	}

	resp, err := server.Register(context.Background(), req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, userID, resp.UserId)
}

func TestRegister_InvalidArguments(t *testing.T) {
	mockAuth := mocks.NewMockAuth(t)
	server := serverAPI{auth: mockAuth}

	tests := []struct {
		name    string
		request *sso.RegisterRequest
		errCode codes.Code
	}{
		{"Missing email", &sso.RegisterRequest{Password: "pass"}, codes.InvalidArgument},
		{"Missing password", &sso.RegisterRequest{Email: "email@example.com"}, codes.InvalidArgument},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			resp, err := server.Register(context.Background(), tc.request)
			assert.Nil(t, resp)
			assert.Error(t, err)
			assert.Equal(t, tc.errCode, status.Code(err))
		})
	}
}

func TestRegister_UserAlreadyExists(t *testing.T) {
	mockAuth := mocks.NewMockAuth(t)
	server := serverAPI{auth: mockAuth}

	email := "existing@example.com"
	password := "password"

	mockAuth.On("Register", mock.Anything, email, password).
		Return("", auth.ErrUserAlreadyExists).Once()

	req := &sso.RegisterRequest{
		Email:    email,
		Password: password,
	}

	resp, err := server.Register(context.Background(), req)

	assert.Nil(t, resp)
	assert.Error(t, err)
	assert.Equal(t, codes.InvalidArgument, status.Code(err))
}

func TestRegister_InternalError(t *testing.T) {
	mockAuth := mocks.NewMockAuth(t)
	server := serverAPI{auth: mockAuth}

	email := "test@example.com"
	password := "password"

	mockAuth.On("Register", mock.Anything, email, password).
		Return("", errors.New("unexpected error")).Once()

	req := &sso.RegisterRequest{
		Email:    email,
		Password: password,
	}

	resp, err := server.Register(context.Background(), req)

	assert.Nil(t, resp)
	assert.Error(t, err)
	assert.Equal(t, codes.Internal, status.Code(err))
}
