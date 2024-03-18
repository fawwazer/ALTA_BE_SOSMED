package main

import (
	"ALTA_BE_SOSMED/config"
	"ALTA_BE_SOSMED/features/user/data"
	"ALTA_BE_SOSMED/features/user/handler"
	"ALTA_BE_SOSMED/features/user/services"
	"ALTA_BE_SOSMED/routes"

	// "github.com/labstack/echo"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	e := echo.New()
	cfg := config.InitConfig()
	db := config.InitSQL(cfg)

	userData := data.New(db)
	userService := services.NewService(userData)
	userHandler := handler.NewUserHandler(userService)

	e.Pre(middleware.RemoveTrailingSlash())
	e.Use(middleware.Logger())
	e.Use(middleware.CORS())
	routes.InitRoute(e, userHandler)
	e.Logger.Fatal(e.Start(":8080"))
}
