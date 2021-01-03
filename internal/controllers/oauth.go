package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/teablog/tea/internal/config"
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
	redirectUri := ctx.DefaultQuery("redirect_uri", config.Global.Host())
	github := account.NewGithub()
	if err := github.TokenV2(code); err != nil {
		ctx.String(http.StatusForbidden, err.Error())
		return
	}
	if err := github.UserV2(); err != nil {
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
