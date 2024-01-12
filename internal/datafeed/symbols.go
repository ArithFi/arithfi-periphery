package datafeed

import (
	"github.com/arithfi/arithfi-periphery/internal/binance"
	"github.com/arithfi/arithfi-periphery/model"
	"github.com/labstack/echo/v4"
	"log"
	"net/http"
	"strconv"
)

func findFilter(filters *[]model.SymbolFilter) string {
	for _, f := range *filters {
		if f.FilterType == "PRICE_FILTER" {
			return f.TickSize
		}
	}
	return ""
}

func Symbols(c echo.Context) error {
	symbol := c.QueryParam("symbol")
	if symbol == "" {
		return c.JSON(http.StatusBadRequest, model.UDFError{S: "error", Errmsg: "symbol: 404 not found"})
	}
	exchangeInfo := binance.GetExchangeInfo()
	for _, symbolInfo := range exchangeInfo.Symbols {
		if symbolInfo.Symbol == symbol {
			tickerSize := findFilter(&symbolInfo.Filters)

			pscale, err := strconv.ParseFloat(tickerSize, 64)
			if err != nil {
				log.Println("Symbols : " + symbolInfo.Symbol)
				log.Println("TickSize : " + tickerSize)
				return c.JSON(http.StatusConflict, model.UDFError{S: "error", Errmsg: "tickerSize calc ERROR"})
			}

			a := &model.Symbol{
				Symbol:               symbolInfo.Symbol,
				Ticker:               symbolInfo.Symbol,
				Name:                 symbolInfo.BaseAsset + " Coin",
				FullName:             "[" + symbolInfo.QuoteAsset + "]" + symbolInfo.BaseAsset + " Coin",
				Description:          symbolInfo.BaseAsset + " / " + symbolInfo.QuoteAsset,
				Exchange:             "BINANCE",
				ListedExchange:       "BINANCE",
				Type:                 "crypto",
				CurrencyCode:         symbolInfo.QuoteAsset,
				Session:              "24x7",
				Timezone:             "UTC",
				Minmovent:            1,
				Minmov:               1,
				Minmovement2:         0,
				Minmov2:              0,
				Pricescale:           int64(1 / pscale),
				SupportedResolutions: []string{"1", "3", "5", "15", "30", "60", "120", "240", "360", "480", "720", "1D", "3D", "1W", "1M"},
				HasIntraday:          true,
				HasDaily:             true,
				HasWeeklyAndMonthly:  true,
				DataStatus:           "streaming",
			}
			return c.JSON(http.StatusOK, a)
		}
	}
	return c.JSON(http.StatusOK, model.UDFError{S: "error, symbol: 404 not found"})
}
