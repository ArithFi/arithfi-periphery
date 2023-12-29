package routes

import (
	"github.com/arithfi/arithfi-periphery/internal/scan"
	"github.com/labstack/echo/v4"
)

func ScanRoutes(e *echo.Echo) {
	e.GET("scan/erc20_transfer_bsc", scan.ERC20TransferBSC)
	e.GET("scan/deposit_withdrawal", scan.DepositWithdrawal)
	e.GET("scan/f_future_trading", scan.FFutureTrading)
}
