package datafeed

import (
	"github.com/arithfi/arithfi-periphery/internal/binance"
	"github.com/arithfi/arithfi-periphery/internal/forex"
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
		log.Println("from : " + c.QueryParam("from"))
		return c.JSON(http.StatusConflict, model.UDFError{S: "error", Errmsg: "from parse error"})
	}
	to, err := strconv.ParseInt(c.QueryParam("to"), 0, 64)
	if err != nil {
		log.Println("to : " + c.QueryParam("to"))
		return c.JSON(http.StatusConflict, model.UDFError{S: "error", Errmsg: "to parse error"})
	}
	symbol = strings.ReplaceAll(symbol, "/", "")

	if strings.Contains(symbol, "USDT") {
		klines := binance.GetKlines(symbol, resolution, from*1000, to*1000)

		result := model.Bar{}
		if klines == nil {
			result.S = "no_data"
			return c.JSON(http.StatusOK, result)
		}
		result.S = "ok"
		for _, data := range *klines {
			result.T = append(result.T, data.OpenTime)
			openPrice, _ := strconv.ParseFloat(data.Open, 64)
			result.O = append(result.O, openPrice)
			highPrice, _ := strconv.ParseFloat(data.High, 64)
			result.H = append(result.H, highPrice)
			lowPrice, _ := strconv.ParseFloat(data.Low, 64)
			result.L = append(result.L, lowPrice)
			closePrice, _ := strconv.ParseFloat(data.Close, 64)
			result.C = append(result.C, closePrice)
			volume, _ := strconv.ParseFloat(data.Volume, 64)
			result.V = append(result.V, volume)
		}

		return c.JSON(http.StatusOK, result)
	} else {
		klines := forex.GetKlines(symbol, resolution, from*1000, to*1000)

		result := model.Bar{}
		if klines == nil || len(*klines) == 0 {
			result.S = "no_data"
			result.T = []int64{}
			result.O = []float64{}
			result.H = []float64{}
			result.L = []float64{}
			result.C = []float64{}
			result.V = []float64{}
			return c.JSON(http.StatusBadRequest, result)
		}
		result.S = "ok"
		for _, data := range *klines {
			result.T = append(result.T, data.OpenTime)
			openPrice, _ := strconv.ParseFloat(data.Open, 64)
			result.O = append(result.O, openPrice)
			highPrice, _ := strconv.ParseFloat(data.High, 64)
			result.H = append(result.H, highPrice)
			lowPrice, _ := strconv.ParseFloat(data.Low, 64)
			result.L = append(result.L, lowPrice)
			closePrice, _ := strconv.ParseFloat(data.Close, 64)
			result.C = append(result.C, closePrice)
			volume, _ := strconv.ParseFloat(data.Volume, 64)
			result.V = append(result.V, volume)
		}
		return c.JSON(http.StatusOK, result)
	}
}
