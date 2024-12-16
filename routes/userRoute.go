package routes

import (
	"main/controllers"
	"main/middlewares"

	"github.com/gin-gonic/gin"
)

func UserRoutes(router *gin.Engine) {
	user := router.Group("/user")
	user.Use(middlewares.RequireAuth)
	{
		user.GET("/all", controllers.GetUsers)
		user.GET("/:id", controllers.GetUserById)
		user.PATCH("/update/:id", middlewares.IsMe ,controllers.UpdateUser)
		user.DELETE("/delete/:id", controllers.DeleteUser)
	}
	auth := router.Group("/auth")
	{
		auth.POST("/register", controllers.Register)
		auth.POST("/login", controllers.GenerateToken)
		auth.GET("/validate",middlewares.RequireAuth ,controllers.Validate)
		auth.GET("/check",middlewares.RequireAuth, middlewares.Admin ,controllers.Validate)
	}
}
