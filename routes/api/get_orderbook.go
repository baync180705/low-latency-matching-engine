package api

import (
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/baync180705/low-latency-matching-engine/algo"
	types "github.com/baync180705/low-latency-matching-engine/types"
)

func GetOrderBook(c echo.Context) error {
	symbol := c.Param("symbol")
	depthParam := c.QueryParam("depth")

	depth := 10
	if depthParam != "" {
		if d, err := strconv.Atoi(depthParam); err == nil && d > 0 {
			depth = d
		}
	}
	registry := algo.GetRegistry()
	book, exists := registry.Books[symbol]
	if !exists {
		return c.JSON(http.StatusNotFound, map[string]string{
			"error": "Symbol not found",
		})
	}
	book.Mu.Lock()
	defer book.Mu.Unlock()

	resp := &types.OrderBookResponse{
		Symbol:    symbol,
		Timestamp: time.Now().UnixMilli(),
		Bids:      aggregateSide(book.BuyHeap, depth),
		Asks:      aggregateSide(book.SellHeap, depth),
	}

	return c.JSON(http.StatusOK, resp)
}

func aggregateSide(heap *types.Heap, depth int) []types.OrderBookLevel {
	if heap == nil || heap.Len() == 0 {
		return []types.OrderBookLevel{}
	}

	levels := []types.OrderBookLevel{}
	count := 0

	for _, price := range heap.PriceHeap {
		if count >= depth {
			break
		}
		timeQ, exists := heap.TimeQueue[price]
		if !exists {
			continue
		}
		totalQty := int64(0)
		for e := timeQ.Front(); e != nil; e = e.Next() {
			order := e.Value.(*types.Order)
			if !order.IsCancelled && !order.IsComplete {
				totalQty += order.Quantity
			}
		}
		if totalQty > 0 {
			levels = append(levels, types.OrderBookLevel{
				Price:    price,
				Quantity: totalQty,
			})
			count++
		}
	}

	return levels
}
