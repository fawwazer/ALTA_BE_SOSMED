package data

import (
	"gorm.io/gorm"
)

type Post struct {
	gorm.Model
	PicturePost string
	Posting     string `gorm:"type:string"`
	Pemilik     string `gorm:"type:varchar(30);"`
}
