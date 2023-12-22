package account_trades

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

type (
	QueryOrderReqType struct {
		OrderId int `json:"orderId" validate:"required"` // OrderId = PositionIndex
	}
)

/*
QueryOrder Check an order's status.
*/
func QueryOrder(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{
		"msg": "ok",
	})
}