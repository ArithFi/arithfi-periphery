package datafeed

import (
	"github.com/arithfi/arithfi-periphery/model"
	"github.com/labstack/echo/v4"
	"net/http"
)

var pricescaleMap = map[string]int64{
	"ETH/USDT":   100,
	"BTC/USDT":   100,
	"BNB/USDT":   100,
	"MATIC/USDT": 10000,
	"ADA/USDT":   10000,
	"DOGE/USDT":  100000,
	"XRP/USDT":   1000,
	"SOL/USDT":   100,
	"LTC/USDT":   100,
	"AVAX/USDT":  1000,
	"AUD/USD":    100000,
	"EUR/USD":    100000,
	"USD/JPY":    1000,
	"USD/CAD":    100000,
	"GBP/USD":    100000,
}

func Symbols(c echo.Context) error {
	symbol := c.QueryParam("symbol")
	if symbol == "" {
		return c.JSON(http.StatusBadRequest, model.UDFError{S: "e", Errmsg: "symbol: 404 not found"})
	}

	pricescale, e := pricescaleMap[symbol]
	if e != true {
		pricescale = 100
	}

	a := &model.Symbol{
		Symbol:               symbol,
		Ticker:               symbol,
		Name:                 symbol,
		FullName:             symbol,
		Description:          symbol,
		Exchange:             "BINANCE",
		ListedExchange:       "BINANCE",
		Type:                 "crypto",
		CurrencyCode:         "USD",
		Session:              "24x7",
		Timezone:             "UTC",
		Minmovent:            1,
		Minmov:               1,
		Minmovement2:         0,
		Minmov2:              0,
		Pricescale:           pricescale,
		SupportedResolutions: []string{"1", "3", "5", "15", "30", "60", "120", "240", "360", "480", "720", "1D", "3D", "1W", "1M"},
		HasIntraday:          true,
		HasDaily:             true,
		HasWeeklyAndMonthly:  true,
		DataStatus:           "streaming",
	}
	return c.JSON(http.StatusOK, a)
}
