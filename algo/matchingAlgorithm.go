package algo

import (
	types "github.com/baync180705/low-latency-matching-engine/types"
)

func MatchingAlgorithm (newOrder *types.Order) {
	globalRegistry := GetRegistry()
	
	currentOrderBook := globalRegistry.Books[newOrder.Symbol]

	currentOrderBook.Mu.Lock()
	defer currentOrderBook.Mu.Unlock()
	if !(currentOrderBook.BuyHeap.Len()>0 && currentOrderBook.SellHeap.Len()>0) {
		return
	}
	buyPrice := currentOrderBook.BuyHeap.PriceHeap[0]
	sellPrice := currentOrderBook.SellHeap.PriceHeap[0]

	if newOrder.Type=="LIMIT" {
		if buyPrice<sellPrice {return}

		for buyPrice>=sellPrice {
			bestBuyOrder := currentOrderBook.BuyHeap.TimeQueue[buyPrice].Front().Value.(*types.Order)
			bestSellOrder := currentOrderBook.SellHeap.TimeQueue[sellPrice].Front().Value.(*types.Order)

			buyQty := bestBuyOrder.Quantity
			sellQty := bestSellOrder.Quantity

			qtyTrade := min(buyQty, sellQty)

			bestBuyOrder.Quantity-=qtyTrade
			bestSellOrder.Quantity-= qtyTrade

			currentOrderBook.BuyHeap.Qty -=qtyTrade
			currentOrderBook.SellHeap.Qty-=qtyTrade

			if bestBuyOrder.Quantity ==0 {
				currentOrderBook.BuyHeap.TimeQueue[buyPrice].Remove(currentOrderBook.BuyHeap.TimeQueue[buyPrice].Front())
				delete(currentOrderBook.OrderIDMap, bestBuyOrder.ID)
				if currentOrderBook.BuyHeap.TimeQueue[buyPrice].Len() == 0 {
					delete(currentOrderBook.BuyHeap.TimeQueue, buyPrice)
					currentOrderBook.BuyHeap.Pop()
				}
			} 

			if bestSellOrder.Quantity ==0 {
				currentOrderBook.SellHeap.TimeQueue[sellPrice].Remove(currentOrderBook.SellHeap.TimeQueue[sellPrice].Front())
				delete(currentOrderBook.OrderIDMap, bestSellOrder.ID)
				if currentOrderBook.SellHeap.TimeQueue[sellPrice].Len() == 0 {
					delete(currentOrderBook.SellHeap.TimeQueue, sellPrice)
					currentOrderBook.SellHeap.Pop()
				}
			}

			if !(currentOrderBook.BuyHeap.Len()>0 && currentOrderBook.SellHeap.Len()>0) {
				return
			}

			buyPrice = currentOrderBook.BuyHeap.PriceHeap[0]
			sellPrice= currentOrderBook.SellHeap.PriceHeap[0]
		}
	} else if newOrder.Type=="MARKET" {
		demandedQty := newOrder.Quantity
		var heap *types.Heap
		var availableQty int64
		if newOrder.Side=="BUY" {
			availableQty = currentOrderBook.SellHeap.Qty
			heap = currentOrderBook.SellHeap
		} else if newOrder.Side =="SELL" {
			availableQty = currentOrderBook.BuyHeap.Qty
			heap = currentOrderBook.BuyHeap
		}

		if demandedQty>availableQty {return}

		temp:= demandedQty
		for temp>0 && heap.Len() > 0{
			price := heap.PriceHeap[0]
			offerOrder := heap.TimeQueue[price].Front().Value.(*types.Order)
			offerQty := offerOrder.Quantity
			tradeQty := min(offerQty, temp)
			temp -= tradeQty
			offerOrder.Quantity -= tradeQty
			heap.Qty-= tradeQty
			if offerOrder.Quantity ==0 {
				heap.TimeQueue[price].Remove(heap.TimeQueue[price].Front())
				if heap.TimeQueue[price].Len() ==0 {
					delete(heap.TimeQueue, price)
					delete(currentOrderBook.OrderIDMap, offerOrder.ID)
					heap.Pop()
				}
			}
		}
		delete(currentOrderBook.OrderIDMap, newOrder.ID)
	}
}