package xzbase

import (
	"github.com/gin-gonic/gin"
	"strconv"
)

type controller interface {
	setContext(ctx *gin.Context)
}
type BaseController struct {
	Ctx *gin.Context
}

func (c*BaseController) setContext(ctx *gin.Context) {
	c.Ctx=ctx
}

func (c*BaseController)QueryInt64(key string,def int64) int64 {
	v:=c.Ctx.Query(key)
	i,e:=strconv.ParseInt(v,0,64)
	if e!=nil{
		return def
	}
	return i
}

