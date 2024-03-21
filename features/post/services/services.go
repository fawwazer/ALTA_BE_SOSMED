package services

import (
	"ALTA_BE_SOSMED/features/post"
	"ALTA_BE_SOSMED/helper"
	"ALTA_BE_SOSMED/middlewares"
	"errors"
	"log"

	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
)

type service struct {
	m post.PostModel
	v *validator.Validate
}

func NewPostService(model post.PostModel) post.PostService {
	return &service{
		m: model,
		v: validator.New(),
	}
}

func (s *service) AddPost(pemilik *jwt.Token, postingBaru post.Post) (post.Post, error) {
	email := middlewares.DecodeToken(pemilik)
	if email == "" {
		log.Println("error decode token:", "token tidak ditemukan")
		return post.Post{}, errors.New("data tidak valid")
	}

	err := s.v.Struct(&postingBaru)
	if err != nil {
		log.Println("error validasi", err.Error())
		return post.Post{}, err
	}

	result, err := s.m.InsertPost(email, postingBaru)
	if err != nil {
		return post.Post{}, errors.New(helper.ServerGeneralError)
	}

	return result, nil
}
