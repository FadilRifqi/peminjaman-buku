package controllers

import (
	"main/database"
	"main/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetUsers(c *gin.Context){
	var users []models.User
	result := database.DB.Find(&users)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": users})
}

func GetUserById(c *gin.Context){
	var user models.User
	result := database.DB.First(&user, c.Param("id"))

	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": user})
}

func UpdateUser(c *gin.Context){
    // Find the User in the Database with url parameter
    var user models.User
    result := database.DB.First(&user, c.Param("id"))

    if result.Error != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
        return
    }

    // Get Email and Username
    var body struct {
        Username string	`json:"username"`
        Email    string	`json:"email"`
    }

    // Validate Email and Username
    if c.Bind(&body) != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Bad Request"})
        return
    }

    // Update the User from body if body fields are not empty, otherwise use existing data
    if body.Username != "" {
        user.Username = body.Username
    }
    if body.Email != "" {
        user.Email = body.Email
    }

    result = database.DB.Save(&user)

    if result.Error != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong"})
        return
    }

    // Respond
    c.JSON(http.StatusOK, gin.H{"message": "User updated successfully"})
}

func DeleteUser(c *gin.Context){
	// Find the User in the Database with url parameter
	var user models.User
	result := database.DB.First(&user, c.Param("id"))

	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Delete the User from the Database
	result = database.DB.Delete(&user)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong"})
		return
	}

	// Respond
	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}
