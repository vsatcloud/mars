package mars

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func New() *gin.Engine {
	r := gin.New()

	r.GET("ping", func(c *gin.Context) { c.String(http.StatusOK, "pong") })

	return r
}
