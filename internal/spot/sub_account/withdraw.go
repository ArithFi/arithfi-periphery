package sub_account

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

type (
	WithdrawReqType struct {
		WalletAddress string `json:"wallet_address" validate:"required"`
	}
)

/*
Withdraw Withdraw
*/
func Withdraw(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]any{
		"success": true,
	})
}
