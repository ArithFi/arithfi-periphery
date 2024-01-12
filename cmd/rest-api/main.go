package main

import (
	_ "github.com/arithfi/arithfi-periphery/configs"
	"github.com/arithfi/arithfi-periphery/internal/datafeed"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"net/http"
	"strconv"
	"time"
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

	// Routes
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello World!")
	})
	e.GET("/time", func(c echo.Context) error {
		return c.String(http.StatusOK, strconv.FormatInt(time.Now().Unix(), 10))
	})
	e.GET("/config", datafeed.GetConfig)
	e.GET("/symbols", datafeed.Symbols)
	e.GET("/symbol_info", datafeed.SymbolInfo)
	e.GET("/search", datafeed.Search)
	e.GET("/history", datafeed.History)

	e.Logger.Fatal(e.Start(":8080"))
}
