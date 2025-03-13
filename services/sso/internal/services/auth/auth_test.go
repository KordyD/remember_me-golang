package auth

import (
	"context"
	"github.com/brianvoe/gofakeit"
	"github.com/golang-jwt/jwt/v5"
	"github.com/kordyd/remember_me-golang/services/sso/internal/models"
	"github.com/kordyd/remember_me-golang/services/sso/internal/services/auth/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
	"log/slog"
	"testing"
	"time"
)

func TestAuthServiceLogin_HappyPath(t *testing.T) {
	const (
		appId  = "test"
		secret = "test"
	)
	ctx := context.Background()
	userId := gofakeit.UUID()
	email := gofakeit.Email()
	password := gofakeit.Password(true, true, true, true, false, 15)

	passHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		t.Fatalf("failed to encrypt pass: %v", err)
	}

	mockedUserStorage := mocks.NewMockUserStorage(t)
	mockedAppStorage := mocks.NewMockAppStorage(t)

	mockedUserStorage.On("SaveUser", ctx, email, mock.AnythingOfType("[]uint8")).Return(userId, nil)
	mockedAppStorage.On("GetApp", ctx, appId).Return(models.App{
		Id:     appId,
		Name:   "test",
		Secret: secret,
	}, nil)
	mockedUserStorage.On("GetUser", ctx, email).Return(models.User{
		Id:       userId,
		Email:    email,
		PassHash: passHash,
	}, nil)

	tokenTTL := 24 * time.Hour

	authService := New(slog.Default(), mockedUserStorage, mockedAppStorage, tokenTTL)
	resUserId, err := authService.Register(ctx, email, password)
	require.NoError(t, err)
	assert.NotEmpty(t, resUserId)

	token, err := authService.Login(ctx, email, password, appId)
	require.NoError(t, err)
	assert.NotEmpty(t, token)

	loginTime := time.Now()

	tokenParsed, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	require.NoError(t, err)

	claims, ok := tokenParsed.Claims.(jwt.MapClaims)
	require.True(t, ok)

	assert.Equal(t, email, claims["email"].(string))
	assert.Equal(t, userId, claims["userId"].(string))
	assert.Equal(t, appId, claims["appId"].(string))

	assert.InDelta(t, loginTime.Add(tokenTTL).Unix(), claims["exp"].(float64), 1)
}
