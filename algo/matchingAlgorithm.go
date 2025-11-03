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

			bestBuyOrder.Quantity = bestBuyOrder.Quantity - qtyTrade
			bestSellOrder.Quantity = bestSellOrder.Quantity - qtyTrade

			if bestBuyOrder.Quantity ==0 {
				currentOrderBook.BuyHeap.TimeQueue[buyPrice].Remove(currentOrderBook.BuyHeap.TimeQueue[buyPrice].Front())
				if currentOrderBook.BuyHeap.TimeQueue[buyPrice].Len() == 0 {
					delete(currentOrderBook.BuyHeap.TimeQueue, buyPrice)
					delete(currentOrderBook.OrderIDMap, bestBuyOrder.ID)
					currentOrderBook.BuyHeap.Pop()
				}
			} 

			if bestSellOrder.Quantity ==0 {
				currentOrderBook.SellHeap.TimeQueue[sellPrice].Remove(currentOrderBook.SellHeap.TimeQueue[sellPrice].Front())
				if currentOrderBook.SellHeap.TimeQueue[sellPrice].Len() == 0 {
					delete(currentOrderBook.SellHeap.TimeQueue, sellPrice)
					delete(currentOrderBook.OrderIDMap, bestSellOrder.ID)
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
		
	}

	

	
}