package tracke

import (
	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
	"github.com/teablog/tea/internal/config"
	"github.com/teablog/tea/internal/consts"
	"github.com/teablog/tea/internal/logger"
)

func Route(api *gin.RouterGroup) {
	t := api.Group("/tracking")
	{
		t.GET("/impress", Impression)
	}
}

func Impression(ctx *gin.Context) {
	ck, _ := ctx.Cookie(consts.CookieUUID)
	if ck == "" {
		uid, err := uuid.NewV4()
		if err != nil {
			logger.Wrapf(err, "failed to generate UUID:")
		} else {
			// 设置cookie
			ctx.SetCookie(consts.CookieUUID, uid.String(), config.Global.CookieMaxAge(), "/", "." + config.Global.Domain(), false, false)
		}
	}
}
