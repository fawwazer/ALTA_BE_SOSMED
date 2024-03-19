package user

import (
	"mime/multipart"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	// "github.com/labstack/echo"
)

type UserController interface {
	Add() echo.HandlerFunc
	Login() echo.HandlerFunc
	Profile() echo.HandlerFunc
	UploadPicture() echo.HandlerFunc
}

type UserService interface {
	Register(newData User) error
	Login(loginData User) (User, string, error)
	Profile(token *jwt.Token) (User, error)
	SaveUploadedFile(file *multipart.FileHeader, path string) error
	UploadPicture(userID int, pictureUrl string, token *jwt.Token) error
}

type UserModel interface {
	AddUser(newData User) error
	UpdateUser(email string, data User) error
	Login(email string) (User, error)
	GetUserByEmail(email string) (User, error)
	GetLastUserID() (int, error)
	UploadPictureURL(userID int, picture string) error
}

type User struct {
	UserID    int
	Nama      string
	Email     string
	Password  string
	Picture   string
	Tgl_lahir string
	Gender    string
	Alamat    string
}

type Login struct {
	Email    string `validate:"required"`
	Password string `validate:"required,alphanum,min=8"`
}

type Register struct {
	UserID    int
	Nama      string `validate:"required,alpha"`
	Email     string `validate:"required"`
	Password  string `validate:"required,alphanum,min=8"`
	Tgl_lahir string `validate:"required"`
	Gender    string `validate:"required"`
	Alamat    string `validate:"required,alphaunicode"`
}
