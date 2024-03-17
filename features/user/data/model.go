package data

type User struct {
	Nama      string
	Email     string `gorm:"type:varchar(100);primaryKey"`
	Password  string
	Tgl_lahir string
	Gender    string
	Alamat    string
}
