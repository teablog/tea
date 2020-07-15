package controllers

import (
	"dyc/internal/config"
	"dyc/internal/middleware"
	"dyc/internal/module/chat"
	"dyc/internal/module/util"
	"github.com/gin-gonic/gin"
	"net/http"
	"path"
)

func Init(engine *gin.Engine) {
	util.Init()
}

func NewRouter(router *gin.Engine) {
	hub := chat.NewHub()
	go hub.Run()
	storageDir := config.GetKey("path::storage_dir").String()
	api := router.Group("/tea")
	{
		// 文章
		api.GET("/articles", Article.List)
		api.GET("/articles/labels", Article.Labels)
		api.GET("/article/:id", Article.View)
		api.GET("/topic/:topic", Topic.List)
		api.GET("/search/articles", Article.Search)
		api.POST("/subscribe", Subscribe.Create)
		// 电影资源
		api.GET("/media/subtype/:subtype", Media.Index)
		api.GET("/search/media", Media.Search)
		api.GET("/video/:id", Media.View)
		api.GET("/oauth/github", Oauth.Github)
		api.POST("/oauth/google", Oauth.Google)
		// 工具
		utils := api.Group("/utils")
		{
			// 测试接口
			utils.GET("/preserve_host", Util.PreserveHost)
			utils.GET("/weather", Util.Weather)
			// ip 地址解析
			utils.GET("/ip/position", Util.Ip)
			// 地区
			region := utils.Group("/region")
			{
				region.GET("/amap", Util.Amap)
				region.GET("/location", Util.Location)
			}
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
