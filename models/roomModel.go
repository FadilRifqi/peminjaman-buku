package models

import "gorm.io/gorm"

type Room struct {
	gorm.Model
	Label string `gorm:"type:varchar(255);unique"`
	RoomMembers []RoomMember
	Chats []Chat
}
