package data

import "ALTA_BE_SOSMED/features/post/data"

type User struct {
	UserID    int
	Nama      string
	Email     string `gorm:"primaryKey;type:varchar(30);"`
	Password  string
	Picture   string
	Tgl_lahir string
	Gender    bool
	Posts     []data.Post `gorm:"foreignKey:Pemilik;references:Email"`
}
