package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username 	string  `gorm:"type:varchar(255)"`
	Email 		string 	`gorm:"unique;type:varchar(255)"`
	Password 	string 	`gorm:"type:varchar(255)"`
	RoleID 		uint 	`gorm:"default:2"`
	Role 		Role 	`gorm:"foreignKey:RoleID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}
