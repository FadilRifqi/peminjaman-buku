package models

import "gorm.io/gorm"

type Friendship struct {
	gorm.Model
	Status		string `gorm:"type:varchar(255);default:'Pending'"`
	UserID 		uint
	FriendID 	uint
	User		User `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Friend		User `gorm:"foreignKey:FriendID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}

func (Friendship) TableName() string {
    return "friendship"
}
