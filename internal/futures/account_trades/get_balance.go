package account_trades

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

type (
	GetBalanceReqType struct {
		WalletAddress string `json:"walletAddress" validate:"required"` // WalletAddress
	}
)

/*
GetBalance Get account balance
*/
func GetBalance(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{
		"msg": "ok",
	})
}
