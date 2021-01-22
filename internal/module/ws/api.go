package ws

import "github.com/gin-gonic/gin"

func Route(api *gin.RouterGroup)  {
	t := api.Group("/ws")
	{
		t.GET("/join", ServeWs)
	}
}
