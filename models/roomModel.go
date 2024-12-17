package models

import "gorm.io/gorm"

type Room struct {
	gorm.Model
	Name string 	`gorm:"type:varchar(255)"`
	Label string 	`gorm:"type:varchar(255);unique"`
	RoomMembers 	[]RoomMember
	Chats 			[]Chat
}
