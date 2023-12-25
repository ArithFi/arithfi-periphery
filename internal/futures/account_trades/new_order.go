package account_trades

import (
	"context"
	"github.com/arithfi/arithfi-periphery/configs"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"log"
	"net/http"
)

type (
	NewOrderReqType struct {
		OrderId         int     `json:"orderId" validate:"required"` // OrderId = PositionIndex
		WalletAddress   string  `json:"walletAddress" validate:"required"`
		Symbol          string  `json:"symbol" validate:"required"` // Symbol = Pair = Product
		Side            string  `json:"side" validate:"required"`   // Side, ['BUY', 'SELL']
		PositionSide    string  `json:"positionSide"`               // PositionSide, ['LONG', 'SHORT']
		Type            string  `json:"type" validate:"required"`   // Type, ['SELF', 'COPY']
		CopyFromAddress string  `json:"copyFromAddress"`
		OpenPrice       float64 `json:"openPrice" validate:"required"`
		MarkPrice       float64 `json:"markPrice"`
		StopLossPrice   float64 `json:"stopLossPrice"`
		TakeProfitPrice float64 `json:"takeProfitPrice"`
		InitialMargin   float64 `json:"initialMargin" validate:"required"` // Initial Margin
		Leverage        int64   `json:"leverage" validate:"required"`      // Leverage
		Volume          float64 `json:"volume" validate:"required"`
		OpenFees        float64 `json:"openFees"`
	}
)

/*
NewOrder Send in a new order.
TODO
*/
func NewOrder(c echo.Context) error {
	var req NewOrderReqType
	if err := c.Bind(&req); err != nil {
		return err
	}

	// if Type = SELF
	// If there is an invitation relationship, it is necessary to associate the invite's information with the order for
	// the convenience of generating commission records.
	if req.Type == "SELF" {
		// create a order item in mongodb, and record relationship info in this order,
		// prepare for commission calculation
		futuresOrderCollection := configs.GetCollection("futures_order")

		one, err := futuresOrderCollection.InsertOne(context.TODO(), bson.D{
			{"_id", req.OrderId},
			{"wallet_address", req.WalletAddress},
			{"symbol", req.Symbol},
			{"side", req.Side},
			{"position_side", req.PositionSide},
			{"type", req.Type},
			{"copy_from_address", req.CopyFromAddress},
			{"open_price", req.OpenPrice},
			{"mark_price", req.MarkPrice},
			{"stop_loss_price", req.StopLossPrice},
			{"take_profit_price", req.TakeProfitPrice},
			{"initial_margin", req.InitialMargin},
			{"leverage", req.Leverage},
			{"volume", req.Volume},
			{"open_fees", req.OpenFees},
		})

		log.Println("Inserted ID", one.InsertedID)

		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{
				"msg": err.Error(),
			})
		}
		// need to record the circulation of funds.
		//
	}

	// if Type = COPY
	// Need to associate and copy the information of the person for easy generation of profit sharing records.
	if req.Type == "COPY" {
		// create a order item in mongodb, and record copy relationship info in this order,
		// prepare for commission calculation

	}

	// Record the circulation of funds.
	// Transaction Fee Circulation -> Contract Account (Instant Settlement)
	// Remaining Opening Amount -> Contract Account (Instantly Available)
	// Contract Account -> Commission Account (only generates records, not immediately credited, manual trigger required)

	return c.JSON(http.StatusOK, map[string]string{
		"msg": "ok",
	})
}
