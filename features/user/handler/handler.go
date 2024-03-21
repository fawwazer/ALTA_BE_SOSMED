package handler

import (
	"ALTA_BE_SOSMED/features/user"
	"ALTA_BE_SOSMED/helper"
	"context"
	"errors"

	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
	// "gorm.io/gorm"
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
				log.Print("error unsupport: ", err.Error())
				return c.JSON(http.StatusUnsupportedMediaType,
					helper.ResponseFormat(http.StatusUnsupportedMediaType, "format data tidak didukung", nil))
			}
			log.Print("error bad request input: ", err.Error())
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
		userID, err := strconv.ParseUint(c.Param("user_id"), 10, 32)
		if err != nil {
			log.Println("error param:", err.Error())
			return c.JSON(http.StatusBadRequest,
				helper.ResponseFormat(http.StatusBadRequest, helper.UserInputError, nil))
		}

		token, ok := c.Get("user").(*jwt.Token)
		if !ok {
			return c.JSON(http.StatusBadRequest,
				helper.ResponseFormat(http.StatusBadRequest, helper.UserInputError, nil))
		}

		result, err := ct.service.Profile(token, uint(userID))
		if err != nil {
			// Jika tidak ada data profil yang ditemukan, kembalikan respons "not found"
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return c.JSON(http.StatusNotFound,
					helper.ResponseFormat(http.StatusNotFound, "data tidak ditemukan", nil))
			}
			// Jika terjadi kesalahan lain selain "record not found",
			// kembalikan respons forbidden
			log.Println("error getting profile:", err.Error())
			return c.JSON(http.StatusForbidden,
				helper.ResponseFormat(http.StatusForbidden, "Anda tidak diizinkan mengakses profil pengguna lain", nil))
		}

		var response ProfileResponse
		response.UserID = result.UserID
		response.Nama = result.Nama
		response.Email = result.Email
		response.Gender = result.Gender
		response.Tgl_lahir = result.Tgl_lahir
		response.Picture = result.Picture

		return c.JSON(http.StatusOK,
			helper.ResponseFormat(http.StatusOK, "berhasil mendapatkan data", response))
	}
}

func (ct *controller) UpdateProfile() echo.HandlerFunc {
	return func(c echo.Context) error {
		userId, err := strconv.ParseInt(c.Param("user_id"), 10, 32)
		if err != nil {
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

		// var pictureURL string
		if file != nil { // Check if file nil
			cld, err := cloudinary.NewFromURL("cloudinary://426244812151882:GBqN6L8Rm77iHHkPXiemVPP_e2Y@dlosajdpy")
			if err != nil {
				log.Print("error connect error: ", err.Error())
				return err
			}
			resp, err := cld.Upload.Upload(context.Background(), file, uploader.UploadParams{})
			if err != nil {
				log.Print("error upload error: ", err.Error())
				return err
			}

			nama := c.FormValue("nama")
			email := c.FormValue("email")
			gender := c.FormValue("gender")
			tgl_lahir := c.FormValue("tgl_lahir")

			// convert string to bool
			convertGender, err := strconv.ParseBool(gender)
			if err != nil {
				log.Print("err convert: ", err)
			}

			var updateProcess user.User
			updateProcess.Nama = nama
			updateProcess.Email = email
			updateProcess.Gender = convertGender
			updateProcess.Tgl_lahir = tgl_lahir

			// log.Print(resp.SecureURL)
			// var updateProcess user.User
			if resp.SecureURL != "" { // Update picture only if URL is not empty
				updateProcess.Picture = resp.SecureURL
			}
			if err := ct.service.UpdateProfile(int(userId), token, updateProcess); err != nil {
				log.Println("error update account:", err.Error())
				if errors.Is(err, gorm.ErrRecordNotFound) {
					return c.JSON(http.StatusNotFound,
						helper.ResponseFormat(http.StatusNotFound, "data tidak ditemukan", nil))
				}
				// Jika terjadi kesalahan lain selain "record not found",
				// kembalikan respons forbidden
				log.Println("error update profile:", err.Error())
				return c.JSON(http.StatusForbidden,
					helper.ResponseFormat(http.StatusForbidden, "Anda tidak diizinkan mengakses profil pengguna lain", nil))
			}
			return c.JSON(http.StatusOK,
				helper.ResponseFormat(http.StatusOK, "Update Profile Success", nil))
		} else if file == nil {
			res, err := ct.service.GetPicture(token)
			if err != nil {
				return c.JSON(http.StatusInternalServerError,
					helper.ResponseFormat(http.StatusInternalServerError, helper.ServerGeneralError, nil))
			}
			// var pictureURL string
			// pictureURL := res.Picture

			nama := c.FormValue("nama")
			email := c.FormValue("email")
			gender := c.FormValue("gender")
			tgl_lahir := c.FormValue("tgl_lahir")

			// convert string to bool
			convertGender, err := strconv.ParseBool(gender)
			if err != nil {
				log.Print("err convert: ", err)
			}

			var updateProcess user.User
			updateProcess.Nama = nama
			updateProcess.Email = email
			updateProcess.Gender = convertGender
			updateProcess.Tgl_lahir = tgl_lahir
			updateProcess.Picture = res.Picture

			// if pictureURL != "" { // Update picture only if URL is not empty
			// 	var updateProcess user.User
			// 	updateProcess.Picture = pictureURL
			// }

			if err := ct.service.UpdateProfile(int(userId), token, updateProcess); err != nil {
				log.Println("error update account:", err.Error())
				if errors.Is(err, gorm.ErrRecordNotFound) {
					return c.JSON(http.StatusNotFound,
						helper.ResponseFormat(http.StatusNotFound, "data tidak ditemukan", nil))
				}
				// Jika terjadi kesalahan lain selain "record not found",
				// kembalikan respons forbidden
				log.Println("error update profile:", err.Error())
				return c.JSON(http.StatusForbidden,
					helper.ResponseFormat(http.StatusForbidden, "Anda tidak diizinkan mengakses profil pengguna lain", nil))
			}

			return c.JSON(http.StatusOK,
				helper.ResponseFormat(http.StatusOK, "Update Profile Success", nil))

		}
		return nil
	}
}

func (ct *controller) DeleteAccount() echo.HandlerFunc {
	return func(c echo.Context) error {
		userID, err := strconv.ParseUint(c.Param("user_id"), 10, 32)
		if err != nil {
			log.Println("error param:", err.Error())
			return c.JSON(http.StatusBadRequest,
				helper.ResponseFormat(http.StatusBadRequest, helper.UserInputError, nil))
		}

		token, ok := c.Get("user").(*jwt.Token)
		if !ok {
			log.Println("error token:", err.Error())
			return c.JSON(http.StatusBadRequest,
				helper.ResponseFormat(http.StatusBadRequest, helper.UserInputError, nil))
		}

		err = ct.service.DeleteAccount(uint(userID), token)
		if err != nil {
			log.Println("error deleting account:", err.Error())
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return c.JSON(http.StatusNotFound,
					helper.ResponseFormat(http.StatusNotFound, "data tidak ditemukan", nil))
			}
			// Jika terjadi kesalahan lain selain "record not found",
			// kembalikan respons forbidden
			log.Println("error Deleting profile:", err.Error())
			return c.JSON(http.StatusForbidden,
				helper.ResponseFormat(http.StatusForbidden, "Anda tidak diizinkan menghapus profil pengguna lain", nil))
		}

		return c.JSON(http.StatusCreated, helper.ResponseFormat(http.StatusCreated, "berhasil menghapus akun", nil))
	}
}
