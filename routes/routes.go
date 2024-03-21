package routes

import (
	"ALTA_BE_SOSMED/config"
	"ALTA_BE_SOSMED/features/comment"
	"ALTA_BE_SOSMED/features/user"

	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
)

func InitRoute(c *echo.Echo, ct1 user.UserController, ct comment.ComController) {
	userRoute(c, ct1)
	comRoute(c, ct)
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

func comRoute(c *echo.Echo, ct comment.ComController) {
	c.GET("/comment", ct.View(), echojwt.WithConfig(echojwt.Config{
		SigningKey: []byte(config.JWTSECRET),
	}))
	c.POST("/comment", ct.Add(), echojwt.WithConfig(echojwt.Config{
		SigningKey: []byte(config.JWTSECRET),
	}))
	c.PUT("/comment/:id", ct.Update(), echojwt.WithConfig(echojwt.Config{
		SigningKey: []byte(config.JWTSECRET),
	}))
	c.DELETE("/comment/:id", ct.Delete(), echojwt.WithConfig(echojwt.Config{
		SigningKey: []byte(config.JWTSECRET),
	}))
}
