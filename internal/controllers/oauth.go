package controllers

import (
	"dyc/internal/helper"
	"dyc/internal/module/account"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

var Oauth *_oauth

type _oauth struct{}

func (*_oauth) Github(ctx *gin.Context) {
	code := ctx.Query("code")
	if code == "" {
		helper.Fail(ctx, errors.Errorf("code参数丢失！"))
		return
	}
	github := account.NewGithub()
	if err := github.Token(code); err != nil {
		helper.Fail(ctx, err)
		return
	}
	if err := github.User(); err != nil {
		helper.Fail(ctx, err)
		return
	}
	_, err := account.Account.Create(github)
	if err != nil {
		helper.Fail(ctx, err)
	}
	ctx.Redirect(302, "https://www.douyacun.com/")
}

func (*_oauth) Google(ctx *gin.Context) {
	google, err := account.NewGoogle(ctx)
	if err != nil {
		helper.Fail(ctx, err)
		return
	}
	data, err := account.Account.Create(google)
	if err != nil {
		helper.Fail(ctx, err)
		return
	}
	helper.Success(ctx, data)
	return
}
