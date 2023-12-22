package account_trades

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

type (
	ModifyOrderReqType struct {
		OrderId         int    `json:"orderId" validate:"required"` // OrderId = PositionIndex
		Side            string `json:"side" validate:"required"`    // Side , ['BUY', 'SELL']
		StopLossPrice   int64  `json:"stopLossPrice"`
		TakeProfitPrice int64  `json:"takeProfitPrice"`
		Volume          int64  `json:"volume" validate:"required"` // Updated Volume
	}
)

/*
ModifyOrder Modify an order.
Include: Add position、Close position、Modify stop loss、Modify take profit
TODO
*/
func ModifyOrder(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{
		"msg": "ok",
	})
}
