package sub_account

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

type (
	GetUserAssetsReqType struct {
		WalletAddress string `json:"wallet_address" validate:"required"`
	}
)

/*
GetUserAssets Get user assets
*/
func GetUserAssets(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]any{
		"success": true,
	})
}
