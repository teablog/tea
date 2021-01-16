package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/teablog/tea/internal/config"
	"github.com/teablog/tea/internal/logger"
	"net/http"
	"strings"
)

func TokenAuthCheck() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		t := ctx.Request.Header.Get("Authorization")
		if t == "" {
			t = ctx.Query("token")
		} else {
			t = strings.TrimLeft(t, "Bearer ")
		}
		logger.Debugf("token: %s", t)
		if t == "" || t != config.Global.Token() {
			ctx.JSON(http.StatusUnauthorized, "未授权访问～")
			ctx.Abort()
			return
		}
		ctx.Next()
	}
}