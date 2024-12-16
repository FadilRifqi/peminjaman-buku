package routes

import (
	"main/controllers"
	"main/middlewares"

	"github.com/gin-gonic/gin"
)

func PeminjamanRoutes(router *gin.Engine){
	peminjaman:=router.Group("/peminjaman")
	peminjaman.Use(middlewares.RequireAuth)
	{
		peminjaman.POST("/add",controllers.CreatePeminjaman)
		peminjaman.GET("/all",controllers.GetAllPeminjaman)
		peminjaman.GET("/:id",middlewares.PeminjamanMiddleware,controllers.GetPeminjamanById)
		peminjaman.GET("/me",controllers.GetAllMyPeminjaman)
		peminjaman.PATCH("/update/:id",middlewares.PeminjamanMiddleware,controllers.UpdatePeminjaman)
		peminjaman.DELETE("/delete/:id",middlewares.PeminjamanMiddleware,controllers.DeletePeminjaman)
	}
}
