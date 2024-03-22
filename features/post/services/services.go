package services

import (
	"ALTA_BE_SOSMED/features/post"
	"ALTA_BE_SOSMED/helper"
	"ALTA_BE_SOSMED/middlewares"
	"context"
	"errors"
	"io"
	"log"
	"mime/multipart"
	"os"

	"github.com/cloudinary/cloudinary-go"
	"github.com/cloudinary/cloudinary-go/api/uploader"
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

func (s *service) AddPost(pemilik *jwt.Token, postingBaru post.Post, file *multipart.FileHeader) (post.Post, error) {
	// Decode token to obtain email
	email := middlewares.DecodeToken(pemilik)
	if email == "" {
		log.Println("error decode token:", "token tidak ditemukan")
		return post.Post{}, errors.New("data tidak valid")
	}

	// Validate the new post data
	err := s.v.Struct(&postingBaru)
	if err != nil {
		log.Println("error validasi", err.Error())
		return post.Post{}, err
	}

	// Save the uploaded file if it exists
	if file != nil { // Check if file exists
		cld, err := cloudinary.NewFromURL("cloudinary://426244812151882:GBqN6L8Rm77iHHkPXiemVPP_e2Y@dlosajdpy")
		if err != nil {
			log.Print("error connect error: ", err.Error())
			return post.Post{}, err
		}
		src, err := file.Open()
		if err != nil {
			return post.Post{}, err
		}
		defer src.Close()
		resp, err := cld.Upload.Upload(context.Background(), src, uploader.UploadParams{})
		if err != nil {
			log.Print("error upload error: ", err.Error())
			return post.Post{}, err
		}
		postingBaru.Picture = resp.SecureURL
	}
	log.Println(postingBaru.Picture)
	// Set the picture URL in the new post

	// Insert the new post into the database
	result, err := s.m.InsertPost(email, postingBaru)
	if err != nil {
		return post.Post{}, errors.New(helper.ServerGeneralError)
	}

	return result, nil
}

func (s *service) SaveUploadedFile(file *multipart.FileHeader, path string) error {
	// Open the uploaded file.
	src, err := file.Open()
	if err != nil {
		log.Print("file open error :", err.Error())
		return err
	}
	defer src.Close()

	// Create a destination file for the uploaded content.
	dst, err := os.Create(path)
	if err != nil {
		log.Print("file create error :", err.Error())
		return err
	}
	defer dst.Close()

	// Copy the uploaded content to the destination file.
	if _, err = io.Copy(dst, src); err != nil {
		log.Print("file copy error :", err.Error())
		return err
	}

	return nil
}
