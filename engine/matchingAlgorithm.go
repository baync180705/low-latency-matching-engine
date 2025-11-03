package engine

import (
	"container/heap"
	"errors"
	"time"

	types "github.com/baync180705/low-latency-matching-engine/types"
	"github.com/google/uuid"
)

func MatchingAlgorithm(newOrder *types.Order) ([]*types.TradeRecord, error) {
	globalRegistry := GetRegistry()
	currentOrderBook := globalRegistry.Books[newOrder.Symbol]

	currentOrderBook.Mu.Lock()
	defer currentOrderBook.Mu.Unlock()

	var tradeResponse []*types.TradeRecord

	if newOrder.Type == "LIMIT" {
		// For LIMIT orders, check if matching is possible
		if !(currentOrderBook.BuyHeap.Len() > 0 && currentOrderBook.SellHeap.Len() > 0) {
			return []*types.TradeRecord{}, nil
		}

		buyPrice := currentOrderBook.BuyHeap.PriceHeap[0]
		sellPrice := currentOrderBook.SellHeap.PriceHeap[0]

		// No cross - can't match
		if buyPrice < sellPrice {
			return []*types.TradeRecord{}, nil
		}

		var totalQtyTrade int64 = 0

		for buyPrice >= sellPrice {
			bestBuyOrder := currentOrderBook.BuyHeap.TimeQueue[buyPrice].Front().Value.(*types.Order)
			bestSellOrder := currentOrderBook.SellHeap.TimeQueue[sellPrice].Front().Value.(*types.Order)

			// Handle cancelled buy orders
			if currentOrderBook.OrderIDMap[bestBuyOrder.ID].IsCancelled {
				currentOrderBook.BuyHeap.TimeQueue[buyPrice].Remove(currentOrderBook.BuyHeap.TimeQueue[buyPrice].Front())
				if currentOrderBook.BuyHeap.TimeQueue[buyPrice].Len() == 0 {
					delete(currentOrderBook.BuyHeap.TimeQueue, buyPrice)
					heap.Pop(currentOrderBook.BuyHeap)
				}
				if currentOrderBook.BuyHeap.Len() == 0 {
					break
				}
				buyPrice = currentOrderBook.BuyHeap.PriceHeap[0]
				continue
			}

			// Handle cancelled sell orders
			if currentOrderBook.OrderIDMap[bestSellOrder.ID].IsCancelled {
				currentOrderBook.SellHeap.TimeQueue[sellPrice].Remove(currentOrderBook.SellHeap.TimeQueue[sellPrice].Front())
				if currentOrderBook.SellHeap.TimeQueue[sellPrice].Len() == 0 {
					delete(currentOrderBook.SellHeap.TimeQueue, sellPrice)
					heap.Pop(currentOrderBook.SellHeap)
				}
				if currentOrderBook.SellHeap.Len() == 0 {
					break
				}
				sellPrice = currentOrderBook.SellHeap.PriceHeap[0]
				continue
			}

			buyQty := bestBuyOrder.Quantity
			sellQty := bestSellOrder.Quantity

			qtyTrade := min(buyQty, sellQty)
			totalQtyTrade += qtyTrade

			bestBuyOrder.Quantity -= qtyTrade
			bestSellOrder.Quantity -= qtyTrade

			currentOrderBook.BuyHeap.Qty -= qtyTrade
			currentOrderBook.SellHeap.Qty -= qtyTrade

			// Remove fully filled orders
			if bestBuyOrder.Quantity == 0 {
				bestBuyOrder.IsComplete = true
				currentOrderBook.BuyHeap.TimeQueue[buyPrice].Remove(currentOrderBook.BuyHeap.TimeQueue[buyPrice].Front())
				if currentOrderBook.BuyHeap.TimeQueue[buyPrice].Len() == 0 {
					delete(currentOrderBook.BuyHeap.TimeQueue, buyPrice)
					heap.Pop(currentOrderBook.BuyHeap)
				}
			}

			if bestSellOrder.Quantity == 0 {
				bestSellOrder.IsComplete = true
				currentOrderBook.SellHeap.TimeQueue[sellPrice].Remove(currentOrderBook.SellHeap.TimeQueue[sellPrice].Front())
				if currentOrderBook.SellHeap.TimeQueue[sellPrice].Len() == 0 {
					delete(currentOrderBook.SellHeap.TimeQueue, sellPrice)
					heap.Pop(currentOrderBook.SellHeap)
				}
			}

			// Create trade record for this match
			var tradePrice int64
			if newOrder.Side == "BUY" {
				tradePrice = sellPrice
			} else if newOrder.Side == "SELL" {
				tradePrice = buyPrice
			}

			trade := &types.TradeRecord{
				TradeID:   uuid.New().String(),
				Price:     tradePrice,
				Quantity:  qtyTrade,
				Timestamp: time.Now().UnixMilli(),
			}
			tradeResponse = append(tradeResponse, trade)

			// Check if we can continue matching
			if !(currentOrderBook.BuyHeap.Len() > 0 && currentOrderBook.SellHeap.Len() > 0) {
				break
			}

			buyPrice = currentOrderBook.BuyHeap.PriceHeap[0]
			sellPrice = currentOrderBook.SellHeap.PriceHeap[0]
		}

		// Mark complete if fully filled
		if totalQtyTrade >= newOrder.Quantity {
			newOrder.IsComplete = true
		}

	} else if newOrder.Type == "MARKET" {
		// Market orders match against the OPPOSITE side
		demandedQty := newOrder.Quantity
		var currHeap *types.Heap

		if newOrder.Side == "BUY" {
			currHeap = currentOrderBook.SellHeap
		} else if newOrder.Side == "SELL" {
			currHeap = currentOrderBook.BuyHeap
		}

		// Check if opposite side has any liquidity
		if currHeap == nil || currHeap.Len() == 0 {
			newOrder.IsCancelled = true
			return nil, errors.New("No liquidity available")
		}

		availableQty := currHeap.Qty

		// Strict mode: reject if insufficient liquidity for full fill
		if demandedQty > availableQty {
			newOrder.IsCancelled = true
			return nil, errors.New("Insufficient liquidity")
		}

		temp := demandedQty
		for temp > 0 && currHeap.Len() > 0 {
			price := currHeap.PriceHeap[0]
			offerOrder := currHeap.TimeQueue[price].Front().Value.(*types.Order)
			offerID := offerOrder.ID

			// Skip cancelled orders
			if currentOrderBook.OrderIDMap[offerID].IsCancelled {
				currHeap.TimeQueue[price].Remove(currHeap.TimeQueue[price].Front())
				if currHeap.TimeQueue[price].Len() == 0 {
					delete(currHeap.TimeQueue, price)
					heap.Pop(currHeap)
				}
				continue
			}

			offerQty := offerOrder.Quantity
			tradeQty := min(offerQty, temp)
			temp -= tradeQty
			offerOrder.Quantity -= tradeQty
			currHeap.Qty -= tradeQty

			// Remove fully filled resting order
			if offerOrder.Quantity == 0 {
				offerOrder.IsComplete = true
				currHeap.TimeQueue[price].Remove(currHeap.TimeQueue[price].Front())
				if currHeap.TimeQueue[price].Len() == 0 {
					delete(currHeap.TimeQueue, price)
					heap.Pop(currHeap)
				}
			}

			trade := &types.TradeRecord{
				TradeID:   uuid.New().String(),
				Price:     price,
				Quantity:  tradeQty,
				Timestamp: time.Now().UnixMilli(),
			}
			tradeResponse = append(tradeResponse, trade)
		}

		// Update market order quantity to reflect unfilled portion
		newOrder.Quantity = temp

		// Mark complete if fully filled
		if temp == 0 {
			newOrder.IsComplete = true
		}
	}

	return tradeResponse, nil
}
