package handler

import (
	"ALTA_BE_SOSMED/features/user/data"
	"ALTA_BE_SOSMED/features/user/services"
	"net/http"

	"github.com/labstack/echo/v4"
)

func FileUpload(c echo.Context) error {
	//upload
	formHeader, err := c.FormFile("file")
	if err != nil {
		return c.JSON(
			http.StatusInternalServerError,
			MediaResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "error",
				Data:       &echo.Map{"data": "Select a file to upload"},
			})
	}

	//get file from header
	formFile, err := formHeader.Open()
	if err != nil {
		return c.JSON(
			http.StatusInternalServerError,
			MediaResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "error",
				Data:       &echo.Map{"data": err.Error()},
			})
	}

	uploadUrl, err := services.NewMediaUpload().FileUpload(data.File{File: formFile})
	if err != nil {
		return c.JSON(
			http.StatusInternalServerError,
			MediaResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "error",
				Data:       &echo.Map{"data": err.Error()},
			})
	}

	return c.JSON(
		http.StatusOK,
		MediaResponse{
			StatusCode: http.StatusOK,
			Message:    "success",
			Data:       &echo.Map{"data": uploadUrl},
		})
}

func RemoteUpload(c echo.Context) error {
	var url data.Url

	//validate the request body
	if err := c.Bind(&url); err != nil {
		return c.JSON(
			http.StatusBadRequest,
			MediaResponse{
				StatusCode: http.StatusBadRequest,
				Message:    "error",
				Data:       &echo.Map{"data": err.Error()},
			})
	}

	uploadUrl, err := services.NewMediaUpload().RemoteUpload(url)
	if err != nil {
		return c.JSON(
			http.StatusInternalServerError,
			MediaResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "error",
				Data:       &echo.Map{"data": "Error uploading file"},
			})
	}

	return c.JSON(
		http.StatusOK,
		MediaResponse{
			StatusCode: http.StatusOK,
			Message:    "success",
			Data:       &echo.Map{"data": uploadUrl},
		})
}
