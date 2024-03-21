package handler

import (
	"github.com/labstack/echo/v4"
)

type MediaResponse struct {
	StatusCode int       `json:"statusCode"`
	Message    string    `json:"message"`
	Data       *echo.Map `json:"data"`
}
