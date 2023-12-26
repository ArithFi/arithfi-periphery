package routes

import (
	"github.com/arithfi/arithfi-periphery/internal/spot/sub_account"
	"github.com/arithfi/arithfi-periphery/internal/spot/wallet"
	"github.com/labstack/echo/v4"
)

func SpotRoutes(e *echo.Echo) {
	// Sub Account, off chain
	e.GET("api/sub_account", sub_account.QuerySubAccount)
	e.GET("api/sub_account/status", sub_account.GetSubAccountStatus)
	e.GET("api/sub_account/account_snapshot", sub_account.DailyAccountSnapshot)
	e.GET("api/sub_account/deposit/history", sub_account.GetDepositHistory)
	e.GET("api/sub_account/withdraw/history", sub_account.GetWithdrawHistory)
	e.GET("api/sub_account/assets", sub_account.GetUserAssets)
	e.GET("api/sub_account/transfer/history", sub_account.GetTransferHistory)

	// Wallet, on chain
	e.GET("api/system/status", wallet.GetSystemStatus)
	e.GET("api/wallet/assets", wallet.GetWalletAssetOnChain)
	e.GET("api/wallet/status", wallet.GetWalletStatus)
}
