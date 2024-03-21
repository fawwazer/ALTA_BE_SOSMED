package data

import "ALTA_BE_SOSMED/features/post"

type User struct {
	UserID    int
	Nama      string
	Email     string `gorm:"primaryKey;type:varchar(30);"`
	Password  string
	Picture   string
	Tgl_lahir string
	Gender    bool
	Posts     []post.Post `gorm:"foreignKey:Pemilik;references:Email"`
}
