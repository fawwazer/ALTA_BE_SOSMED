package data

import (
	"ALTA_BE_SOSMED/features/comment/data"
)

type User struct {
	UserId    int
	Nama      string
	Email     string `gorm:"type:varchar(30);primaryKey"`
	Password  string
	Picture   string
	Tgl_lahir string
	Gender    string
	Comment   []data.Comment `gorm:"foreignKey:Pemiliks;references:Email"`
}
