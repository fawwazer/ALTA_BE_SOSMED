package user

import "github.com/labstack/echo"

type UserController interface {
	Add() echo.HandlerFunc
	Login() echo.HandlerFunc
	Profile() echo.HandlerFunc
}

type User struct {
	Nama      string
	Email     string
	Password  string
	Tgl_lahir string
	gender    string
	alamat    string
}

type Login struct {
	Email    string `validate:"required,email"`
	Password string `validate:"required,alphanum,min=8"`
}

type Register struct {
	Nama      string `validate:"required,alpha"`
	Email     string `validate:"required,email"`
	Password  string `validate:"required,alphanum,min=8"`
	Tgl_lahir string
	gender    string
	alamat    string
}
