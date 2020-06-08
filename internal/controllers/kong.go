package controllers

import (
	"dyc/internal/helper"
	"github.com/gin-gonic/gin"
)

var Kong *_kong

type _kong struct{}

func (*_kong) PreserveHost(ctx *gin.Context) {
	helper.Success(ctx, gin.H{
		"header": ctx.Request.Header,
		"host":   ctx.Request.Host,
	})
}
