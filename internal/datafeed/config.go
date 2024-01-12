package datafeed

import (
	"github.com/arithfi/arithfi-periphery/model"
	"github.com/labstack/echo/v4"
	"net/http"
)

func GetConfig(c echo.Context) error {
	config := &model.Config{
		SupportsSearch:         true,
		SupportsGroupRequest:   false,
		SupportsMarks:          false,
		SupportsTimescaleMarks: false,
		SupportsTime:           true,
		SupportedResolutions:   []string{"1", "3", "5", "15", "30", "60", "120", "240", "360", "480", "720", "1D", "3D", "1W", "1M"},
		Exchanges:              []model.Info{{Value: "Binance", Name: "Binance Exchange", Desc: "Binance"}},
		SymbolsTypes:           []model.Info{{Value: "crypto", Name: "Crypto"}},
	}

	return c.JSON(http.StatusOK, config)
}
