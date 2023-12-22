package sub_account

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

type (
	GetSubAccountStatusReqType struct {
		WalletAddress string `json:"wallet_address" validate:"required"`
	}
)

/*
GetSubAccountStatus Get sub_account status
*/
func GetSubAccountStatus(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{
		"msg": "ok",
	})
}
