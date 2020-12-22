package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/teablog/tea/internal/helper"
	"github.com/teablog/tea/internal/module/account"
	"net/http"
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
		ctx.String(http.StatusForbidden, err.Error())
		return
	}
	if err := github.User(); err != nil {
		helper.Fail(ctx, err)
		return
	}
	user, err := account.NewAccount().Create(ctx, github)
	if err != nil {
		helper.Fail(ctx, err)
	}
	user.SetCookie(ctx)
	ctx.Redirect(302, redirectUri)
}

func (*_oauth) Google(ctx *gin.Context) {
	google, err := account.NewGoogle(ctx)
	if err != nil {
		helper.Fail(ctx, err)
		return
	}
	user, err := account.NewAccount().Create(ctx, google)
	if err != nil {
		helper.Fail(ctx, err)
		return
	}
	user.SetCookie(ctx)
	helper.Success(ctx, user)
	return
}
