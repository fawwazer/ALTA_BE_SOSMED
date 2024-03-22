package services

import (
	"ALTA_BE_SOSMED/features/post"
	"ALTA_BE_SOSMED/helper"
	"ALTA_BE_SOSMED/middlewares"
	"context"
	"errors"
	"log"
	"mime/multipart"

	"github.com/cloudinary/cloudinary-go"
	"github.com/cloudinary/cloudinary-go/api/uploader"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
)

type service struct {
	m post.PostModel
	v *validator.Validate
}

// SaveUploadedFile implements post.PostService.
func (s *service) SaveUploadedFile(file *multipart.FileHeader, path string) error {
	panic("unimplemented")
}

func NewPostService(model post.PostModel) post.PostService {
	return &service{
		m: model,
		v: validator.New(),
	}
}

func (s *service) AddPost(pemilik *jwt.Token, postingBaru post.Post, file *multipart.FileHeader) (post.Post, error) {
	// Decode token to obtain email
	email := middlewares.DecodeToken(pemilik)
	if email == "" {
		log.Println("error decode token:", "token tidak ditemukan")
		return post.Post{}, errors.New("data tidak valid")
	}

	// Validate the new post data
	if err := s.v.Struct(&postingBaru); err != nil {
		log.Println("error validasi", err.Error())
		return post.Post{}, err
	}

	// Upload image to Cloudinary
	uploadResp, err := uploadToCloudinary(file)
	if err != nil {
		log.Println("cloudinary upload error:", err.Error())
		return post.Post{}, err
	}

	// Set picture URL to Cloudinary image URL
	postingBaru.Picture = uploadResp.URL

	// Insert the new post into the database
	result, err := s.m.InsertPost(email, postingBaru)
	if err != nil {
		return post.Post{}, errors.New(helper.ServerGeneralError)
	}

	return result, nil
}

func uploadToCloudinary(file *multipart.FileHeader) (*uploader.UploadResult, error) {
	cld, err := cloudinary.NewFromURL("cloudinary://975476473639685:K7MSOZOWkrlRiX8rhm4ybiNRCkc@duwhyyuy8")
	if err != nil {
		return nil, err
	}
	return cld.Upload.Upload(context.Background(), file, uploader.UploadParams{})
}
