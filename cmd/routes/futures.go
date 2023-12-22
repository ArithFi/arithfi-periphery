package routes

import (
	"github.com/arithfi/arithfi-periphery/internal/futures/account_trades"
	"github.com/arithfi/arithfi-periphery/internal/futures/market_data"
	"github.com/labstack/echo/v4"
)

func FuturesRoutes(e *echo.Echo) {
	// Market data
	e.GET("futures/ping", market_data.Ping)
	e.GET("futures/time", market_data.Time)
	//e.GET("futures/exchangeInfo", market_data.ExchangeInfo)
	//e.GET("futures/kline", market_data.Kline)
	//e.GET("futures/ticker/24hr", market_data.Ticker24hr)
	//e.GET("futures/ticker/price", market_data.TickerPrice)

	// Account trades
	e.GET("futures/order", account_trades.QueryOrder)
	e.POST("futures/order", account_trades.NewOrder)
	e.PUT("futures/order", account_trades.ModifyOrder)
	e.DELETE("futures/order", account_trades.CancelOrder)
	e.GET("futures/orderAmendment", account_trades.GetOrderAmendment)
	e.GET("futures/openOrder", account_trades.QueryCurrentOpenOrder)
	e.GET("futures/openOrders", account_trades.QueryCurrentOpenOrders)
}
