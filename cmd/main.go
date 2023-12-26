package main

import (
	"github.com/arithfi/arithfi-periphery/cmd/routes"
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
	// Hosts
	hosts := map[string]*Host{}

	futures := echo.New()

	// Middleware
	futures.Use(middleware.Logger())
	futures.Use(middleware.Recover())
	futures.Use(middleware.Gzip())
	futures.Use(middleware.CORS())

	// Routes
	routes.FuturesRoutes(futures)

	hosts["fapi.localhost:8080"] = &Host{futures}

	spot := echo.New()

	// Middleware
	spot.Use(middleware.Logger())
	spot.Use(middleware.Recover())
	spot.Use(middleware.Gzip())
	spot.Use(middleware.CORS())

	routes.SpotRoutes(spot)

	hosts["api.localhost:8080"] = &Host{spot}

	e := echo.New()
	e.Any("/*", func(c echo.Context) (err error) {
		req := c.Request()
		res := c.Response()
		host := hosts[req.Host]

		if host == nil {
			err = echo.ErrNotFound
		} else {
			host.Echo.ServeHTTP(res, req)
		}

		return
	})

	e.Logger.Fatal(e.Start("localhost:8080"))
}
