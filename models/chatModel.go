package models

import "gorm.io/gorm"

type Chat struct {
	gorm.Model
	RoomID 		uint
	Room		Room `gorm:"foreignKey:RoomID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	UserID 		uint
	Sender		User `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Message 	string
}
