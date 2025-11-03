package algo

import (
	"container/heap"
	"sync"
	"time"

	"github.com/baync180705/low-latency-matching-engine/algo/handlers"
	types "github.com/baync180705/low-latency-matching-engine/types"
	"github.com/google/uuid"
)

var globalRegistry *types.Regsitry
var once sync.Once

func GetRegistry() *types.Regsitry {
	once.Do(func() {
		globalRegistry = types.NewRegistry()
	})
	return globalRegistry
}

func SubmitOrderEntry(order *types.OrderInput) (*types.Order, error) {
	// First step is to validate whether the input matches the correct standards eg: quantity cannot be negative etc.
	if err := handlers.ValidateInput(order); err != nil {
		return nil, err
	}
	
	globalRegistry := GetRegistry()
	registry := globalRegistry
	
	newOrder := types.Order{
		ID:          uuid.New().String(),
		Symbol:      order.Symbol,
		Side:        order.Side,
		Type:        order.Type,
		Price:       order.Price,
		Quantity:    order.Quantity,
		Timestamp:   time.Now().UnixMilli(),
		IsComplete:  false,
		IsCancelled: false,
		InitQty:     order.Quantity,
	}
	
	registry.Mu.Lock()
	book, exists := registry.Books[newOrder.Symbol]
	// Check if a field corresponding to the given value exists in the map, if it does not, then initailize all the maps and heaps
	if !exists {
		book = types.NewOrderBook(newOrder.Symbol)
		registry.Books[newOrder.Symbol] = book
	}
	registry.Mu.Unlock()
	
	book.Mu.Lock()
	defer book.Mu.Unlock()
	
	// Always add to OrderIDMap for tracking
	book.OrderIDMap[newOrder.ID] = &newOrder
	
	// ONLY add LIMIT orders to the book heaps
	// Market orders execute immediately or get rejected
	if newOrder.Type == "LIMIT" {
		if newOrder.Side == "BUY" {
			if !book.BuyHeap.PriceLevelExists(newOrder.Price) {
				heap.Push(book.BuyHeap, newOrder.Price)
			}
			_, listExists := book.BuyHeap.TimeQueue[newOrder.Price]
			if !listExists {
				orderList := types.NewOrderList()
				book.BuyHeap.TimeQueue[newOrder.Price] = orderList
			}
			book.BuyHeap.TimeQueue[newOrder.Price].PushBack(&newOrder)
			book.BuyHeap.Qty += newOrder.Quantity
		} else if newOrder.Side == "SELL" {
			if !book.SellHeap.PriceLevelExists(newOrder.Price) {
				heap.Push(book.SellHeap, newOrder.Price)
			}
			_, listExists := book.SellHeap.TimeQueue[newOrder.Price]
			if !listExists {
				orderList := types.NewOrderList()
				book.SellHeap.TimeQueue[newOrder.Price] = orderList
			}
			book.SellHeap.TimeQueue[newOrder.Price].PushBack(&newOrder)
			book.SellHeap.Qty += newOrder.Quantity
		}
	}
	
	return &newOrder, nil
}