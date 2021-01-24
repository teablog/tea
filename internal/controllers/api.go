package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/teablog/tea/internal/middleware"
	"github.com/teablog/tea/internal/module/account"
	"github.com/teablog/tea/internal/module/article"
	"github.com/teablog/tea/internal/module/tools"
	"github.com/teablog/tea/internal/module/tracke"
	"github.com/teablog/tea/internal/module/ws"
	"net/http"
)

func Init(engine *gin.Engine) {
	tools.Init()
}

func NewRouter(router *gin.Engine) {
	api := router.Group("/api")
	{
		// websocket
		ws.Route(api)
		// 数据上报
		tracke.Route(api)
		// 文章
		article.Route(api)
		// 账户
		account.Route(api)
		// 工具
		tool := api.Group("/tools")
		{
			// ip 地址解析
			tool.GET("/location/ip", Tools.Ip)
			tool.GET("/location/latitude-longitude", Tools.Amap)
			tool.GET("/location", Tools.Location)
		}
		// 帮助中心
		help := api.Group("/helper")
		{
			help.GET("/token", Help.Token)
		}
		api.GET("/seo/sitemap", Seo.SiteMap)
		// 需要简单鉴权的API
		tokenAuth := api.Group("/", middleware.TokenAuthCheck())
		{
			tokenAuth.POST("/logstash/collect", Logstash.Collect)
		}
	}
	router.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "OK")
	})
}
