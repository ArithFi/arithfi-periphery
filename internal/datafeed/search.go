package datafeed

import (
	"github.com/arithfi/arithfi-periphery/internal/binance"
	"github.com/arithfi/arithfi-periphery/model"
	"github.com/labstack/echo/v4"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func Filter(arr *[]model.SymbolInfo, f func(model.SymbolInfo) bool) {
	result := make([]model.SymbolInfo, 0)
	for _, v := range *arr {
		if f(v) {
			result = append(result, v)
		}
	}
	*arr = result
}

func Search(c echo.Context) error {
	query := c.QueryParam("query")
	ptype := c.QueryParam("type")
	exchange := c.QueryParam("exchange")
	limit := c.QueryParam("limit")

	symbols := binance.GetExchangeInfo().Symbols
	if query != "" {
		Filter(&symbols, func(v model.SymbolInfo) bool {
			return strings.Contains(v.Symbol, query)
		})
	}
	if ptype != "" {
		// Filter(&symbols, func(v model.SymbolInfo) bool {
		// 	return v.type ==ptype
		// })
	}
	if exchange != "" {
		// Filter(&symbols, func(v model.SymbolInfo) bool {
		// 	return strings.Contains(v.Symbol, query)
		// })
	}
	if limit != "" {
		lm, err := strconv.Atoi(limit)
		if err != nil {
			log.Println("limit : " + limit)
			return c.JSON(http.StatusConflict, model.UDFError{S: "error", Errmsg: "limit parse error"})
		}
		return c.JSON(http.StatusOK, symbols[:lm])
	}

	return c.JSON(http.StatusOK, symbols)
}
