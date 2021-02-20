package middleware

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/vsatcloud/mars/models"

	"github.com/rs/zerolog/log"

	"github.com/gin-gonic/gin"
)

func OperationRecord() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 不处理get相关接口,所以请不要再get接口中涉及到数据修改的操作
		if c.Request.Method == http.MethodGet {
			c.Next()
			return
		}

		var body []byte
		var userId uint
		var err error
		body, err = ioutil.ReadAll(c.Request.Body)
		if err != nil {
			log.Warn().Msg(err.Error())
		} else {
			c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(body))
		}

		record := models.OperationRecord{
			Ip:     c.ClientIP(),
			Method: c.Request.Method,
			Path:   c.Request.URL.Path,
			Agent:  c.Request.UserAgent(),
			Body:   string(body),
			UserID: userId,
		}
		// 存在某些未知错误 TODO
		//values := c.Request.Header.Values("content-type")
		//if len(values) >0 && strings.Contains(values[0], "boundary") {
		//	record.Body = "file"
		//}
		writer := responseBodyWriter{
			ResponseWriter: c.Writer,
			body:           &bytes.Buffer{},
		}
		c.Writer = writer
		now := time.Now()

		c.Next()

		latency := time.Now().Sub(now)
		record.ErrorMessage = c.Errors.ByType(gin.ErrorTypePrivate).String()
		record.Status = c.Writer.Status()
		record.Latency = latency
		record.Resp = writer.body.String()

		if err := record.Save(); err != nil {
			log.Warn().Msg("create operation record error:" + err.Error())
		}
	}
}

type responseBodyWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (r responseBodyWriter) Write(b []byte) (int, error) {
	r.body.Write(b)
	return r.ResponseWriter.Write(b)
}
