package sub_account

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

type (
	DepositReqType struct {
		WalletAddress string `json:"wallet_address" validate:"required"`
	}
)

/*
Deposit Deposit
TODO
*/
func Deposit(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]any{
		"success": true,
	})
}
