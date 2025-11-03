package algo

import (
	"errors"
	"time"

	types "github.com/baync180705/low-latency-matching-engine/types"
	"github.com/google/uuid"
)

func RunPipeline(input *types.OrderInput) (*types.Order, []*types.TradeRecord, error) {
	order, err := SubmitOrderEntry(input)

	if err!= nil {
		return order, nil, errors.New("Failed to submit your order request")
	}

	globalRegistry:= GetRegistry()
	var heap *types.Heap 
	if order.Side=="BUY" {
		heap = globalRegistry.Books[order.Symbol].BuyHeap
	} else if order.Side =="SELL" {
		heap = globalRegistry.Books[order.Symbol].SellHeap
	}

	if heap.Len() == 0 {
		return order, nil, errors.New("Order queued — no active price levels yet")
	}

	price := heap.PriceHeap[0]
	entry := heap.TimeQueue[price].Front().Value.(*types.Order)

	var tradeResponse []*types.TradeRecord

	defaultRecd := &types.TradeRecord{
		TradeID: uuid.New().String() ,
		Price: 0,
		Quantity: 0,
		Timestamp: time.Now().UnixMilli(),
	}

	if order.Type == "LIMIT" {
		if entry.ID ==order.ID {
			recd, err := MatchingAlgorithm(order)
			if err!=nil {
				tradeResponse = append(tradeResponse, defaultRecd)
				return order, tradeResponse, err
			}
			tradeResponse = recd
		} else {
			tradeResponse = append(tradeResponse, defaultRecd)
		}
	} else if order.Type =="MARKET" {
		recd, err := MatchingAlgorithm(order)
		if err!=nil {
			tradeResponse = append(tradeResponse, defaultRecd)
			return order, tradeResponse, err
		}
		tradeResponse = recd
	}
	
	return order, tradeResponse, nil
}