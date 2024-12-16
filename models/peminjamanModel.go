package models

import (
	"time"

	"gorm.io/gorm"
)

type Peminjaman struct {
	gorm.Model
	UserID 			uint
	Peminjam   		User `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	BukuID 			uint
	Buku   			Buku `gorm:"foreignKey:BukuID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	TanggalPinjam 	time.Time
    TanggalKembali 	time.Time
    BatasKembali  	time.Time
	Status 			string `gorm:"type:varchar(255),default:'Dipinjam'"`
	Denda 			uint `gorm:"default:0"`
}


func (Peminjaman) TableName() string {
    return "peminjaman"
}
