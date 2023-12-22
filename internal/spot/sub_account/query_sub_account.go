package sub_account

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

type (
	QuerySubAccountReqType struct {
		WalletAddress string `json:"wallet_address" validate:"required"`
	}
)

/*
QuerySubAccount Query sub_account information
*/
func QuerySubAccount(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{
		"msg": "ok",
	})
}
