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

var ResolutionMap = map[string]string{
	"1":   "1m",
	"3":   "3m",
	"5":   "5m",
	"15":  "15m",
	"30":  "30m",
	"60":  "1h",
	"120": "2h",
	"240": "4h",
	"360": "6h",
	"480": "8h",
	"720": "12h",
	"D":   "1d",
	"1D":  "1d",
	"3D":  "3d",
	"W":   "1w",
	"1W":  "1w",
	"M":   "1M",
	"1M":  "1M",
}

func History(c echo.Context) error {
	symbol := c.QueryParam("symbol")
	resolution := ResolutionMap[c.QueryParam("resolution")]
	from, err := strconv.ParseInt(c.QueryParam("from"), 0, 64)
	if err != nil {
		log.Println("limit : " + c.QueryParam("from"))
		return c.JSON(http.StatusConflict, model.UDFError{S: "error", Errmsg: "from parse error"})
	}
	to, err := strconv.ParseInt(c.QueryParam("to"), 0, 64)
	if err != nil {
		log.Println("limit : " + c.QueryParam("to"))
		return c.JSON(http.StatusConflict, model.UDFError{S: "error", Errmsg: "to parse error"})
	}
	symbol = strings.ReplaceAll(symbol, "/", "")
	log.Println("symbol", symbol)
	klines := binance.GetKlines(symbol, resolution, from*1000, to*1000)

	result := make([]model.Bar, len(*klines))
	for i, data := range *klines {
		result[i] = model.Bar{
			S: "ok",
			T: data.OpenTime / 1000,
			C: data.Close,
			O: data.Open,
			H: data.High,
			L: data.Low,
			V: data.Volume,
		}
	}

	return c.JSON(http.StatusOK, result)
}
