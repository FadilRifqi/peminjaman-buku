package controllers

import (
	"main/database"
	"main/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

//FIXME: How to get the room ID from the Frontend?

func SendChat(c *gin.Context) {
	var body struct{
		Message string `json:"message"`
		ReceiverID uint `json:"receiver_id"`
	}
	user, _ := c.Get("user")
	roomId := c.Param("id")
	var chat models.Chat

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if body.Message == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Message is required"})
		return
	}

	if body.ReceiverID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Receiver ID is required"})
		return
	}

	// check if the receiver exists
	var receiver models.User
	receiverResult := database.DB.First(&receiver, body.ReceiverID)
	if receiverResult.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Receiver not found"})
		return
	}

	// check if the room exists
	var room models.Room
	roomResult := database.DB.First(&room, roomId)
	if roomResult.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Room not found"})
		return
	}

	// parse the room id
	paresdRoomId,err := strconv.ParseUint(roomId, 10, 32)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid room ID"})
		return
	}

	chat.UserID = user.(models.User).ID
	chat.RoomID = uint(paresdRoomId)
	chat.Message = body.Message

	result := database.DB.Create(&chat)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send chat"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{})
}

func GetMyChats(c *gin.Context) {
	roomId := c.Param("id")

	var chats []models.Chat
	result := database.DB.Where("room_id = ?", roomId).Find(&chats)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get chats"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"chats": chats})
}
