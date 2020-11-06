package mars

import (
	"github.com/gin-gonic/gin"
)

const (
	VERSION = "1.0.0"
)

func New() *gin.Engine {
	r := gin.New()

	//r.Use(ginzerolog.Logger("gin"))
	//r.Use(gin.Recovery())

	//r.GET("ping", func(c *gin.Context) { c.String(http.StatusOK, "pong") })

	return r
}
