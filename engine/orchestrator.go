package engine

import (
	"errors"

	types "github.com/baync180705/low-latency-matching-engine/types"
)

func RunPipeline(input *types.OrderInput) (*types.Order, []*types.TradeRecord, error) {
	order, err := SubmitOrderEntry(input)
	if err != nil {
		return order, nil, errors.New("Failed to submit your order request")
	}

	globalRegistry := GetRegistry()
	book := globalRegistry.Books[order.Symbol]

	// Market orders always try to match immediately
	if order.Type == "MARKET" {
		tradeResponse, err := MatchingAlgorithm(order)
		return order, tradeResponse, err
	}

	// For LIMIT orders, check if this order can trigger matching
	// Only the order at the top of its side can trigger matching
	var heap *types.Heap
	if order.Side == "BUY" {
		heap = book.BuyHeap
	} else if order.Side == "SELL" {
		heap = book.SellHeap
	}

	if heap.Len() == 0 {
		return order, []*types.TradeRecord{}, nil
	}

	price := heap.PriceHeap[0]
	if heap.TimeQueue[price].Len() == 0 {
		return order, []*types.TradeRecord{}, nil
	}

	entry := heap.TimeQueue[price].Front().Value.(*types.Order)

	// Only match if this order is at the top of the book
	if entry.ID == order.ID {
		tradeResponse, err := MatchingAlgorithm(order)
		return order, tradeResponse, err
	}

	// Order is not at the top, just queued
	return order, []*types.TradeRecord{}, nil
}
