package middleware

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/teablog/tea/internal/config"
	"github.com/teablog/tea/internal/consts"
	"github.com/teablog/tea/internal/helper"
	"github.com/teablog/tea/internal/logger"
)

type UUIDClaim struct {
	UUID string
	jwt.StandardClaims
}

func CookieUUID() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ck, _ := ctx.Cookie(consts.CookieUUIDV2)
		if ck != "" {
			claim := new(UUIDClaim)
			token, err := jwt.ParseWithClaims(ck, claim, func(token *jwt.Token) (interface{}, error) {
				return []byte(config.Global.Token()), nil
			})
			if err != nil || !token.Valid {
				logger.Wrapf(err, "jwt parse token err")
				helper.FailWithCode(ctx, errors.New("invalid token"), 401)
				ctx.Abort()
				return
			}
			logger.Debugf("uuid: %s", claim.UUID)
			ctx.Set(consts.CookieUUIDV2, claim.UUID)
		} else {
			helper.FailWithCode(ctx, errors.New("invalid token"), 401)
			ctx.Abort()
			return
		}
		ctx.Next()
	}
}
