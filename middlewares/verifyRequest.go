package middlewares

import (
	resutil "github.com/gangjun06/mySchoolApp-api/utils/res"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

func VerifyRequest(data interface{}) gin.HandlerFunc {
	return func(c *gin.Context) {
		r := resutil.New(c)
		if err := c.ShouldBindWith(data, binding.JSON); err != nil {
			r.SendError(resutil.ERR_BAD_REQUEST, err.Error())
			return
		}
		c.Set("body", data)
	}
}

func VerifyQuery(data interface{}) gin.HandlerFunc {
	return func(c *gin.Context) {
		r := resutil.New(c)
		if err := c.ShouldBindQuery(data); err != nil {
			r.SendError(resutil.ERR_BAD_REQUEST, err.Error())
			return
		}
		c.Set("query", data)
	}
}
