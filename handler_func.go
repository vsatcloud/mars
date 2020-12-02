package mars

import "github.com/gin-gonic/gin"

type HandlerFunc func(*Context)

func HandlerWarp(handler HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		context := new(Context)
		context.Context = c
		handler(context)
	}
}
