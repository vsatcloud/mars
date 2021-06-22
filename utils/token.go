package utils

import (
	"time"

	"github.com/dgrijalva/jwt-go"
)

//生成token
func GenerateToken(data map[string]interface{}, secret string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	for key, value := range data {
		token.Claims.(jwt.MapClaims)[key] = value
	}
	//设置token过期时间
	token.Claims.(jwt.MapClaims)["expired_at"] = time.Now().Add(time.Hour * 24 * 7).Unix()

	return token.SignedString([]byte(secret))
}

//解析token
func ParseToken(tokenString string, secret string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err == nil && token.Valid {
		return token, nil
	} else {
		return nil, err
	}
}
