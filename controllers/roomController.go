package controllers

import (
	"main/database"
	"main/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetRoomIDFromLabel(c *gin.Context) {
	var room models.Room
	label := c.Param("label")
	result := database.DB.Where("label = ?", label).First(&room)

	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Room not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": room.ID})
}
