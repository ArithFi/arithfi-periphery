package routes

import (
	"github.com/arithfi/arithfi-periphery/internal/futures/market_data"
	"github.com/arithfi/arithfi-periphery/internal/futures/summary"
	"github.com/labstack/echo/v4"
)

func FuturesRoutes(e *echo.Echo) {
	// Market data
	e.GET("fapi/ping", market_data.Ping)
	e.GET("fapi/time", market_data.Time)

	e.POST("fapi/update_daily_deposit_withdraw_metrics", summary.UpdateDailyDepositWithdrawMetrics)
	e.POST("fapi/update_daily_destroy_metrics", summary.UpdateDailyFuturesDestroyMetrics)
	e.POST("fapi/update_daily_futures_metrics", summary.UpdateDailyFuturesMetrics)
}
