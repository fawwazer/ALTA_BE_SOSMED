package handler

import (
	"ALTA_BE_SOSMED/features/post"
	"ALTA_BE_SOSMED/helper"
	"strings"

	"log"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

type controller struct {
	s post.PostService
}

func NewHandler(service post.PostService) post.PostController {
	return &controller{
		s: service,
	}
}

func (ct *controller) Add() echo.HandlerFunc {
	return func(c echo.Context) error {
		// Bind post data from request
		var input PostRequest
		err := c.Bind(&input)
		if err != nil {
			log.Println("error bind data:", err.Error())
			return c.JSON(http.StatusBadRequest,
				helper.ResponseFormat(http.StatusBadRequest, helper.UserInputError, nil))
		}

		// Retrieve JWT token from request context
		token, ok := c.Get("user").(*jwt.Token)
		if !ok {
			return c.JSON(http.StatusBadRequest, helper.ResponseFormat(http.StatusBadRequest, helper.UserInputError, nil))
		}

		// Retrieve uploaded image file from request
		file, err := c.FormFile("picture")
		if err != nil && err != http.ErrMissingFile {
			log.Println("error form file:", err.Error())
			return c.JSON(http.StatusBadRequest,
				helper.ResponseFormat(http.StatusBadRequest, "Invalid data! Please provide a valid picture file.", nil))
		}

		// Check if file is nil or unsupported
		if file == nil {
			return c.JSON(http.StatusBadRequest, helper.ResponseFormat(http.StatusBadRequest, "No file uploaded!", nil))
		}

		// Create new post object
		var inputPost post.Post
		inputPost.Posting = input.Posting

		// Call service to add post with Cloudinary file upload response
		result, err := ct.s.AddPost(token, inputPost, file)
		if err != nil {
			log.Println("error insert db:", err.Error())
			if strings.Contains(err.Error(), "unsupported file type") {
				return c.JSON(http.StatusBadRequest, helper.ResponseFormat(http.StatusBadRequest, "Unsupported file type! Please upload a valid image file.", nil))
			}
			return c.JSON(http.StatusInternalServerError, helper.ResponseFormat(http.StatusInternalServerError, helper.ServerGeneralError, nil))
		}

		// Return success response with newly added post
		return c.JSON(http.StatusCreated, helper.ResponseFormat(http.StatusCreated, "berhasil menambahkan kegiatan", result))
	}
}
