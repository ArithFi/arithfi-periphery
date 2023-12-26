package routes

import (
	"github.com/arithfi/arithfi-periphery/internal/futures/account_trades"
	"github.com/arithfi/arithfi-periphery/internal/futures/market_data"
	"github.com/labstack/echo/v4"
)

func FuturesRoutes(e *echo.Echo) {
	// Market data
	e.GET("ping", market_data.Ping)
	e.GET("time", market_data.Time)
	//e.GET("exchangeInfo", market_data.ExchangeInfo)
	//e.GET("kline", market_data.Kline)
	//e.GET("ticker/24hr", market_data.Ticker24hr)
	//e.GET("ticker/price", market_data.TickerPrice)

	// Account trades
	e.GET("order", account_trades.QueryOrder)
	e.GET("orderAmendment", account_trades.GetOrderAmendment)
	e.GET("openOrder", account_trades.QueryCurrentOpenOrder)
	e.GET("openOrders", account_trades.QueryCurrentOpenOrders)
	e.GET("account", account_trades.GetAccount)
	e.GET("balance", account_trades.GetBalance)
	e.GET("userTrades", account_trades.GetTradeList)
	e.GET("allOrders", account_trades.GetAllOrders)
}
