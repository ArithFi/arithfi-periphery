package sub_account

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

type (
	DeleteSubAccountReqType struct {
		WalletAddress string `json:"wallet_address" validate:"required"`
	}
)

/*
DeleteSubAccount Delete a sub_account
TODO
*/
func DeleteSubAccount(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]any{
		"success": true,
	})
}
