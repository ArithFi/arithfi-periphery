package account_trades

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

type (
	GetPositionMarginChangeHistoryReqType struct {
		WalletAddress string `json:"walletAddress" validate:"required"` // WalletAddress
	}
)

/*
GetPositionMarginChangeHistory Add position margin
*/
func GetPositionMarginChangeHistory(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{
		"msg": "ok",
	})
}