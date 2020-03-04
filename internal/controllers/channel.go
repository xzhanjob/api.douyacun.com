package controllers

import (
	"dyc/internal/consts"
	"dyc/internal/helper"
	"dyc/internal/module/chat"
	"dyc/internal/validate"
	"github.com/gin-gonic/gin"
)

var Channel *_channel

type _channel struct{}

func (*_channel) Create(ctx *gin.Context) {
	var (
		v *validate.ChannelCreateValidator
		err error
	)
	if err = ctx.ShouldBindJSON(&v); err != nil {
		helper.Fail(ctx, err)
		return
	}
	if err = validate.DoValidate(v); err != nil {
		helper.Fail(ctx, err)
		return
	}
	if v.Type == consts.TypeChannelPrivate {
		if c, ok := chat.Channel.Get(ctx, v); ok {
			helper.Success(ctx, c)
			return
		}
	}
	c, err := chat.Channel.Create(ctx, v)
	if err != nil {
		helper.Fail(ctx, err)
		return
	}
	helper.Success(ctx, c)
}

func (*_channel) List(ctx *gin.Context)  {
	data, err := chat.Channel.List(ctx)
	if err != nil {
		helper.Fail(ctx, err)
		return
	}
	helper.Success(ctx, data)
	return
}