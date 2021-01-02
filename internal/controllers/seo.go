package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/teablog/tea/internal/helper"
	"github.com/teablog/tea/internal/module/seo"
)

var Seo _seo

type _seo struct{}

func (s *_seo) SiteMap(ctx *gin.Context) {
	if err := seo.Sitemap.Generate(ctx); err != nil {
		helper.Fail(ctx, err)
		return
	}
	helper.Success(ctx, "success")
	return
}

