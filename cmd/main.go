package main

import (
	"github.com/arithfi/arithfi-periphery/cmd/routes"
	_ "github.com/arithfi/arithfi-periphery/configs"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	// Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Routes
	routes.FuturesRoutes(e)

	// Start server
	e.Logger.Fatal(e.Start(":8080"))
}
