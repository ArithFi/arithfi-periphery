package account_trades

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

type (
	GetTradeListReqType struct {
		Symbol        string `json:"symbol" validate:"required"`        // Symbol
		WalletAddress string `json:"walletAddress" validate:"required"` // WalletAddress
		StartTime     int    `json:"startTime"`                         // StartTime
		EndTime       int    `json:"endTime"`                           // EndTime
		FromId        int    `json:"fromId"`                            // FromId
		Limit         int    `json:"limit"`                             // Limit
	}
)

/*
GetTradeList Get trades for a specific account and symbol.
*/
func GetTradeList(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{
		"msg": "ok",
	})
}
