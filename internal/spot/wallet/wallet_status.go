package wallet

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

type (
	GetWalletStatusReqType struct {
		WalletAddress string `json:"wallet_address" validate:"required"`
	}
)

/*
GetWalletStatus Get account status
*/
func GetWalletStatus(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]any{
		"success": true,
	})
}
