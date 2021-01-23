package tracke

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
	"github.com/teablog/tea/internal/config"
	"github.com/teablog/tea/internal/consts"
	"github.com/teablog/tea/internal/logger"
	"github.com/teablog/tea/internal/middleware"
	"time"
)

func Route(api *gin.RouterGroup) {
	t := api.Group("/tracking")
	{
		t.GET("/impress", Impression)
	}
}

func Impression(ctx *gin.Context) {
	ck, _ := ctx.Cookie(consts.CookieUUIDV2)
	if ck == "" {
		uid, err := uuid.NewV4()
		if err != nil {
			logger.Wrapf(err, "failed to generate UUID:")
			uid, _ = uuid.NewV1()
		}
		// 设置cookie
		claim := middleware.UUIDClaim{
			UUID: uid.String(),
			StandardClaims: jwt.StandardClaims{
				ExpiresAt: time.Now().Add(time.Second * time.Duration(config.Global.CookieMaxAge())).Unix(),
			},
		}
		token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), claim)
		tt, err := token.SignedString([]byte(config.Global.Token()))
		if err != nil {
			logger.Wrapf(err, "toke signed string err")
		}
		ctx.SetCookie(consts.CookieUUIDV2, tt, config.Global.CookieMaxAge(), "/", "."+config.Global.Domain(), false, false)
		ctx.Set(consts.CookieUUIDV2, uid.String())
	}
}

