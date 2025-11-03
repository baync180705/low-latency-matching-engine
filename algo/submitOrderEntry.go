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

func GetRegistry() *types.Regsitry {
    once.Do(func() {
        globalRegistry = types.NewRegistry()
    })
    return globalRegistry
} // I could have simply exprted globalRegistry by making it GlobalRegistry, but there could be an edge case when the GlobalRegistry gets used in the other files before being initialzied properly. The GetRegistry function ensures that the registry is always initialized before it gets exported.

func SubmitOrderEntry(order *types.OrderInput) (*types.Order, error) {
	// First step is to validate whether the input matches the correct standards eg: quantity cannot be negative etc.
	if err := handlers.ValidateInput(order); err !=nil {
		return nil, err
	}

	globalRegistry := GetRegistry() // This will get executed only once in the program lifecyle. This is because the globalRegistry has to be initialized the 1st time the program executes.

	registry := globalRegistry //Earlier I tried to create a copy by value, but the globalRegistry struct contains sync.RWMutex field and in go copying any struct which contains RWMutex field is not permitted.

	newOrder := types.Order{
        ID:        uuid.New().String(),
        Symbol:    order.Symbol,
        Side:      order.Side,
        Type:      order.Type,
        Price:     order.Price,
        Quantity:  order.Quantity,
        Timestamp: time.Now().UnixMilli(),
    } //Here I added ID and Timestamp fields so that the newOrder confirms to the Order dataType.

	registry.Mu.Lock()
	book, exists := registry.Books[newOrder.Symbol]
	
	//Check if a field corresponding to the given value exists in the map, if it does not, then initailize all the maps and heaps. If they are not initialized we cannot perform push, pop or lookup operations since all they will host is a nil value
	if !exists {
		book = types.NewOrderBook(newOrder.Symbol)
		registry.Books[newOrder.Symbol] = book
	}

	registry.Mu.Unlock()

	book.Mu.Lock()
	defer book.Mu.Unlock() //We will keep modifying the book throughout the function, hence it will be unlocked at the end of the function

	book.OrderIDMap[newOrder.ID] = &newOrder

	//Remaining code populates the fields based on whether the order is from the BUY side or SELL side.
	if newOrder.Side=="BUY" {
		if !book.BuyHeap.PriceLevelExists(newOrder.Price) {book.BuyHeap.Push(newOrder.Price)}

		_, listExists := book.BuyHeap.TimeQueue[newOrder.Price]
		if !listExists {
			orderList := types.NewOrderList()
			book.BuyHeap.TimeQueue[newOrder.Price] = orderList 
		}

		book.BuyHeap.TimeQueue[newOrder.Price].PushBack(&newOrder)
		book.BuyHeap.Qty+=newOrder.Quantity
	} else if newOrder.Side=="SELL" {
		if !book.SellHeap.PriceLevelExists(newOrder.Price) {book.SellHeap.Push(newOrder.Price)}

		_, listExists := book.SellHeap.TimeQueue[newOrder.Price]
		if !listExists {
			orderList := types.NewOrderList()
			book.SellHeap.TimeQueue[newOrder.Price] = orderList 
		}

		book.SellHeap.TimeQueue[newOrder.Price].PushBack(&newOrder)
		book.SellHeap.Qty+=newOrder.Quantity
	}
	
	return &newOrder, nil
} 