package handler

import (
	"ALTA_BE_SOSMED/features/user"
	"ALTA_BE_SOSMED/helper"
	"net/http"
	"strconv"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	// "github.com/labstack/echo"
)

type controller struct {
	service user.UserService
}

func NewUserHandler(s user.UserService) user.UserController {
	return &controller{
		service: s,
	}
}

func (ct *controller) Add() echo.HandlerFunc {
	return func(c echo.Context) error {
		var input user.User
		err := c.Bind(&input)
		if err != nil {
			if strings.Contains(err.Error(), "unsupport") {
				return c.JSON(http.StatusUnsupportedMediaType,
					helper.ResponseFormat(http.StatusUnsupportedMediaType, "format data tidak didukung", nil))
			}
			return c.JSON(http.StatusBadRequest,
				helper.ResponseFormat(http.StatusBadRequest, "data yang dikirmkan tidak sesuai", nil))
		}
		err = ct.service.Register(input)
		if err != nil {
			var code = http.StatusInternalServerError
			if strings.Contains(err.Error(), "validation") {
				code = http.StatusBadRequest
			}
			return c.JSON(code,
				helper.ResponseFormat(code, err.Error(), nil))
		}
		return c.JSON(http.StatusCreated,
			helper.ResponseFormat(http.StatusCreated, "selamat data sudah terdaftar", nil))
	}
}

func (ct *controller) Login() echo.HandlerFunc {
	return func(c echo.Context) error {
		var input LoginRequest
		err := c.Bind(&input)
		if err != nil {
			if strings.Contains(err.Error(), "unsupport") {
				return c.JSON(http.StatusUnsupportedMediaType,
					helper.ResponseFormat(http.StatusUnsupportedMediaType, "format data tidak didukung", nil))
			}
			return c.JSON(http.StatusBadRequest,
				helper.ResponseFormat(http.StatusBadRequest, "data yang dikirmkan tidak sesuai", nil))
		}

		var processData user.User
		processData.Email = input.Email
		processData.Password = input.Password

		result, token, err := ct.service.Login(processData)
		if err != nil {
			var code = http.StatusInternalServerError
			if strings.Contains(err.Error(), "validation") || strings.Contains(err.Error(), "cek kembali") {
				code = http.StatusBadRequest
			}
			return c.JSON(code,
				helper.ResponseFormat(code, err.Error(), nil))
		}

		var responseData LoginResponse
		responseData.Email = result.Email
		responseData.Nama = result.Nama
		responseData.Token = token

		return c.JSON(http.StatusOK,
			helper.ResponseFormat(http.StatusOK, "berhasil login", responseData))

	}
}

func (ct *controller) Profile() echo.HandlerFunc {
	return func(c echo.Context) error {
		token, ok := c.Get("user").(*jwt.Token)
		if !ok {
			return c.JSON(http.StatusBadRequest,
				helper.ResponseFormat(http.StatusBadRequest, helper.UserInputError, nil))
		}
		result, err := ct.service.Profile(token)
		if err != nil {
			var code = http.StatusInternalServerError
			if strings.Contains(err.Error(), "validation") || strings.Contains(err.Error(), "cek kembali") {
				code = http.StatusBadRequest
			}
			return c.JSON(code,
				helper.ResponseFormat(code, err.Error(), nil))
		}

		return c.JSON(http.StatusOK,
			helper.ResponseFormat(http.StatusOK, "berhasil mendapatkan data", result))
	}
}

func (ct *controller) UploadPicture() echo.HandlerFunc {
	return func(c echo.Context) error {
		userId, _ := strconv.ParseInt(c.Param("user_id"), 10, 32)

		token, ok := c.Get("user").(*jwt.Token)
		if !ok {
			return c.JSON(http.StatusBadRequest,
				helper.ResponseFormat(http.StatusBadRequest, helper.UserInputError, nil))
		}

		// Retrieve the uploaded file from the request.
		file, err := c.FormFile("image")
		if err != nil {
			return c.JSON(http.StatusBadRequest,
				helper.ResponseFormat(http.StatusBadRequest, "Invalid data! The data type must be images!", nil))
		}

		// Define the file path to save the uploaded image.
		pathImage := "path/to/your/project-profile/picture" + file.Filename

		// Save the uploaded file to the specified path.
		if err := ct.service.SaveUploadedFile(file, pathImage); err != nil {
			return c.JSON(http.StatusInternalServerError,
				helper.ResponseFormat(http.StatusInternalServerError, helper.ServerGeneralError, nil))
		}

		// Construct the URL for the saved picture.
		baseURL := "http://localhost:8000"
		pictureURL := baseURL + "/picture/" + file.Filename

		// // Update the user's profile with the picture URL using the user service.
		// if err := c.userService.UpdatePictureURL(userID, pictureURL); err != nil {
		//    return e.JSON(http.StatusInternalServerError, &models.Response{
		// 		Message: "Error uploading the cover image URL",
		// 		Status:  false,
		// 	})
		// }

		if err := ct.service.UploadPicture(int(userId), pictureURL, token); err != nil {
			return c.JSON(http.StatusInternalServerError,
				helper.ResponseFormat(http.StatusInternalServerError, helper.ServerGeneralError, nil))
		}

		return c.JSON(http.StatusOK,
			helper.ResponseFormat(http.StatusOK, "Upload Photo Success", nil))

	}
}
