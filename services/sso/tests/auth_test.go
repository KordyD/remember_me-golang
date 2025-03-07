package tests

import (
	"github.com/brianvoe/gofakeit"
	"github.com/golang-jwt/jwt/v5"
	"github.com/kordyd/remember_me-golang/protos/gens/go/sso"
	"github.com/kordyd/remember_me-golang/sso/tests/suite"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

const (
	secret = "test"
	appId  = "ec62edbe-2b8f-4d8e-a46f-885f1808403f"
)

func TestLogin_HappyPath(t *testing.T) {
	st, ctx := suite.New(t)
	email := gofakeit.Email()
	password := gofakeit.Password(true, true, true, true, false, 15)
	regResp, err := st.AuthClient.Register(ctx, &sso.RegisterRequest{
		Email:    email,
		Password: password,
	})
	require.NoError(t, err)
	userId := regResp.GetUserId()
	assert.NotEmpty(t, userId)

	logResp, err := st.AuthClient.Login(ctx, &sso.LoginRequest{
		Email:    email,
		Password: password,
		AppId:    appId,
	})
	require.NoError(t, err)

	loginTime := time.Now()

	token := logResp.GetToken()
	assert.NotEmpty(t, token)

	tokenParsed, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	require.NoError(t, err)

	claims, ok := tokenParsed.Claims.(jwt.MapClaims)
	require.True(t, ok)

	assert.Equal(t, email, claims["email"].(string))
	assert.Equal(t, userId, claims["userId"].(string))
	assert.Equal(t, appId, claims["appId"].(string))

	assert.InDelta(t, loginTime.Add(st.Cfg.TokenTTL).Unix(), claims["exp"].(float64), 1)

}

func TestRegisterLogin_DuplicatedRegistration(t *testing.T) {
	st, ctx := suite.New(t)

	email := gofakeit.Email()
	password := gofakeit.Password(true, true, true, true, false, 15)

	respReg, err := st.AuthClient.Register(ctx, &sso.RegisterRequest{
		Email:    email,
		Password: password,
	})
	require.NoError(t, err)
	require.NotEmpty(t, respReg.GetUserId())

	respReg, err = st.AuthClient.Register(ctx, &sso.RegisterRequest{
		Email:    email,
		Password: password,
	})
	require.Error(t, err)
	assert.Empty(t, respReg.GetUserId())
	assert.ErrorContains(t, err, "user already exists")
}

func TestLogin_FailCases(t *testing.T) {
	st, ctx := suite.New(t)

	tests := []struct {
		name        string
		email       string
		password    string
		appId       string
		expectedErr string
	}{
		{
			name:        "Login with Empty Password",
			email:       gofakeit.Email(),
			password:    "",
			appId:       appId,
			expectedErr: "missing password",
		},
		{
			name:        "Login with Empty Email",
			email:       "",
			password:    gofakeit.Password(true, true, true, true, false, 15),
			appId:       appId,
			expectedErr: "missing email",
		},
		{
			name:        "Login with Both Empty Email and Password",
			email:       "",
			password:    "",
			appId:       appId,
			expectedErr: "missing email",
		},
		{
			name:        "Login with Non-Matching Password",
			email:       gofakeit.Email(),
			password:    gofakeit.Password(true, true, true, true, false, 15),
			appId:       appId,
			expectedErr: "invalid credentials",
		},
		{
			name:        "Login without AppID",
			email:       gofakeit.Email(),
			password:    gofakeit.Password(true, true, true, true, false, 15),
			appId:       "",
			expectedErr: "missing app_id",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := st.AuthClient.Register(ctx, &sso.RegisterRequest{
				Email:    gofakeit.Email(),
				Password: gofakeit.Password(true, true, true, true, false, 15),
			})
			require.NoError(t, err)

			_, err = st.AuthClient.Login(ctx, &sso.LoginRequest{
				Email:    tt.email,
				Password: tt.password,
				AppId:    tt.appId,
			})
			require.Error(t, err)
			require.Contains(t, err.Error(), tt.expectedErr)
		})
	}
}
