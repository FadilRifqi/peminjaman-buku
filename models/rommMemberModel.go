package models

import "gorm.io/gorm"

type RoomMember struct {
	gorm.Model
	RoomID uint
	UserID uint
	Room   Room `gorm:"foreignKey:RoomID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	User   User `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}

func (RoomMember) TableName() string {
	return "room_member"
}
