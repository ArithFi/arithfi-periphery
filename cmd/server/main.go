package main

import (
	"github.com/arithfi/arithfi-periphery/cmd/server/handlers"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	// Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Routes => handler
	e.GET("/", handlers.Hello)

	e.POST("/events", handlers.HandleEvents)

	// Start server
	e.Logger.Fatal(e.Start(":8080"))
}
