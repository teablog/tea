package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/teablog/tea/internal/config"
	"github.com/teablog/tea/internal/middleware"
	"github.com/teablog/tea/internal/module/chat"
	"github.com/teablog/tea/internal/module/tools"
	"net/http"
	"path"
)

func Init(engine *gin.Engine) {
	tools.Init()
}

func NewRouter(router *gin.Engine) {
	hub := chat.NewHub()
	go hub.Run()
	storageDir := config.GetKey("path::storage_dir").String()
	api := router.Group("/api")
	{
		// 文章
		api.GET("/articles", Article.List)
		api.GET("/articles/labels", Article.Labels)
		api.GET("/article/:id", Article.View)
		api.GET("/topic/:topic", Topic.List)
		api.GET("/search/articles", Article.Search)
		api.POST("/subscribe", Subscribe.Create)
		// 账户
		api.GET("/oauth/github", Oauth.Github)
		api.POST("/oauth/google", Oauth.Google)
		// 工具
		tool := api.Group("/tools")
		{
			// ip 地址解析
			tool.GET("/location/ip", Tools.Ip)
			tool.GET("/location", Tools.Location)
			tool.GET("/location/latitude-longitude", Tools.Amap)
		}
		// websocket
		auth := api.Group("/", middleware.LoginCheck())
		{
			auth.GET("/ws/join", func(context *gin.Context) {
				WS.Join(context, hub)
			})
			auth.POST("/ws/channel", Channel.Create)
			auth.GET("/ws/channel/subscribe", Channel.subscribe)
			auth.GET("/ws/channel/messages", Channel.Messages)
			auth.GET("/account/list", Account.List)
		}
		api.GET("/seo/sitemap", Seo.SiteMap)
	}
	router.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "OK")
	})
	// 静态文件
	router.Static("/images", path.Join(storageDir, "images"))
	router.StaticFile("/sitemap.xml", path.Join(storageDir, "seo"))
	router.StaticFile("/robots.txt", storageDir)
	router.StaticFile("/logo.png", storageDir)
	router.Static("/ext_dict", path.Join(storageDir, "ext_dict"))
}
