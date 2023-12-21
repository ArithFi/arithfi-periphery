package main

import (
	"github.com/arithfi/arithfi-periphery/cmd/server/configs"
	"github.com/arithfi/arithfi-periphery/cmd/server/routes"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	// Echo instance
	e := echo.New()

	//run database
	configs.ConnectDB()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Routes
	routes.Routes(e)

	// Start server
	e.Logger.Fatal(e.Start(":8080"))
}
