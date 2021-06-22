package middleware

import (
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/vsatcloud/mars"
)

type Auth struct {
	TokenSecret string
	SkipList    []string // ["GET /api/v1/user/login", ...]
}

func (j *Auth) Authenticator() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 跳过
		for _, item := range j.SkipList {
			routes := strings.Split(item, " ")
			if strings.Contains(c.Request.URL.Path, routes[1]) && c.Request.Method == routes[0] {
				c.Next()
				return
			}

		}

		tokenBearer := c.Request.Header.Get("Authorization")
		if len(tokenBearer) == 0 || !strings.HasPrefix(tokenBearer, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusOK, gin.H{
				"code":    mars.CodeFailedAuthVerify,
				"message": "身份验证失败",
			})
			return
		}

		token := strings.Split(tokenBearer, " ")[1]

		jwtT, err := ParseToken(token, j.TokenSecret)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusOK, gin.H{
				"code":    mars.CodeFailedAuthVerify,
				"message": "身份验证失败",
			})
			return
		}

		claims := jwtT.Claims.(jwt.MapClaims)
		userID := uint(claims["user_id"].(float64))
		c.Set("user_id", userID)

		//token过期
		expiredAt := int64(claims["expired_at"].(float64))
		if time.Unix(expiredAt, 0).Before(time.Now()) {
			c.AbortWithStatusJSON(http.StatusOK, gin.H{
				"code":    mars.CodeTokenExpired,
				"message": mars.CodeMsg[mars.CodeTokenExpired],
			})
			return
		}

		//权限设置
		authorityID := claims["authority_id"].(string)
		c.Set("authority_id", authorityID)

		c.Next()
	}
}

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
