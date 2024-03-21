package comment

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type ComController interface {
	View() echo.HandlerFunc
	Add() echo.HandlerFunc
	Update() echo.HandlerFunc
	Delete() echo.HandlerFunc
}

type ComModels interface {
	GetComByOwner(pemilik string) ([]Comment, error)
	InsertCom(pemilik string, kegiatanBaru Comment) (Comment, error)
	Update(pemilik string, comID uint, data Comment) (Comment, error)
	DeleteCom(DeleteCom Comment) error
}

type ComService interface {
	View(token *jwt.Token) ([]Comment, error)
	Update(pemilik string, comID uint, inputCom Comment) error
	AddCom(pemilik *jwt.Token, kegiatanBaru Comment) (Comment, error)
	DeleteCom(deleteID Comment) error
}

type Comment struct {
	gorm.Model
	Comment  string
	Pemiliks string
}
type ComRequest struct {
	Comment  string
	Pemiliks string
}
type Update struct {
	gorm.Model
	Comment string
	Pemilik string
}

type DeleteCom struct {
	gorm.Model
}
