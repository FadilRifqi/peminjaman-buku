package middlewares

import (
	"main/database"
	"main/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func IsRoomMember(c *gin.Context){
	// Retrieve the user from the context
	user, _ := c.Get("user")

	u, ok := user.(models.User)
	if !ok {
		c.AbortWithStatus(http.StatusForbidden)
		return
	}

	// Retrieve the room id from the context
	roomId := c.Param("id")

	// Check if the user is a member of the room
	var room models.Room
	result := database.DB.First(&room, roomId)
	if result.Error != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "Room not found"})
		return
	}

	// Check if the user is a member of the room
	var roomMember models.RoomMember
	result = database.DB.Where("room_id = ? AND user_id = ?", roomId, u.ID).First(&roomMember)
	if result.Error != nil {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
		return
	}

	c.Next()
}
