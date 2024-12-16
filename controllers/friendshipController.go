package controllers

import (
	"main/database"
	"main/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func SendFriendRequest(c *gin.Context)  {
	user,_ := c.Get("user")

	friendID := c.Param("id")
	//Convert friendID to uint
	uintFriendID, err := strconv.ParseUint(friendID, 10, 32)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid friend ID"})
		return
	}

	//check if friend exists
	var friend models.User
	friendModel := database.DB.First(&friend, uint(uintFriendID))

	if friendModel.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Friend not found"})
		return
	}

	// Check if the user is trying to add themselves
	if user.(models.User).ID == uint(uintFriendID) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "You can't add yourself as a friend"})
		return
	}

	// Check if the friendship already exists
	var friendship models.Friendship
	result := database.DB.Where("user_id = ? AND friend_id = ?", user.(models.User).ID, uint(uintFriendID)).First(&friendship)

	if result.RowsAffected > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Already Sent a Friend Request"})
		return
	}


	// Create a new friendship
	friendship = models.Friendship{
		UserID: user.(models.User).ID,
		FriendID: uint(uintFriendID),
	}

	// Check if the friend already sent a friend request
	result = database.DB.Where("user_id = ? AND friend_id = ?", uint(uintFriendID), user.(models.User).ID).First(&friendship)

	if result.Error == nil {
		friendship.Status = "Accepted"
		database.DB.Save(&friendship)

		// Create a new friendship
		friendship = models.Friendship{
			UserID: user.(models.User).ID,
			FriendID: uint(uintFriendID),
			Status: "Accepted",
		}

		database.DB.Create(&friendship)

		c.JSON(http.StatusCreated, gin.H{})
		return
	}

	database.DB.Create(&friendship)

	c.JSON(http.StatusCreated, gin.H{})
}

func AcceptFriendRequest(c *gin.Context){
	user, _ := c.Get("user")

	friendRequestID := c.Param("id")
	// Convert friendRequestID to uint
	uintFriendRequestID, err := strconv.ParseUint(friendRequestID, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid friend request ID"})
		return
	}

	// Check if the friend request doesnt exists
	var friendRequest models.Friendship
	result := database.DB.Where("id = ? AND friend_id = ?", uint(uintFriendRequestID), user.(models.User).ID).First(&friendRequest)

	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Friend Request not found"})
		return
	}

	// Check if the friend request is already accepted
	if friendRequest.Status == "Accepted" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Friend Request already accepted"})
		return
	}

	// Update the friendship
	friendRequest.Status = "Accepted"
	database.DB.Save(&friendRequest)

	// Create a new friendship
	friendship := models.Friendship{
		UserID: user.(models.User).ID,
		FriendID: friendRequest.UserID,
		Status: "Accepted",
	}

	database.DB.Create(&friendship)

	c.JSON(http.StatusOK, gin.H{})
}

func DeleteFriend(c *gin.Context){
	user, _ := c.Get("user")

	friendID := c.Param("id")
	// Convert friendID to uint
	uintFriendID, err := strconv.ParseUint(friendID, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid friend ID"})
		return
	}

	// Check if the friend exists
	var friend models.User
	friendModel := database.DB.First(&friend, uint(uintFriendID))

	if friendModel.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Friend not found"})
		return
	}

	// Check if the friendship exists
	var friendship models.Friendship
	result := database.DB.Where("user_id = ? AND friend_id = ?", user.(models.User).ID, uint(uintFriendID)).First(&friendship)

	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Friendship not found"})
		return
	}

	// Check if the friend has a friendship with the user
	var friendFriendship models.Friendship
	fResult := database.DB.Where("user_id = ? AND friend_id = ?", uint(uintFriendID), user.(models.User).ID).First(&friendFriendship)

	if fResult.Error == nil {
		database.DB.Delete(&friendFriendship)
	}
	database.DB.Delete(&friendship)

	c.JSON(http.StatusOK, gin.H{})
}
