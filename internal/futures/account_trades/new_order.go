package account_trades

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

type (
	Request struct {
		OrderId         int    `json:"orderId" validate:"required"` // OrderId = PositionIndex
		WalletAddress   string `json:"walletAddress" validate:"required"`
		Symbol          string `json:"symbol" validate:"required"` // Symbol = Pair = Product
		PositionSide    string `json:"positionSide"`               // PositionSide, ['LONG', 'SHORT']
		Type            string `json:"type" validate:"required"`   // Type, ['SELF', 'COPY']
		CopyFromAddress string `json:"copyFromAddress"`
		OpenPrice       int64  `json:"openPrice" validate:"required"`
		MarkPrice       int64  `json:"markPrice"`
		StopLossPrice   int64  `json:"stopLossPrice"`
		TakeProfitPrice int64  `json:"takeProfitPrice"`
		InitialMargin   int64  `json:"initialMargin" validate:"required"` // Initial Margin
		Volume          int64  `json:"volume" validate:"required"`
		OpenFees        int64  `json:"openFees"`
	}
)

/*
NewOrder Send in a new order.
*/
func NewOrder(c echo.Context) error {
	// 获取必要的订单信息，并校验

	// 如果是自开单
	// 如果有邀请关系，需要关联邀请人信息到该订单，便于返佣记录生成；

	// 如果是复制单
	// 需要关联复制人信息，便于分润记录生成

	// 记录资金的流转
	// 手续费流转 -》 合约账户（立即到账）
	// 剩余开仓金额 -》合约账户（立即到账）
	// 合约账户 -》 返佣人账户（只生成记录，不立即到账，需要手动触发）
	return c.JSON(http.StatusOK, map[string]string{
		"msg": "ok",
	})
}
