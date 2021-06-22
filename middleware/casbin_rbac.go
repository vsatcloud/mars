package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/vsatcloud/mars"
	"github.com/vsatcloud/mars/casbin"
	"github.com/vsatcloud/mars/models"
)

func CasbinRBAC(db models.Database) gin.HandlerFunc {
	return func(c *gin.Context) {
		authorityID, _ := c.Get("authority_id")
		// 获取请求的URI
		obj := c.Request.URL.RequestURI()
		// 获取请求方法
		act := c.Request.Method
		// 获取用户的角色
		sub := authorityID.(string)
		e := casbin.Casbin(db)
		// 判断策略中是否存在
		success := e.Enforce(sub, obj, act)
		if success {
			c.Next()
		} else {
			c.AbortWithStatusJSON(http.StatusOK, gin.H{
				"code":    mars.CodeNoPerm,
				"message": mars.CodeMsg[mars.CodeNoPerm],
			})
			return
		}
	}
}
