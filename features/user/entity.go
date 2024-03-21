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
	UpdateProfile() echo.HandlerFunc
	DeleteAccount() echo.HandlerFunc
}

type UserService interface {
	Register(newData User) error
	Login(loginData User) (User, string, error)
	Profile(token *jwt.Token) (User, error)
	SaveUploadedFile(file *multipart.FileHeader, path string) error
	UpdateProfile(userID int, token *jwt.Token, newData User) error
	DeleteAccount(userID uint, token *jwt.Token) error
}

type UserModel interface {
	AddUser(newData User) error
	UpdateUser(email string, data User) error
	Login(email string) (User, error)
	GetUserByEmail(email string) (User, error)
	GetLastUserID() (int, error)
	Update(userID int, newData User) error
	Delete(userID uint) error
}

type User struct {
	UserID    int
	Nama      string
	Email     string
	Password  string
	Picture   string
	Tgl_lahir string
	Gender    bool
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
	Gender    bool   `validate:"required"`
}
