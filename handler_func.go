package mars

import (
	"github.com/gin-gonic/gin"
)

type HandlerFunc func(*Context)

func HandlerWarp(handler HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		context := new(Context)
		context.Context = c
		handler(context)
	}
}

type HandlerFuncJson func(*Context) error

func HandlerWarpJson(handler HandlerFuncJson) gin.HandlerFunc {
	return func(c *gin.Context) {
		context := new(Context)
		context.Context = c
		err := handler(context)
		if err != nil {
			context.SystemError(err)
		}

		context.ResponseJson()
	}
}
