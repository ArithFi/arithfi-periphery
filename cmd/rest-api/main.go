package main

import (
	_ "github.com/arithfi/arithfi-periphery/configs"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type (
	Host struct {
		Echo *echo.Echo
	}
)

func main() {
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.Gzip())
	e.Use(middleware.CORS())

	e.Logger.Fatal(e.Start(":8080"))
}
