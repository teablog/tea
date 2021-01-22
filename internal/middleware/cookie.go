package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
	"github.com/teablog/tea/internal/config"
	"github.com/teablog/tea/internal/consts"
	"github.com/teablog/tea/internal/logger"
)

func CookieUUID() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ck, err := ctx.Cookie(consts.CookieUUID)
		if err != nil {
			logger.Wrapf(err, "ctx cooke read err")
		}
		if ck == "" {
			uid, err := uuid.NewV4()
			if err != nil {
				logger.Wrapf(err, "failed to generate UUID:")
			} else {
				// 设置cookie
				ctx.SetCookie(consts.CookieUUID, uid.String(), config.Global.CookieMaxAge(), "/", config.Global.Domain(), false, false)
			}
			// 避免首次访问，取不到cookie
			ctx.Set(consts.CookieUUID, uid.String())
		} else {
			ctx.Set(consts.CookieUUID, ck)
		}
		ctx.Next()
	}
}
