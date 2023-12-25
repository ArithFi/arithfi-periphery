package market_data

import (
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

/*
PingTest Test connectivity to the Rest API.
Weight: 1
Parameters: NONE
*/
func TestPing(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	h := Ping(c)
	if assert.NoError(t, h) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, "{}\n", rec.Body.String())
	}
}
