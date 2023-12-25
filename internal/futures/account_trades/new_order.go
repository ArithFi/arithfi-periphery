package account_trades

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

type (
	NewOrderReqType struct {
		OrderId         int    `json:"orderId" validate:"required"` // OrderId = PositionIndex
		WalletAddress   string `json:"walletAddress" validate:"required"`
		Symbol          string `json:"symbol" validate:"required"` // Symbol = Pair = Product
		Side            string `json:"side" validate:"required"`   // Side, ['BUY', 'SELL']
		PositionSide    string `json:"positionSide"`               // PositionSide, ['LONG', 'SHORT']
		Type            string `json:"type" validate:"required"`   // Type, ['SELF', 'COPY']
		CopyFromAddress string `json:"copyFromAddress"`
		OpenPrice       int64  `json:"openPrice" validate:"required"`
		MarkPrice       int64  `json:"markPrice"`
		StopLossPrice   int64  `json:"stopLossPrice"`
		TakeProfitPrice int64  `json:"takeProfitPrice"`
		InitialMargin   int64  `json:"initialMargin" validate:"required"` // Initial Margin
		Leverage        int64  `json:"leverage" validate:"required"`      // Leverage
		Volume          int64  `json:"volume" validate:"required"`
		OpenFees        int64  `json:"openFees"`
	}
)

/*
NewOrder Send in a new order.
TODO
*/
func NewOrder(c echo.Context) error {
	// Obtain necessary order information and verify.

	// if Type = SELF
	// If there is an invitation relationship, it is necessary to associate the invite's information with the order for
	// the convenience of generating commission records.

	// if Type = COPY
	// Need to associate and copy the information of the person for easy generation of profit sharing records.

	// Record the circulation of funds.
	// Transaction Fee Circulation -> Contract Account (Instant Settlement)
	// Remaining Opening Amount -> Contract Account (Instantly Available)
	// Contract Account -> Commission Account (only generates records, not immediately credited, manual trigger required)

	var req NewOrderReqType
	if err := c.Bind(&req); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, map[string]string{
		"msg": "ok",
	})
}
