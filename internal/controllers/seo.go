package controllers

import (
	"dyc/internal/helper"
	"dyc/internal/module/seo"
	"github.com/gin-gonic/gin"
)

var Seo _seo

type _seo struct{}

func (s *_seo) SiteMap(ctx *gin.Context) {
	if err := seo.Sitemap.Generate(ctx); err != nil {
		helper.Fail(ctx, err)
		return
	}
	helper.Success(ctx, "success")
}
