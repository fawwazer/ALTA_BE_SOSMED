package routes

import (
	"ALTA_BE_SOSMED/config"
	post "ALTA_BE_SOSMED/features/post"
	user "ALTA_BE_SOSMED/features/user"

	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	// "github.com/labstack/echo"
)

func InitRoute(c *echo.Echo, ct1 user.UserController, pc post.PostController) {
	userRoute(c, ct1)
	postRoute(c, pc)
}

func userRoute(c *echo.Echo, ct1 user.UserController) {
	c.POST("/users", ct1.Add())
	c.POST("/login", ct1.Login())
	c.GET("/profile", ct1.Profile(), echojwt.WithConfig(echojwt.Config{
		SigningKey: []byte(config.JWTSECRET),
	}))
	c.PUT("/profile/:user_id", ct1.UpdateProfile(), echojwt.WithConfig(echojwt.Config{
		SigningKey: []byte(config.JWTSECRET),
	}))
	c.DELETE("/users/:user_id", ct1.DeleteAccount(), echojwt.WithConfig(echojwt.Config{
		SigningKey: []byte(config.JWTSECRET),
	}))
}

func postRoute(c *echo.Echo, pc post.PostController) {
	c.POST("/posting", pc.Add(), echojwt.WithConfig(echojwt.Config{
		SigningKey: []byte(config.JWTSECRET),
	}))
}
