package sub_account

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

type (
	GetTransferHistoryReqType struct {
		WalletAddress string `json:"wallet_address" validate:"required"`
	}
)

/*
GetTransferHistory Get transfer history
*/
func GetTransferHistory(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{
		"msg": "ok",
	})
}
