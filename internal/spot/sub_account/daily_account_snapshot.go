package sub_account

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

type (
	DailyAccountSnapshotReqType struct {
		WalletAddress string `json:"wallet_address" validate:"required"`
	}
)

/*
DailyAccountSnapshot Get daily account snapshot
*/
func DailyAccountSnapshot(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{
		"msg": "ok",
	})
}
