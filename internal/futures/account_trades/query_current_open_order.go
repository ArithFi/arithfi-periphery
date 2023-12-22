package account_trades

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

type (
	QueryCurrentOrderReqType struct {
		OrderId int `json:"orderId" validate:"required"` // OrderId = PositionIndex
	}
)

/*
QueryCurrentOpenOrder Query current open order
*/
func QueryCurrentOpenOrder(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{
		"msg": "ok",
	})
}
