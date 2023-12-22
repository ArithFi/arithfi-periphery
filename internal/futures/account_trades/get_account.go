package account_trades

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

type (
	GetAccountReqType struct {
		WalletAddress string `json:"walletAddress" validate:"required"` // WalletAddress
	}
)

/*
GetAccount Get account information
return assets and position
*/
func GetAccount(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{
		"msg": "ok",
	})
}
