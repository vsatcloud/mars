package mars

import (
	"errors"
	"io"
	"net/http"
	"strings"

	"github.com/vsatcloud/mars/proto"

	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
	"github.com/mcuadros/go-defaults"
	"github.com/rs/zerolog/log"
)

type Context struct {
	*gin.Context
	Result proto.Result
	Err    error
}

func (c *Context) BindParams(obj interface{}) error {
	defaults.SetDefaults(obj)

	contentType := c.ContentType()
	if contentType == "application/json" && c.Request.ContentLength > 0 {
		c.Err = c.ShouldBindJSON(obj)
		if errors.Is(c.Err, io.EOF) {
			return nil
		}
	} else {
		c.Err = c.ShouldBind(obj)
	}
	if c.Err != nil {
		c.Result.Code = proto.CodeErrParams
		return c.Err
	}

	valid := validation.Validation{}
	var b bool
	b, c.Err = valid.Valid(obj)
	if c.Err != nil {
		c.Result.Code = proto.CodeErrParams
		return nil
	}
	if !b {
		if valid.HasErrors() {
			var errstrs []string
			for _, err := range valid.Errors {
				errstrs = append(errstrs, err.Field+" "+err.String())
			}

			c.Result.Code = proto.CodeErrParams
			return errors.New(strings.Join(errstrs, ";"))
		}
	}

	log.Debug().Interface("param", obj)

	return nil
}

func (c *Context) ResponseJson() {
	if c.Result.Message == "" {
		c.Result.Message = proto.CodeMsg[c.Result.Code]
	}
	c.JSON(http.StatusOK, c.Result)
}

func (c *Context) SystemError(err error) {
	c.Err = err
	c.Result.Code = proto.CodeErrSystem
	c.Result.Detail = c.Err.Error()
}

func (c *Context) SetData(data interface{}) error {
	c.Result.Data = data
	return nil
}

func (c *Context) RespOK(data interface{}) error {
	c.Result.Data = data
	return nil
}

func (c *Context) RespCode(code int) error {
	c.Result.Code = code
	return nil
}

func (c *Context) RespMsg(msg string) error {
	c.Result.Code = proto.CodeErrLogic
	c.Result.Message = msg
	return nil
}

func (c *Context) SetCode(code int) error {
	c.Result.Code = code
	return nil
}

func (c *Context) LogicError(msg string) error {
	c.Result.Code = proto.CodeErrLogic
	c.Result.Message = msg

	return nil
}
