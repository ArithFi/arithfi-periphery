package routes

import (
	"github.com/arithfi/arithfi-periphery/cmd/server/handlers"
	"github.com/labstack/echo/v4"
)

func Routes(e *echo.Echo) {
	// Routes => handler
	e.GET("/", handlers.Hello)

	e.POST("/events", handlers.HandleEvents)
}
