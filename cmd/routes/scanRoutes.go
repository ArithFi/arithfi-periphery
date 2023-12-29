package routes

import (
	"github.com/arithfi/arithfi-periphery/internal/scan"
	"github.com/labstack/echo/v4"
)

func ScanRoutes(e *echo.Echo) {
	// scan erc20_transfer_bsc
	e.GET("scan/erc20_transfer_bsc", scan.ERC20TransferBSC)
}
