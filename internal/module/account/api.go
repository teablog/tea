package account

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/teablog/tea/internal/helper"
	"strings"
)

func NameExists(ctx *gin.Context) {
	name := ctx.Query("name")
	name = strings.Trim(name, " \r\n")
	if name == "" {
		helper.Fail(ctx, errors.New("起个像样点的名字吧～"))
	}
}
