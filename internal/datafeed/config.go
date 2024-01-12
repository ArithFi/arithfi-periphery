package datafeed

import (
	"github.com/arithfi/arithfi-periphery/model"
	"github.com/labstack/echo/v4"
	"net/http"
)

func GetConfig(c echo.Context) error {
	config := &model.Config{
		// supports_search: Set it to true if your data feed supports symbol search and individual symbol resolve logic.
		SupportsSearch: true,
		// supports_group_request: Set it to true if your data feed provides full information on symbol group only and is not able to perform symbol search or individual symbol resolve.
		SupportsGroupRequest:   false,
		SupportsMarks:          false,
		SupportsTimescaleMarks: false,
		SupportedResolutions:   []string{"1", "3", "5", "15", "30", "60", "120", "240", "360", "480", "720", "1D", "3D", "1W", "1M"},
	}

	return c.JSON(http.StatusOK, config)
}
