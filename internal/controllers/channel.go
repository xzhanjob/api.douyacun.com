package controllers

import (
	"dyc/internal/helper"
	"dyc/internal/module/chat"
	"dyc/internal/validate"
	"github.com/gin-gonic/gin"
)

var Channel *_channel

type _channel struct{}

func (*_channel) Private(ctx *gin.Context) {
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
	c, err := chat.Channel.Create(ctx, v)
	if err != nil {
		helper.Fail(ctx, err)
		return
	}
	helper.Success(ctx, c)
}
