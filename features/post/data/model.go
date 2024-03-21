package data

import (
	"gorm.io/gorm"
)

type Post struct {
	gorm.Model
	PicturePost string
	Posting     string `gorm:type:text`
	Pemilik     string
}
