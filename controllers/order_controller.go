package controller

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/baync180705/low-latency-matching-engine/algo"
	types "github.com/baync180705/low-latency-matching-engine/types"
)

func SubmitOrder(c echo.Context) error {
	var input types.OrderInput
	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, types.OrderResponse{
			Error: "Invalid request body",
		})
	}

	order, trades, err := algo.RunPipeline(&input)
	if err != nil {
		//This is for matching or validation failure
		return c.JSON(http.StatusBadRequest, types.OrderResponse{
			Error: err.Error(),
		})
	}

	var filledQty int64
	for _, t := range trades {
		filledQty += t.Quantity
	}
	remainingQty := order.Quantity - filledQty

	resp := types.OrderResponse{
		OrderID: order.ID,
		Trades:  trades,
	}

	switch {
	case filledQty == 0:
		resp.Status = "ACCEPTED"
		resp.Message = "Order added to book"
		return c.JSON(http.StatusCreated, resp)

	case remainingQty > 0:
		resp.Status = "PARTIAL_FILL"
		resp.FilledQuantity = filledQty
		resp.RemainingQuantity = remainingQty
		return c.JSON(http.StatusAccepted, resp)

	default:
		resp.Status = "FILLED"
		resp.FilledQuantity = filledQty
		return c.JSON(http.StatusOK, resp)
	}
}
