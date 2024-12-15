package database

import "main/models"

func SyncDB(){
	DB.AutoMigrate(&models.User{})
	DB.AutoMigrate(&models.Buku{})
	DB.AutoMigrate(&models.Peminjaman{})
}
