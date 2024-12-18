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

	// Find Room with Label = userID-friendID
	var room models.Room
	roomResult := database.DB.Where("label = ?", strconv.Itoa(int(user.(models.User).ID)) + "-" + strconv.Itoa(int(friendRequest.UserID))).First(&room)

	if roomResult.Error != nil {
		// Create New Room
		room = models.Room{
			Label: strconv.Itoa(int(user.(models.User).ID)) + "-" + strconv.Itoa(int(friendRequest.UserID)),
		}
		database.DB.Create(&room)
		// Create Room Member
		roomMember := models.RoomMember{
			RoomID: room.ID,
			UserID: user.(models.User).ID,
		}
		database.DB.Create(&roomMember)

		roomMember = models.RoomMember{
			RoomID: room.ID,
			UserID: friendRequest.UserID,
		}
		database.DB.Create(&roomMember)
		return
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

func GetMyFriends(c *gin.Context) {
    user, _ := c.Get("user")

    var friendships []models.Friendship
    database.DB.Where("user_id = ? AND status = ?", user.(models.User).ID, "Accepted").Find(&friendships)

    type UserResponse struct {
        ID          	uint   `json:"id"`
        Username    	string `json:"username"`
        Email       	string `json:"email"`
        LastMessage  	string `json:"lastMessage"`
		Time 	  		string `json:"time"`
    }

    var friends []UserResponse
    for _, friendship := range friendships {
        var friend models.User
        database.DB.First(&friend, friendship.FriendID)

        friendResponse := UserResponse{
            ID:       friend.ID,
            Username: friend.Username,
            Email:    friend.Email,
        }

        // Search for the room with either user-friend or friend-user label
        var room models.Room
        userID := strconv.Itoa(int(user.(models.User).ID))
        friendID := strconv.Itoa(int(friendship.FriendID))
        label1 := userID + "-" + friendID
        label2 := friendID + "-" + userID

        roomResult := database.DB.Where("label = ? OR label = ?", label1, label2).First(&room)
        if roomResult.Error == nil {
            // Find the latest chat message for that room
            var LatestChat models.Chat
            database.DB.Where("room_id = ?", room.ID).Order("created_at desc").First(&LatestChat)
            friendResponse.LastMessage = LatestChat.Message
			friendResponse.Time = LatestChat.CreatedAt.Format("2006-01-02 15:04:05")
        }

        friends = append(friends, friendResponse)
    }

    c.JSON(http.StatusOK, gin.H{"data": friends})
}
