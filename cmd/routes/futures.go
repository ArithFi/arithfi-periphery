package routes

import (
	"github.com/arithfi/arithfi-periphery/internal/futures/market_data"
	"github.com/labstack/echo/v4"
)

func FuturesRoutes(e *echo.Echo) {
	// General Info
	e.GET("futures/ping", market_data.Ping)
	e.GET("futures/time", market_data.Time)
	//e.GET("futures/exchangeInfo", market_data.ExchangeInfo)
	//e.GET("futures/kline", market_data.Kline)
	//e.GET("futures/ticker/24hr", market_data.Ticker24hr)
	//e.GET("futures/ticker/price", market_data.TickerPrice)
}
