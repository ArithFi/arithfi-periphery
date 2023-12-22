package sub_account

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

type (
	ModifySubAccountReqType struct {
		WalletAddress string `json:"wallet_address" validate:"required"`
	}
)

/*
ModifySubAccount Modify a sub_account
*/
func ModifySubAccount(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]any{
		"success": true,
	})
}
