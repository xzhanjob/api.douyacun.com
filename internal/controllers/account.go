package controllers

import (
	"dyc/internal/helper"
	"dyc/internal/module/account"
	"github.com/gin-gonic/gin"
)

var Account *_account

type _account struct{}

func (*_account) List(ctx *gin.Context) {
	q := ctx.DefaultQuery("q", "")
	data, err := account.NewAccount().All(q)
	if err != nil {
		helper.Fail(ctx, err)
		return
	}
	helper.Success(ctx, data)
	return
}
