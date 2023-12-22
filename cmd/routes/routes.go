package routes

import (
	"github.com/arithfi/arithfi-periphery/internal/account_trades"
	"github.com/arithfi/arithfi-periphery/internal/market_data"
	"github.com/labstack/echo/v4"
)

func Routes(e *echo.Echo) {
	// General Info
	e.GET("/ping", market_data.Ping)
	e.GET("/time", market_data.Time)
	e.GET("/exchangeInfo", market_data.ExchangeInfo)
	e.GET("/kline", market_data.Kline)
	e.GET("/ticker/24hr", market_data.Ticker24hr)
	e.GET("/ticker/price", market_data.TickerPrice)

	e.POST("/events", account_trades.HandleEvents)
}
