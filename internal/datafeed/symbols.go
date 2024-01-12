package datafeed

import (
	"github.com/arithfi/arithfi-periphery/model"
	"github.com/labstack/echo/v4"
	"net/http"
	"strings"
)

func Symbols(c echo.Context) error {
	symbol := c.QueryParam("symbol")
	if symbol == "" {
		return c.JSON(http.StatusBadRequest, model.UDFError{S: "error", Errmsg: "symbol: 404 not found"})
	}
	pair1 := strings.Split(symbol, "/")[0]
	pair2 := strings.Split(symbol, "/")[1]

	a := &model.Symbol{
		Symbol:               pair1 + pair2,
		Ticker:               pair1 + pair2,
		Name:                 pair1 + pair2,
		FullName:             symbol,
		Description:          symbol,
		Exchange:             "BINANCE",
		ListedExchange:       "BINANCE",
		Type:                 "crypto",
		CurrencyCode:         pair2,
		Session:              "24x7",
		Timezone:             "UTC",
		Minmovent:            1,
		Minmov:               1,
		Minmovement2:         0,
		Minmov2:              0,
		Pricescale:           6,
		SupportedResolutions: []string{"1", "3", "5", "15", "30", "60", "120", "240", "360", "480", "720", "1D", "3D", "1W", "1M"},
		HasIntraday:          true,
		HasDaily:             true,
		HasWeeklyAndMonthly:  true,
		DataStatus:           "streaming",
	}
	return c.JSON(http.StatusOK, a)
}
