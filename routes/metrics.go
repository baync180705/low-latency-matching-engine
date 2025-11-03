package routes

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/baync180705/low-latency-matching-engine/engine"
	"github.com/baync180705/low-latency-matching-engine/metrics"
)

// I return a live metrics snapshot combining counters and registry state.
func GetMetricsHandler(c echo.Context) error {
	snap := metrics.GetSnapshot(engine.GetRegistry())
	return c.JSON(http.StatusOK, snap)
}
