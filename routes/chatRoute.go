package routes

import (
	"main/controllers"
	"main/middlewares"

	"github.com/gin-gonic/gin"
)

func ChatRoutes(router *gin.Engine) {
	chatRoutes := router.Group("/chat")
	chatRoutes.Use(middlewares.RequireAuth)
	{
		chatRoutes.GET("/:id",middlewares.IsRoomMember ,controllers.GetMyChats)
		chatRoutes.POST("/send/:id",middlewares.IsRoomMember ,controllers.SendChat)
	}
}

func RoomRoutes(router *gin.Engine) {
	roomRoutes := router.Group("/room")
	roomRoutes.Use(middlewares.RequireAuth)
	{
		roomRoutes.GET("/:label",controllers.GetRoomIDFromLabel)
	}
}
