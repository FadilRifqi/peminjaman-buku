package controllers

import (
	"main/database"
	"main/models"
	"net/http"

	"github.com/gin-gonic/gin"
)


func AddBuku(c *gin.Context){
	user, _ := c.Get("user")

    var body struct {
        Judul   string `json:"judul"`
        Penulis string `json:"penulis"`
        Tahun   int    `json:"tahun"`
        UserID  uint
    }

    if err := c.Bind(&body); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Bad Request"})
        return
    }

    // Manually validate each field
    if body.Judul == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Field judul is required"})
        return
    }
    if body.Penulis == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Field penulis is required"})
        return
    }
    if body.Tahun == 0 {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Field tahun is required"})
        return
    }

    // Create a new Buku
	buku := models.Buku{
		Judul:   body.Judul,
		Penulis: body.Penulis,
		Tahun:   body.Tahun,
		UserID:  user.(models.User).ID,
	}

	// Save the Buku in the Database
	result := database.DB.Create(&buku)


	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong"})
		return
	}


    c.JSON(http.StatusCreated, gin.H{})
}

func GetBukuById(c *gin.Context){
	var buku models.Buku
	result := database.DB.First(&buku, c.Param("id"))

	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Buku not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": buku})
}

func GetAllBukus(c *gin.Context){
	var bukus []models.Buku
	result := database.DB.Find(&bukus)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": bukus})
}

func GetMyBukus(c *gin.Context){
	user, _ := c.Get("user")

	var bukus []models.Buku
	result := database.DB.Where("user_id = ?", user.(models.User).ID).Find(&bukus)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": bukus})
}

func UpdateBuku(c *gin.Context){
	var buku models.Buku
	result := database.DB.First(&buku, c.Param("id"))

	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Buku not found"})
		return
	}

	var body struct {
		Judul   string `json:"judul"`
		Penulis string `json:"penulis"`
		Tahun   int    `json:"tahun"`
	}

	if err := c.Bind(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bad Request"})
		return
	}

	if body.Judul != "" {
		buku.Judul = body.Judul
	}
	if body.Penulis != "" {
		buku.Penulis = body.Penulis
	}
	if body.Tahun != 0 {
		buku.Tahun = body.Tahun
	}

	result = database.DB.Save(&buku)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong"})
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}

func DeleteBuku(c *gin.Context){
	var buku models.Buku
	result := database.DB.First(&buku, c.Param("id"))

	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{})
		return
	}

	database.DB.Delete(&buku)

	c.JSON(http.StatusOK, gin.H{})
}
