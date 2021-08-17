package middleware

import (
	"net/http"
	"strings"
	"time"

	"github.com/vsatcloud/mars/proto"

	"github.com/vsatcloud/mars/utils"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

type Auth struct {
	TokenSecret string
	SkipList    []string // ["GET /api/v1/user/login", ...]
	HeaderKey   string
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

		key := "Authorization"
		if j.HeaderKey != "" {
			key = j.HeaderKey
		}

		tokenBearer := c.Request.Header.Get(key)
		if len(tokenBearer) == 0 || !strings.HasPrefix(tokenBearer, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusOK, gin.H{
				"code":    proto.CodeFailedAuthVerify,
				"message": "身份验证失败",
			})
			return
		}

		token := strings.Split(tokenBearer, " ")[1]

		jwtT, err := utils.ParseToken(token, j.TokenSecret)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusOK, gin.H{
				"code":    proto.CodeFailedAuthVerify,
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
				"code":    proto.CodeTokenExpired,
				"message": proto.CodeMsg[proto.CodeTokenExpired],
			})
			return
		}

		//权限设置
		roles := claims["roles"].(string)
		c.Set("roles", roles)

		c.Next()
	}
}
