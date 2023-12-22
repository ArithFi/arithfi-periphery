package account_trades

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

type (
	QueryCurrentOrdersReqType struct {
		WalletAddress string `json:"walletAddress" validate:"required"` // WalletAddress
	}
)

/*
QueryCurrentOpenOrders Query current open orders
*/
func QueryCurrentOpenOrders(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{
		"msg": "ok",
	})
}
