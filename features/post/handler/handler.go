package handler

import (
	"ALTA_BE_SOSMED/features/post"
	"ALTA_BE_SOSMED/helper"
	"log"
	"net/http"
	"strings"

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
		var input PostRequest
		err := c.Bind(&input)
		if err != nil {
			log.Println("error bind data:", err.Error())
			if strings.Contains(err.Error(), "unsupport") {
				return c.JSON(http.StatusUnsupportedMediaType,
					helper.ResponseFormat(http.StatusUnsupportedMediaType, helper.UserInputFormatError, nil))
			}
			return c.JSON(http.StatusBadRequest,
				helper.ResponseFormat(http.StatusBadRequest, helper.UserInputError, nil))
		}

		token, ok := c.Get("user").(*jwt.Token)
		if !ok {
			return c.JSON(http.StatusBadRequest,
				helper.ResponseFormat(http.StatusBadRequest, helper.UserInputError, nil))
		}

		// Retrieve the uploaded file from the request.
		file, err := c.FormFile("image")
		if err != nil && err != http.ErrMissingFile { // Check if error is not due to missing file
			log.Println("error form file: ", err.Error())
			return c.JSON(http.StatusBadRequest,
				helper.ResponseFormat(http.StatusBadRequest, "Invalid data! The data type must be images!", nil))
		}

		if file != nil { // Check if file exists
			// Define the file path to save the uploaded image.

			// Save the uploaded file to the specified path.
			// if err := ct.s.SaveUploadedFile(file, ""); err != nil {
			// 	log.Print("error save uploaded file: ", err.Error())
			// 	return c.JSON(http.StatusInternalServerError,
			// 		helper.ResponseFormat(http.StatusInternalServerError, helper.ServerGeneralError, nil))
			// }

			// Construct the URL for the saved picture.

		}

		var inputProcess post.Post
		inputProcess.Posting = input.Posting

		result, err := ct.s.AddPost(token, inputProcess, file)
		if err != nil {
			log.Println("error insert db:", err.Error())
			return c.JSON(http.StatusInternalServerError, helper.ResponseFormat(http.StatusInternalServerError, helper.ServerGeneralError, nil))
		}

		return c.JSON(http.StatusCreated, helper.ResponseFormat(http.StatusCreated, "berhasil menambahkan kegiatan", result))
	}
}
