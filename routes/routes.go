package routes

import (
	"ALTA_BE_SOSMED/config"
	"ALTA_BE_SOSMED/features/user"
	"ALTA_BE_SOSMED/features/user/handler"

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
	c.GET("/profile/:user_id", ct1.Profile(), echojwt.WithConfig(echojwt.Config{
		SigningKey: []byte(config.JWTSECRET),
	}))
	c.PUT("/profile/:user_id", ct1.UpdateProfile(), echojwt.WithConfig(echojwt.Config{
		SigningKey: []byte(config.JWTSECRET),
	}))
	c.DELETE("/users/:user_id", ct1.DeleteAccount(), echojwt.WithConfig(echojwt.Config{
		SigningKey: []byte(config.JWTSECRET),
	}))

	c.POST("/file", handler.FileUpload)
	c.POST("/remote", handler.RemoteUpload)
}
