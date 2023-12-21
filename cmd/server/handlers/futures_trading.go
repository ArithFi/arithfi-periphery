package handlers

import (
	"github.com/arithfi/arithfi-periphery/cmd/server/configs"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
)

var futures_order_collection *mongo.Collection = configs.GetCollection("futures_order")
var futures_tx_collection *mongo.Collection = configs.GetCollection("futures_tx")
var futures_tokentxn_collection *mongo.Collection = configs.GetCollection("futures_tokentxn")
var futures_user_collection *mongo.Collection = configs.GetCollection("user")

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
		// 市价开单成功，需要创建仓位，记录 tx 动作，
		// 记录资金流转; 个人可用(手续费、剩余金额) -> 期货系统合约可用;
		// 如果存在邀请关系的话，期货系统合约 -> 邀请人可用余额

		return c.JSON(http.StatusOK, e)
	}

	if e.OrderType == "LIMIT_ORDER_FEE" {
		// 限价开单成功，需要创建仓位，记录tx动作，
		// 记录资金变化; 从 个人冻结账户(手续费、剩余金额) -> 期货系统合约可用余额
		//
		// 如果存在邀请关系的话，期货系统合约 -> 邀请人可用余额 (不立即到账)
		// 如果是Copy单，存在一个 return 剩余金额（返回），不会立即到账， 期货系统合约 -> 个人可用余额
		return c.JSON(http.StatusOK, e)
	}

	if e.OrderType == "TP_ORDER_FEE" {
		// TP 关单成功，删掉用户的活跃订单
		// 期货系统合约 -> 个人可用余额
		// 执行费个人可用余额 -> 期货系统合约可用余额
		// 如果 Copy 单，分润 期货系统合约 -> 到邀请人 (不立即到账)
		// 如果 Copy 单，存在 return 剩余金额，总额不会立即到账， 期货系统合约 -> 个人可用余额
		return c.JSON(http.StatusOK, e)
	}

	if e.OrderType == "SL_ORDER_FEE" {
		// TP 关单成功，删掉用户的活跃订单
		// 期货系统合约 -> 个人可用余额
		// 执行费个人可用余额 -> 期货系统合约可用余额
		// 如果 Copy 单，分润 期货系统合约 -> 到邀请人 (不立即到账)
		// 如果 Copy 单，存在 return 剩余金额，总额不会立即到账， 期货系统合约 -> 个人可用余额
		return c.JSON(http.StatusOK, e)
	}

	if e.OrderType == "MARKET_CLOSE_FEE" {
		// TP 关单成功，删掉用户的活跃订单
		// 期货系统合约 -> 个人可用余额
		// 如果 Copy 单，分润 期货系统合约 -> 到邀请人 (不立即到账)
		// 如果 Copy 单，存在 return 剩余金额，总额不会立即到账， 期货系统合约 -> 个人可用余额
		return c.JSON(http.StatusOK, e)
	}

	if e.OrderType == "TPSL_EDIT" {
		// Nothing， Record
		return c.JSON(http.StatusOK, e)
	}

	if e.OrderType == "LIMIT_REQUEST" {
		// 有资金流转，个人可用 -> 个人冻结
		// Nothing， Record
		return c.JSON(http.StatusOK, e)
	}

	if e.OrderType == "LIMIT_CANCEL" {
		// 有资金流转，个人冻结 -> 个人可用
		// Nothing， Record
		return c.JSON(http.StatusOK, e)
	}

	if e.OrderType == "LIMIT_EDIT" {
		// Nothing， Record
		return c.JSON(http.StatusOK, e)
	}

	if e.OrderType == "MARKET_LIQUIDATION" {
		// 爆仓，删掉用户的活跃订单
		// 有资金流转，期货系统合约 -> 个人可用
		return c.JSON(http.StatusOK, e)
	}

	if e.OrderType == "MARKET_ORDER_ADD" {
		// 个人可用 -> 期货系统合约
		// 更新用户的活跃订单
		return c.JSON(http.StatusOK, e)
	}

	// 充值事件
	// 系统合约 -> 个人可用

	// 提现事件
	// 个人可用 -> 系统合约

	// 空投事件
	// 运营账户 -> 个人可用

	return c.NoContent(http.StatusBadRequest)
}
