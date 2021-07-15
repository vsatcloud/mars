package mars

import (
	"errors"
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
	if contentType == "application/json" {
		if c.Request.ContentLength > 0 {
			c.Err = c.ShouldBindJSON(obj)
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

func (c *Context) SetData(data interface{}) {
	c.Result.Data = data
}

func (c *Context) SetCode(code int) {
	c.Result.Code = code
}
