package api

import (
	"net/http"
	"log"

	"github.com/baync180705/low-latency-matching-engine/engine"
	types "github.com/baync180705/low-latency-matching-engine/types"
	"github.com/labstack/echo/v4"
)

func SubmitOrder(c echo.Context) error {
	var input types.OrderInput
	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, types.OrderResponse{
			Error: "Invalid request body",
		})
	}
	
	order, trades, err := engine.RunPipeline(&input)
	
	// Handle errors from matching (rejections)
	if err != nil {
		// Check if order was cancelled due to insufficient liquidity
		if order != nil && order.IsCancelled {
			return c.JSON(http.StatusOK, types.OrderResponse{
				OrderID: order.ID,
				Status:  "REJECTED",
				Error:   err.Error(),
			})
		}
		
		return c.JSON(http.StatusBadRequest, types.OrderResponse{
			Error: err.Error(),
		})
	}
	
	// Handle no trades (order added to book)
	if len(trades) == 0 {
		return c.JSON(http.StatusOK, types.OrderResponse{
			OrderID: order.ID,
			Status:  "ACCEPTED",
			Message: "No match found — order added to order book",
		})
	}
	
	// Calculate filled quantity
	var filledQty int64
	for _, t := range trades {
		filledQty += t.Quantity
	}
	remainingQty := input.Quantity - filledQty
	
	var resp types.OrderResponse
	
	// Skip default/empty trades
	if trades[0].Quantity == 0 && trades[0].Price == 0 {
		resp = types.OrderResponse{OrderID: order.ID}
	} else {
		resp = types.OrderResponse{
			OrderID: order.ID,
			Trades:  trades,
		}
	}
	
	log.Printf("Filled Quantity: %d, Remaining Quantity: %d, Initial Qty: %d", 
		filledQty, remainingQty, input.Quantity)
	
	switch {
	case order.IsCancelled:
		resp.Status = "REJECTED"
		resp.Error = "Order was rejected"
		return c.JSON(http.StatusOK, resp)
		
	case order.IsComplete && filledQty == input.Quantity:
		resp.Status = "FILLED"
		resp.FilledQuantity = filledQty
		return c.JSON(http.StatusOK, resp)
		
	case filledQty > 0 && remainingQty > 0:
		resp.Status = "PARTIAL_FILL"
		resp.FilledQuantity = filledQty
		resp.RemainingQuantity = remainingQty
		return c.JSON(http.StatusAccepted, resp)
		
	default:
		resp.Status = "ACCEPTED"
		resp.Message = "Order added to book"
		return c.JSON(http.StatusCreated, resp)
	}
}