package xzbase

import (
	"github.com/gin-gonic/gin"
	"strconv"
)

type Controller interface {
	SetContext(ctx *gin.Context)
}
type controller struct {
	ctx *gin.Context
}

func (c*controller) SetContext(ctx *gin.Context) {
	c.ctx=ctx
}

func (c*controller)QueryInt64(key string,def int64) int64 {
	v:=c.ctx.Query(key)
	i,e:=strconv.ParseInt(v,0,64)
	if e!=nil{
		return def
	}
	return i
}

