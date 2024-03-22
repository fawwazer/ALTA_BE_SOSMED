package post

import (
	"mime/multipart"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

type PostController interface {
	Add() echo.HandlerFunc
}

type PostModel interface {
	InsertPost(pemilik string, postingBaru Post) (Post, error)
	UpdatePost(pemilik string, todoID uint, data Post) (Post, error)
	GetPostByOwner(pemilik string) ([]Post, error)
}

type PostService interface {
	AddPost(pemilik *jwt.Token, postingBaru Post, file *multipart.FileHeader) (Post, error)
	SaveUploadedFile(file *multipart.FileHeader, path string) error
}

type Post struct {
	Posting string
	Picture string
}
