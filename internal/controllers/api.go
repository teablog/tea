package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/teablog/tea/internal/middleware"
	"github.com/teablog/tea/internal/module/account"
	"github.com/teablog/tea/internal/module/tools"
	"github.com/teablog/tea/internal/module/ws"
	"net/http"
)

func Init(engine *gin.Engine) {
	tools.Init()
}

func NewRouter(router *gin.Engine) {
	hub := ws.NewHub()
	go hub.Run()
	api := router.Group("/api")
	{
		// 文章
		api.GET("/articles", Article.List)
		api.GET("/articles/labels", Article.Labels)
		api.GET("/article/:id", Article.View)
		api.GET("/topic/:topic", Topic.List)
		api.GET("/search/articles", Article.Search)
		api.POST("/subscribe", Subscribe.Create)
		// 授权
		//api.GET("/oauth/github", Oauth.Github)
		//api.POST("/oauth/google", Oauth.Google)
		// 账户
		acct := api.Group("/account")
		{
			acct.GET("/name/exists", account.NameExists)
			acct.POST("/register", account.Register)
		}
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
		// websocket
		wss := api.Group("/", middleware.CookieUUID())
		{
			wss.GET("/ws/join", func(ctx *gin.Context) {
				ws.ServeWs(ctx, hub)
			})
		}
		api.GET("/ws/article/messages", Article.Messages)
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
