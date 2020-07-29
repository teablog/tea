package controllers

import (
	"github.com/teablog/tea/internal/helper"
	"github.com/teablog/tea/internal/module/media"
	"github.com/gin-gonic/gin"
	"strconv"
)

var Media _Media

type _Media struct {
	Movie _Movie
	TV    _TV
}

type _Movie struct{}
type _TV struct{}

func (*_Media) Index(ctx *gin.Context) {
	subtype := ctx.Param("subtype")
	page, err := strconv.Atoi(ctx.Query("page"))
	if err != nil {
		page = 1
	}
	total, data, err := media.Resource.Index(page, subtype)
	if err != nil {
		helper.Fail(ctx, err)
		return
	}
	helper.Success(ctx, gin.H{"total": total, "data": data})
	return
}

func (*_Media) View(ctx *gin.Context) {
	id := ctx.Param("id")
	data, err := media.Resource.View(id)
	if err != nil {
		helper.Fail(ctx, err)
		return
	}
	res, err := media.Resource.ToArticle(data)
	if err != nil {
		helper.Fail(ctx, err)
		return
	}
	helper.Success(ctx, gin.H{"data": res})
	return
}

func (*_Media) Search(ctx *gin.Context) {
	q := ctx.Query("q")
	page, err := strconv.Atoi(ctx.Query("page"))
	if err != nil {
		page = 1
	}
	total, data, err := media.Resource.Search(page, q)
	if err != nil {
		helper.Fail(ctx, err)
	}
	helper.Success(ctx, gin.H{"total": total, "data": data})
}

