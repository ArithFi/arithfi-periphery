package sub_account

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

type (
	NewSubAccountReqType struct {
		WalletAddress string `json:"wallet_address" validate:"required"`
	}
)

/*
NewSubAccount Create a new sub_account
TODO
*/
func NewSubAccount(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{
		"msg": "ok",
	})
}
