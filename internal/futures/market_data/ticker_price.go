package market_data

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

/*
TickerPrice Latest price for a symbol or symbols.
Weight:

	1 for a single symbol;
	2 when the symbol parameter is omitted

Parameters:

	symbol	STRING	NO
	pair	STRING	NO

- Symbol and pair cannot be sent together
- If a pair is sent,tickers for all symbols of the pair will be returned
- If either a pair or symbol is sent, tickers for all symbols of all pairs will be returned
*/
func TickerPrice(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{
		"message": "pong",
	})
}
