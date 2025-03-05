package jwt

import (
	"github.com/golang-jwt/jwt/v5"
	"time"
)

func NewUserJwt(userId string, email string, duration time.Duration, appId string, appSecret string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email":  email,
		"exp":    time.Now().Add(duration).Unix(),
		"appId":  appId,
		"userId": userId,
	})
	tokenString, err := token.SignedString([]byte(appSecret))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}
