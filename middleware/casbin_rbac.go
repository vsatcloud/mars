package middleware

import (
	"net/http"
	"strings"

	"github.com/vsatcloud/mars/proto"

	"gorm.io/gorm"

	"github.com/gin-gonic/gin"
	"github.com/vsatcloud/mars/casbin"
)

type Casbin struct {
	Db       *gorm.DB
	SkipList []string // ["GET /api/v1/user/login", ...]

}

func (cas *Casbin) CasbinRBAC() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 跳过
		for _, item := range cas.SkipList {
			routes := strings.Split(item, " ")
			if strings.Contains(c.Request.URL.Path, routes[1]) && c.Request.Method == routes[0] {
				c.Next()
				return
			}

		}

		roles, _ := c.Get("roles")
		// 获取请求的URI
		obj := c.Request.URL.RequestURI()
		// 获取请求方法
		act := c.Request.Method
		// 获取用户的角色
		sub := roles.(string)
		e := casbin.Casbin(cas.Db)
		// 判断策略中是否存在
		success, _ := e.Enforce(sub, obj, act)
		if success {
			c.Next()
		} else {
			c.AbortWithStatusJSON(http.StatusOK, gin.H{
				"code":    proto.CodeNoPerm,
				"message": proto.CodeMsg[proto.CodeNoPerm],
			})
			return
		}
	}
}
