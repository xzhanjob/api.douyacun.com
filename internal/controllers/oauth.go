package controllers

import (
	"dyc/internal/helper"
	"dyc/internal/module/account"
	"encoding/json"
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
	redirectUri := ctx.DefaultQuery("redirect_uri", "https://www.douyacun.com/")
	github := account.NewGithub()
	if err := github.Token(code); err != nil {
		helper.Fail(ctx, err)
		return
	}
	if err := github.User(); err != nil {
		helper.Fail(ctx, err)
		return
	}
	user, err := account.Account.Create(github)
	if err != nil {
		helper.Fail(ctx, err)
	}
	data, err := json.Marshal(user)
	if err == nil {
		ctx.SetCookie("douyacun", string(data), 604800, "/", "douyacun.com", false, true)
	}
	ctx.Redirect(302, redirectUri)
}

func (*_oauth) Google(ctx *gin.Context) {
	google, err := account.NewGoogle(ctx)
	if err != nil {
		helper.Fail(ctx, err)
		return
	}
	user, err := account.Account.Create(google)
	if err != nil {
		helper.Fail(ctx, err)
		return
	}
	data, err := json.Marshal(user)
	if err == nil {
		ctx.SetCookie("douyacun", string(data), 604800, "/", "douyacun.com", false, true)
	}
	helper.Success(ctx, data)
	return
}
