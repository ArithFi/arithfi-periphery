package routes

import (
	"github.com/arithfi/arithfi-periphery/internal/futures/account_trades"
	"github.com/arithfi/arithfi-periphery/internal/futures/market_data"
	"github.com/labstack/echo/v4"
)

func FuturesRoutes(e *echo.Echo) {
	// Market data
	e.GET("fapi/ping", market_data.Ping)
	e.GET("fapi/time", market_data.Time)
	//e.GET("exchangeInfo", market_data.ExchangeInfo)
	//e.GET("kline", market_data.Kline)
	//e.GET("ticker/24hr", market_data.Ticker24hr)
	//e.GET("ticker/price", market_data.TickerPrice)

	// Account trades
	e.GET("fapi/order", account_trades.QueryOrder)
	e.GET("fapi/orderAmendment", account_trades.GetOrderAmendment)
	e.GET("fapi/openOrder", account_trades.QueryCurrentOpenOrder)
	e.GET("fapi/openOrders", account_trades.QueryCurrentOpenOrders)
	e.GET("fapi/account", account_trades.GetAccount)
	e.GET("fapi/balance", account_trades.GetBalance)
	e.GET("fapi/userTrades", account_trades.GetTradeList)
	e.GET("fapi/allOrders", account_trades.GetAllOrders)
}
