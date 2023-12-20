// handlers/handlers.go

package handlers

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

// Hello function to return Hello, World!
func Hello(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, World!")
}

type Event struct {
	CreateTime           string `json:"create_time" form:"create_time" query:"create_time"`
	UpdateTime           string `json:"update_time" form:"update_time" query:"update_time"`
	TimeStamp            string `json:"time_stamp" form:"time_stamp" query:"time_stamp"`
	Product              string `json:"product" form:"product" query:"product"`
	PositionIndex        int    `json:"position_index" form:"position_index" query:"position_index"`
	Currency             string `json:"currency" form:"currency" query:"currency"`
	OrderType            string `json:"order_type" form:"order_type" query:"order_type"`
	Mode                 string `json:"mode" form:"mode" query:"mode"`
	Direction            string `json:"direction" form:"direction" query:"direction"`
	Margin               string `json:"margin" form:"margin" query:"margin"`
	Volume               string `json:"volume" form:"volume" query:"volume"`
	StopLossPrice        string `json:"stop_loss_price" form:"stop_loss_price" query:"stop_loss_price"`
	TakeProfitPrice      string `json:"take_profit_price" form:"take_profit_price" query:"take_profit_price"`
	Fees                 string `json:"fees" form:"fees" query:"fees"`
	ExecutionFees        string `json:"execution_fees" form:"execution_fees" query:"execution_fees"`
	SellValue            string `json:"sell_value" form:"sell_value" query:"sell_value"`
	WalletAddress        string `json:"wallet_address" form:"wallet_address" query:"wallet_address"`
	KolAddress           string `json:"kol_address" form:"kol_address" query:"kol_address"`
	Status               string `json:"status" form:"status" query:"status"`
	ClearStatus          string `json:"clear_status" form:"clear_status" query:"clear_status"`
	Leverage             string `json:"leverage" form:"leverage" query:"leverage"`
	LiquidationVolume    string `json:"liquidation_volume" form:"liquidation_volume" query:"liquidation_volume"`
	AvailableBalance     string `json:"available_balance" form:"available_balance" query:"available_balance"`
	CopyAccountBalance   string `json:"copy_account_balance" form:"copy_account_balance" query:"copy_account_balance"`
	Profit               string `json:"profit" form:"profit" query:"profit"`
	CopyProfitCommission string `json:"copy_profit_commission" form:"copy_profit_commission" query:"copy_profit_commission"`
}

// HandleEvents function to handle events
func HandleEvents(c echo.Context) error {
	e := new(Event)
	if err := c.Bind(e); err != nil {
		return c.String(http.StatusBadRequest, "bad request")
	}

	if e.OrderType == "MARKET_ORDER_FEE" {
		return c.JSON(http.StatusOK, e)
	}

	if e.OrderType == "LIMIT_ORDER_FEE" {
		return c.JSON(http.StatusOK, e)
	}

	if e.OrderType == "SL_ORDER_FEE" {
		return c.JSON(http.StatusOK, e)
	}

	if e.OrderType == "MARKET_CLOSE_FEE" {
		return c.JSON(http.StatusOK, e)
	}

	if e.OrderType == "LIMIT_REQUEST" {
		return c.JSON(http.StatusOK, e)
	}

	if e.OrderType == "LIMIT_CANCEL" {
		return c.JSON(http.StatusOK, e)
	}

	if e.OrderType == "LIMIT_EDIT" {
		return c.JSON(http.StatusOK, e)
	}

	if e.OrderType == "MARKET_LIQUIDATION" {
		return c.JSON(http.StatusOK, e)
	}

	if e.OrderType == "MARKET_ORDER_ADD" {
		return c.JSON(http.StatusOK, e)
	}

	if e.OrderType == "SL_ORDER_FEE" {
		return c.JSON(http.StatusOK, e)
	}

	if e.OrderType == "TPSL_EDIT" {
		return c.JSON(http.StatusOK, e)
	}

	if e.OrderType == "TP_ORDER_FEE" {
		return c.JSON(http.StatusOK, e)
	}

	return c.NoContent(http.StatusBadRequest)
}
