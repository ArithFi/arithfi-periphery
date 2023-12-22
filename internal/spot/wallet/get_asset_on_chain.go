package wallet

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

type (
	GetAssetOnChainReqType struct {
		WalletAddress string `json:"wallet_address" validate:"required"`
	}
)

/*
GetWalletAssetOnChain Get asset on chain
*/
func GetWalletAssetOnChain(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]any{
		"success": true,
	})
}
