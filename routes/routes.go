package routes

import (
	"ALTA_BE_SOSMED/config"
	"ALTA_BE_SOSMED/features/user"

	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	// "github.com/labstack/echo"
)

func InitRoute(c *echo.Echo, ct1 user.UserController) {
	userRoute(c, ct1)
}

func userRoute(c *echo.Echo, ct1 user.UserController) {
	c.POST("/users", ct1.Add())
	c.POST("/login", ct1.Login())
	c.GET("/profile", ct1.Profile(), echojwt.WithConfig(echojwt.Config{
		SigningKey: []byte(config.JWTSECRET),
	}))
	c.POST("/profile/:user_id/upload", ct1.UploadPicture(), echojwt.WithConfig(echojwt.Config{
		SigningKey: []byte(config.JWTSECRET),
	}))
}
