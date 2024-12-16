package controllers

import (
	"main/database"
	"main/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func GetAllPeminjaman(c *gin.Context) {
	var peminjaman []models.Peminjaman
	database.DB.Find(&peminjaman)
	c.JSON(http.StatusOK, gin.H{"data": peminjaman})
}

func GetAllMyPeminjaman(c *gin.Context) {
	user, _ := c.Get("user")
	var peminjaman []models.Peminjaman
	database.DB.Where("user_id = ?", user.(models.User).ID).Find(&peminjaman)
	c.JSON(http.StatusOK, gin.H{"data": peminjaman})
}

func GetPeminjamanById(c *gin.Context) {
	var peminjaman models.Peminjaman
	result := database.DB.First(&peminjaman, c.Param("id"))

	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Peminjaman not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": peminjaman})
}

func CreatePeminjaman(c *gin.Context) {
	user, _ := c.Get("user")

	var body struct {
		BukuID 			uint `json:"buku_id"`
		TanggalPinjam 	time.Time `json:"tanggal_pinjam"`
		BatasKembali 	time.Time `json:"batas_kembali"`
	}

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bad Request"})
		return
	}

	if body.BukuID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Buku ID cannot be empty"})
		return
	}

	if body.TanggalPinjam.IsZero() {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Tanggal Pinjam cannot be empty"})
		return
	}

	if body.BatasKembali.IsZero() {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Batas Kembali cannot be empty"})
		return
	}

	peminjaman := models.Peminjaman{
		UserID: user.(models.User).ID,
		BukuID: body.BukuID,
		TanggalPinjam: body.TanggalPinjam,
		BatasKembali: body.BatasKembali,
	}

	result := database.DB.Create(&peminjaman)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{})
}

func UpdatePeminjaman(c *gin.Context) {
	var body struct {
		TanggalKembali 		time.Time `json:"tanggal_kembali"`
	}

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bad Request"})
		return
	}

	if body.TanggalKembali.IsZero() {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Tanggal Kembali cannot be empty"})
		return
	}

	var peminjaman models.Peminjaman
	result := database.DB.First(&peminjaman, c.Param("id"))

	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Peminjaman not found"})
		return
	}
	peminjaman.TanggalKembali = body.TanggalKembali
	if peminjaman.TanggalKembali.After(peminjaman.BatasKembali) {
		peminjaman.Status = "Dikembalikan Terlambat"
		// tambah denda per hari terlambat + 2000
		denda := int(peminjaman.TanggalKembali.Sub(peminjaman.BatasKembali).Hours()/24) * 2000
		peminjaman.Denda = uint(denda)
	}

	result = database.DB.Save(&peminjaman)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong"})
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}

func DeletePeminjaman(c *gin.Context) {
	var peminjaman models.Peminjaman
	result := database.DB.First(&peminjaman, c.Param("id"))

	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Peminjaman not found"})
		return
	}

	result = database.DB.Delete(&peminjaman)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong"})
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}
