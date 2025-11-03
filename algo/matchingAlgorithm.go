package algo

import (
	"errors"
	"time"

	types "github.com/baync180705/low-latency-matching-engine/types"
	"github.com/google/uuid"
)

func MatchingAlgorithm (newOrder *types.Order) ([]*types.TradeRecord, error) {
	globalRegistry := GetRegistry()
	
	currentOrderBook := globalRegistry.Books[newOrder.Symbol]

	currentOrderBook.Mu.Lock()
	defer currentOrderBook.Mu.Unlock()
	if !(currentOrderBook.BuyHeap.Len()>0 && currentOrderBook.SellHeap.Len()>0) {
		return nil, errors.New("No match found — order added to order book")
	}
	buyPrice := currentOrderBook.BuyHeap.PriceHeap[0]
	sellPrice := currentOrderBook.SellHeap.PriceHeap[0]

	var tradeResponse []*types.TradeRecord 

	if newOrder.Type=="LIMIT" {
		if buyPrice<sellPrice {return nil, errors.New("No match found — order added to order book")}
		var totalQtyTrade int64 =0

		for buyPrice>=sellPrice {
			bestBuyOrder := currentOrderBook.BuyHeap.TimeQueue[buyPrice].Front().Value.(*types.Order)
			bestSellOrder := currentOrderBook.SellHeap.TimeQueue[sellPrice].Front().Value.(*types.Order)

			buyQty := bestBuyOrder.Quantity
			sellQty := bestSellOrder.Quantity

			qtyTrade := min(buyQty, sellQty)
			totalQtyTrade += qtyTrade

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
				break
			}

			buyPrice = currentOrderBook.BuyHeap.PriceHeap[0]
			sellPrice= currentOrderBook.SellHeap.PriceHeap[0]
		}

		var userTxnQty int64
		if newOrder.Quantity > totalQtyTrade {
			userTxnQty = totalQtyTrade
		} else {
			userTxnQty = newOrder.Quantity
		}

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

		if demandedQty>availableQty {
			delete(currentOrderBook.OrderIDMap, newOrder.ID)
			return nil, errors.New("Market order could not be filled — insufficient liquidity")
		}

		temp:= demandedQty
		for temp>0 && heap.Len() > 0{
			price := heap.PriceHeap[0]
			offerOrder := heap.TimeQueue[price].Front().Value.(*types.Order)
			offerID := offerOrder.ID
			if _, exists := currentOrderBook.OrderIDMap[offerID]; !exists {
				heap.Pop()
				continue
			}
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

			trade := &types.TradeRecord{
				TradeID: uuid.New().String() ,
				Price: price,
				Quantity: tradeQty,
				Timestamp: time.Now().UnixMilli(),
			}
			tradeResponse = append(tradeResponse, trade)
		}
		delete(currentOrderBook.OrderIDMap, newOrder.ID)
	}
	return tradeResponse, nil
}