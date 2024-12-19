package routes

import (
	"main/controllers"

	"github.com/gin-gonic/gin"
)

func WebSocketRoutes(router *gin.Engine) {
	websocket := router.Group("/ws")
	{
		websocket.GET("/:id", controllers.HandleWebSocket)
	}
}
