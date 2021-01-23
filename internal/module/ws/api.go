package ws

import (
	"github.com/gin-gonic/gin"
	"github.com/teablog/tea/internal/middleware"
)

func Route(api *gin.RouterGroup) {
	t := api.Group("/ws", middleware.CookieUUID())
	{
		t.GET("/join", ServeWs)
	}
}
