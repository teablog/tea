package helper

import (
	"dyc/internal/consts"
	"github.com/gin-gonic/gin"
	"strconv"
)

var Pagenation _PageNation

type _PageNation struct {
	Page  uint `json:"page"`
	Limit uint `json:"limit"`
	Total uint `json:"total"`
}

func (_PageNation) Init(ctx *gin.Context) _PageNation {
	var (
		page _PageNation
	)
	pUint64, err := strconv.ParseUint(ctx.Param("page"), 10, 64)
	if err != nil {
		page.Page = consts.DefaultPage
	} else {
		page.Page = uint(pUint64)
		if page.Page == 0 {
			page.Page = 1
		}
	}
	lUint64, err := strconv.ParseUint(ctx.Param("limit"), 10, 64)
	if err != nil {
		page.Limit = consts.DefaultPage
	} else {
		page.Limit = uint(lUint64)
		if page.Limit == 0 {
			page.Limit = consts.DefaultPage
		}
	}
	return page
}
