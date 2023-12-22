package sub_account

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

type (
	GetWithdrawHistoryReqType struct {
		WalletAddress string `json:"wallet_address" validate:"required"`
	}
)

/*
GetWithdrawHistory Get withdraw history
*/
func GetWithdrawHistory(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]any{
		"success": true,
	})
}
