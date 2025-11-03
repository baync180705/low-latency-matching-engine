package algo

import (
	"errors"
	"time"

	types "github.com/baync180705/low-latency-matching-engine/types"
	"github.com/google/uuid"
	"container/heap"
)

func MatchingAlgorithm (newOrder *types.Order) ([]*types.TradeRecord, error) {
	globalRegistry := GetRegistry()
	
	currentOrderBook := globalRegistry.Books[newOrder.Symbol]

	currentOrderBook.Mu.Lock()
	defer currentOrderBook.Mu.Unlock()
	if !(currentOrderBook.BuyHeap.Len()>0 && currentOrderBook.SellHeap.Len()>0) {
		return []*types.TradeRecord{}, nil
	}
	buyPrice := currentOrderBook.BuyHeap.PriceHeap[0]
	sellPrice := currentOrderBook.SellHeap.PriceHeap[0]

	var tradeResponse []*types.TradeRecord 

	if newOrder.Type=="LIMIT" {
		if buyPrice<sellPrice {return []*types.TradeRecord{}, nil}
		var totalQtyTrade int64 =0

		for buyPrice>=sellPrice {
			bestBuyOrder := currentOrderBook.BuyHeap.TimeQueue[buyPrice].Front().Value.(*types.Order)
			bestSellOrder := currentOrderBook.SellHeap.TimeQueue[sellPrice].Front().Value.(*types.Order)

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

			bestBuyOrder.Quantity-=qtyTrade
			bestSellOrder.Quantity-= qtyTrade

			currentOrderBook.BuyHeap.Qty -=qtyTrade
			currentOrderBook.SellHeap.Qty-=qtyTrade

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

			if !(currentOrderBook.BuyHeap.Len()>0 && currentOrderBook.SellHeap.Len()>0) {
				break
			}

			buyPrice = currentOrderBook.BuyHeap.PriceHeap[0]
			sellPrice= currentOrderBook.SellHeap.PriceHeap[0]

			userTxnQty := min(newOrder.Quantity, totalQtyTrade)

			var restingPrice int64
			if newOrder.Side == "BUY" {
				restingPrice = sellPrice
			} else if newOrder.Side == "SELL" {
				restingPrice = buyPrice
			}

			trade := &types.TradeRecord{
				TradeID:uuid.New().String() ,
				Price: restingPrice,
				Quantity: userTxnQty,
				Timestamp: time.Now().UnixMilli(),
			}

			tradeResponse = append(tradeResponse, trade)
		}
		if totalQtyTrade == newOrder.Quantity {
			newOrder.IsComplete = true
		}
	} else if newOrder.Type=="MARKET" {
		demandedQty := newOrder.Quantity
		var currHeap *types.Heap
		var availableQty int64
		if newOrder.Side=="BUY" {
			availableQty = currentOrderBook.SellHeap.Qty
			currHeap = currentOrderBook.SellHeap
		} else if newOrder.Side =="SELL" {
			availableQty = currentOrderBook.BuyHeap.Qty
			currHeap = currentOrderBook.BuyHeap
		}

		if demandedQty > availableQty {
			newOrder.IsCancelled = true
			return nil, errors.New("Market order could not be filled — insufficient liquidity")
		}

		temp:= demandedQty
		for temp>0 && currHeap.Len() > 0{
			price := currHeap.PriceHeap[0]
			offerOrder := currHeap.TimeQueue[price].Front().Value.(*types.Order)
			offerID := offerOrder.ID
			if currentOrderBook.OrderIDMap[offerID].IsCancelled {
				heap.Pop(currHeap)
				continue
			}
			offerQty := offerOrder.Quantity
			tradeQty := min(offerQty, temp)
			temp -= tradeQty
			offerOrder.Quantity -= tradeQty
			currHeap.Qty-= tradeQty
			if offerOrder.Quantity ==0 {
				offerOrder.IsComplete = true
				currHeap.TimeQueue[price].Remove(currHeap.TimeQueue[price].Front())
				if currHeap.TimeQueue[price].Len() ==0 {
					delete(currHeap.TimeQueue, price)
					heap.Pop(currHeap)
				}
			}

			trade := &types.TradeRecord{
				TradeID: uuid.New().String() ,
				Price: price,
				Quantity: tradeQty,
				Timestamp: time.Now().UnixMilli(),
			}
			tradeResponse = append(tradeResponse, trade)
		}
		newOrder.IsComplete = true
	}
	return tradeResponse, nil
}