package models

import "gorm.io/gorm"

type Chat struct {
	gorm.Model
	RoomID 		uint
	Room		Room `gorm:"foreignKey:RoomID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	UserID 		uint
	Message 	string
}
