package api

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/baync180705/low-latency-matching-engine/algo"
	types "github.com/baync180705/low-latency-matching-engine/types"
)

func GetOrderStatus(c echo.Context) error {
	orderID := c.Param("order_id")
	globalRegistry := algo.GetRegistry()

	var foundOrder *types.Order

	for _, book := range globalRegistry.Books {
		if order, exists := book.OrderIDMap[orderID]; exists {
			foundOrder = order
			break
		}
	}
	if foundOrder == nil {
		return c.JSON(http.StatusNotFound, map[string]string{
			"error": "Order not found",
		})
	}
	var status string
	switch {
	case foundOrder.IsCancelled:
		status = "CANCELLED"
	case foundOrder.IsComplete && foundOrder.Quantity == 0:
		status = "FILLED"
	case foundOrder.IsComplete && foundOrder.Quantity > 0:
		status = "PARTIAL_FILL"
	default:
		status = "ACCEPTED"
	}

	filledQty := foundOrder.InitQty - foundOrder.Quantity

	response := types.StatusResponse{
		OrderID:        foundOrder.ID,
		Symbol:         foundOrder.Symbol,
		Side:           foundOrder.Side,
		Type:           foundOrder.Type,
		Price:          foundOrder.Price,
		Quantity:       foundOrder.InitQty,
		FilledQuantity: filledQty,
		Status:         status,
		Timestamp:      foundOrder.Timestamp,
	}

	return c.JSON(http.StatusOK, response)
}
