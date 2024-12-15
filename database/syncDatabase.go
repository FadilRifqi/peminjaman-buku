package database

import "main/models"

func SyncDB(){
	DB.AutoMigrate(&models.User{})
}
