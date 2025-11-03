package api

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/baync180705/low-latency-matching-engine/engine"
	types "github.com/baync180705/low-latency-matching-engine/types"
)

func CancelOrder(c echo.Context) error {
	orderID := c.Param("order_id")
	registry := engine.GetRegistry()

	var foundOrder *types.Order

	for _, book := range registry.Books {
		book.Mu.Lock()
		order, exists := book.OrderIDMap[orderID]
		book.Mu.Unlock()
		if exists {
			foundOrder = order
			break
		}
	}
	if foundOrder == nil {
		return c.JSON(http.StatusNotFound, map[string]string{
			"error": "Order not found",
		})
	}
	if foundOrder.IsComplete {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Cannot cancel: order already filled",
		})
	}

	foundOrder.IsCancelled = true

	return c.JSON(http.StatusOK, map[string]string{
		"order_id": foundOrder.ID,
		"status":   "CANCELLED",
	})
}
