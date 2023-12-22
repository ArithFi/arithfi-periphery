package account_trades

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

type (
	CancelOrderReqType struct {
		OrderId int `json:"orderId" validate:"required"` // OrderId = PositionIndex
	}
)

/*
CancelOrder Cancel an order.
TODO
*/
func CancelOrder(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{
		"msg": "ok",
	})
}
