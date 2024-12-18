package controllers

import (
	"main/database"
	"main/models"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func GetRoomIDFromLabel(c *gin.Context) {
    var room models.Room
    label := c.Param("label")
    result := database.DB.Where("label = ?", label).First(&room)

    if result.Error != nil {
        // If no room is found, try the reversed label
        parts := strings.Split(label, "-")
        if len(parts) == 2 {
            reversedLabel := parts[1] + "-" + parts[0]
            result = database.DB.Where("label = ?", reversedLabel).First(&room)
            if result.Error != nil {
                c.JSON(http.StatusNotFound, gin.H{"error": "Room not found"})
                return
            }
        } else {
            c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid label format"})
            return
        }
    }

    c.JSON(http.StatusOK, gin.H{"data": room.ID})
}
