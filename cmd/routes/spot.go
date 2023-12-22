package routes

import (
	"github.com/arithfi/arithfi-periphery/internal/spot/sub_account"
	"github.com/arithfi/arithfi-periphery/internal/spot/wallet"
	"github.com/labstack/echo/v4"
)

func SpotRoutes(e *echo.Echo) {
	// Sub Account, off chain
	e.GET("sub_account", sub_account.QuerySubAccount)
	e.POST("sub_account", sub_account.NewSubAccount)
	e.DELETE("sub_account", sub_account.DeleteSubAccount)
	e.PUT("sub_account", sub_account.ModifySubAccount)
	e.GET("sub_account/status", sub_account.GetSubAccountStatus)
	e.GET("sub_account/account_snapshot", sub_account.DailyAccountSnapshot)
	e.POST("sub_account/deposit", sub_account.Deposit)
	e.GET("sub_account/deposit/history", sub_account.GetDepositHistory)
	e.POST("sub_account/withdraw", sub_account.Withdraw)
	e.GET("sub_account/withdraw/history", sub_account.GetWithdrawHistory)
	e.GET("sub_account/assets", sub_account.GetUserAssets)
	e.GET("sub_account/transfer/history", sub_account.GetTransferHistory)

	// Wallet, on chain
	e.GET("system/status", wallet.GetSystemStatus)
	e.GET("wallet/assets", wallet.GetWalletAssetOnChain)
	e.GET("wallet/status", wallet.GetWalletStatus)
}
