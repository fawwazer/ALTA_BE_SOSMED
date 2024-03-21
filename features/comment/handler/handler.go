package handler

import (
	"ALTA_BE_SOSMED/features/comment"
	"ALTA_BE_SOSMED/helper"
	"ALTA_BE_SOSMED/middlewares"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type controller struct {
	s comment.ComService
}

func NewHandler(service comment.ComService) comment.ComController {
	return &controller{
		s: service,
	}
}

func (ct *controller) View() echo.HandlerFunc {
	return func(c echo.Context) error {
		token, ok := c.Get("user").(*jwt.Token)
		if !ok {
			return c.JSON(http.StatusBadRequest,
				helper.ResponseFormat(http.StatusBadRequest, helper.UserInputError, nil))
		}
		result, err := ct.s.View(token)
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

func (ct *controller) Add() echo.HandlerFunc {
	return func(c echo.Context) error {
		var input ComRequest
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

		var inputProcess comment.Comment
		inputProcess.Comment = input.Comment
		inputProcess.Pemiliks = input.Pemiliks
		_, err = ct.s.AddCom(token, inputProcess)
		if err != nil {
			log.Println("error insert db:", err.Error())
			return c.JSON(http.StatusInternalServerError, helper.ResponseFormat(http.StatusInternalServerError, helper.ServerGeneralError, nil))
		}

		return c.JSON(http.StatusCreated, helper.ResponseFormat(http.StatusCreated, "berhasil menambahkan data", nil))
	}
}

func (ct *controller) Update() echo.HandlerFunc {
	return func(c echo.Context) error {
		var input ComRequest
		readID := c.Param("id")
		cnv, err := strconv.Atoi(readID)
		if err != nil {
			return c.JSON(http.StatusBadRequest,
				helper.ResponseFormat(http.StatusBadRequest, "data yang dikirmkan tidak sesuai", nil))
		}
		err = c.Bind(&input)
		if err != nil {
			log.Println("error bind data:", err.Error())
			if strings.Contains(err.Error(), "unsupport") {
				return c.JSON(http.StatusUnsupportedMediaType,
					helper.ResponseFormat(http.StatusUnsupportedMediaType, "format data tidak didukung", nil))
			}
			return c.JSON(http.StatusBadRequest,
				helper.ResponseFormat(http.StatusBadRequest, "data yang dikirmkan tidak sesuai", nil))
		}

		hp := middlewares.DecodeToken(c.Get("user").(*jwt.Token))

		if hp == "" {
			log.Println("error decode token:", "hp tidak ditemukan")
			return c.JSON(http.StatusUnauthorized, helper.ResponseFormat(http.StatusUnauthorized, "tidak dapat mengakses fitur ini", nil))
		}

		var inputProcess comment.Comment
		inputProcess.Comment = input.Comment
		inputProcess.Pemiliks = input.Pemiliks

		err = ct.s.Update(hp, uint(cnv), inputProcess)

		if err != nil {
			log.Println("error update db:", err.Error())
			return c.JSON(http.StatusInternalServerError, helper.ResponseFormat(http.StatusInternalServerError, "terjadi kesalahan pada proses server", nil))
		}

		return c.JSON(http.StatusOK, helper.ResponseFormat(http.StatusOK, "berhasil mengubah data", nil))
	}

}

func (del *controller) Delete() echo.HandlerFunc {
	return func(c echo.Context) error {
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			return err
		}
		converID := comment.Comment{Model: gorm.Model{ID: uint(id)}}
		err = del.s.DeleteCom(converID)
		if err != nil {
			var code = http.StatusInternalServerError
			if strings.Contains(err.Error(), "validation") {
				code = http.StatusBadRequest
			}
			return c.JSON(code,
				helper.ResponseFormat(code, err.Error(), nil))
		}
		return c.JSON(http.StatusOK,
			helper.ResponseFormat(http.StatusOK, "berhasil menghapus data", nil))
	}
}
