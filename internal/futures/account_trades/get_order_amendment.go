package account_trades

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

type (
	GetOrderAmendmentReqType struct {
		OrderId int `json:"orderId" validate:"required"` // OrderId = PositionIndex
	}
)

/*
GetOrderAmendment Get order modification history
*/
func GetOrderAmendment(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{
		"msg": "ok",
	})
}
