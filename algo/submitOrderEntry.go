package algo

import (
	"time"
	"sync"

	"github.com/baync180705/low-latency-matching-engine/algo/handlers"
	types "github.com/baync180705/low-latency-matching-engine/types"
	"github.com/google/uuid"
)

var globalRegistry *types.Regsitry
var once sync.Once


func SubmitOrderEntry(order *types.OrderInput) error {
	if err := handlers.ValidateInput(order); err !=nil {
		return err
	}

	once.Do(func() {
		globalRegistry = types.NewRegistry()
	}) // This will get executed only once in the program lifecyle. This is because the globalRegistry has to be initialized the 1st time the program executes.

	registry := globalRegistry //Earlier I tried to create a copy by value, but the globalRegistry struct contains sync.RWMutex field and in go copying any struct which contains RWMutex field is not permitted.

	newOrder := types.Order{
        ID:        uuid.New().String(),
        Symbol:    order.Symbol,
        Side:      order.Side,
        Type:      order.Type,
        Price:     order.Price,
        Quantity:  order.Quantity,
        Timestamp: time.Now().UnixMilli(),
    }

	registry.Mu.Lock()
	book, exists := registry.Books[newOrder.Symbol]

	if !exists {
		book = types.NewOrderBook(newOrder.Symbol)
		registry.Books[newOrder.Symbol] = book
	}

	registry.Mu.Unlock()

	book.Mu.Lock()
	defer book.Mu.Unlock() //We will keep modifying the book, hence it will be unlocked at the end of the function

	book.OrderIDMap[newOrder.ID] = &newOrder
	if newOrder.Side=="BUY" {
		if !book.BuyHeap.PriceLevelExists(newOrder.Price) {book.BuyHeap.Push(newOrder.Price)}

		_, listExists := book.BuyHeap.TimeQueue[newOrder.Price]
		if !listExists {
			orderList := types.NewOrderList()
			book.BuyHeap.TimeQueue[newOrder.Price] = orderList 
		}

		book.BuyHeap.TimeQueue[newOrder.Price].PushBack(&newOrder)
	} else {
		if !book.SellHeap.PriceLevelExists(newOrder.Price) {book.SellHeap.Push(newOrder.Price)}

		_, listExists := book.SellHeap.TimeQueue[newOrder.Price]
		if !listExists {
			orderList := types.NewOrderList()
			book.SellHeap.TimeQueue[newOrder.Price] = orderList 
		}

		book.SellHeap.TimeQueue[newOrder.Price].PushBack(&newOrder)
	}
	
	return nil
} 