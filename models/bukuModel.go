package models

import "gorm.io/gorm"

type Buku struct {
	gorm.Model
	Judul  		string `gorm:"type:varchar(255)"`
	Penulis 	string `gorm:"type:varchar(255)"`
	Tahun  		int
	UserID 		uint
	Pemilik 	User `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}
