package models

import "gorm.io/gorm"

type Role struct {
	gorm.Model
	Name string `gorm:"type:varchar(255)"`
}
