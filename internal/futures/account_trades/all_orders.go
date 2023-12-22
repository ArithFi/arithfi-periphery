package account_trades

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

type (
	GetAllOrdersReqType struct {
		Symbol        string `json:"symbol" validate:"required"`        // Symbol
		WalletAddress string `json:"walletAddress" validate:"required"` // WalletAddress
		StartTime     int    `json:"startTime"`                         // StartTime
		EndTime       int    `json:"endTime"`                           // EndTime
		FromId        int    `json:"fromId"`                            // FromId
		Limit         int    `json:"limit"`                             // Limit
	}
)

/*
GetAllOrders Get all account orders
*/
func GetAllOrders(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{
		"msg": "ok",
	})
}
