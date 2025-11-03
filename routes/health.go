package routes

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/baync180705/low-latency-matching-engine/algo"
)

var startTime = time.Now()

func HealthCheck(c echo.Context) error {
	globalRegistry := algo.GetRegistry()
	uptime := int64(time.Since(startTime).Seconds())
	totalOrders := 0
	for _, book := range globalRegistry.Books {
		book.Mu.Lock()
		totalOrders += len(book.OrderIDMap)
		book.Mu.Unlock()
	}
	response := map[string]interface{}{
		"status":           "healthy",
		"uptime_seconds":   uptime,
		"orders_processed": totalOrders,
	}
	return c.JSON(http.StatusOK, response)
}
