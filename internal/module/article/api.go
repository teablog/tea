package article

import "github.com/gin-gonic/gin"

func Route(api *gin.RouterGroup) {
	as := api.Group("/articles")
	{
		as.GET("", Art.List)
		as.GET("/labels", Art.Labels)
		as.GET("/search", Search.List)
	}
	a := api.Group("/article")
	{
		a.GET("", Art.Get)
		a.GET("/messages", Msg.FindMessages)
		a.POST("/message", Msg.Comment)
	}
	t := api.Group("/topic")
	{
		t.GET("", Topics.List)
	}
}
