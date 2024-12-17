package database

import "main/models"

func SyncDB(){
	DB.AutoMigrate(&models.User{})
	DB.AutoMigrate(&models.Buku{})
	DB.AutoMigrate(&models.Peminjaman{})
	DB.AutoMigrate(&models.Room{})
	DB.AutoMigrate(&models.Friendship{})
	DB.AutoMigrate(&models.Chat{})
	DB.AutoMigrate(&models.RoomMember{})
}
