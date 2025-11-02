package dataType

import (
	"sync"
)

type Regsitry struct {
	Books map[string]*OrderBook
	Mu sync.RWMutex
}

type OrderBook struct {
	BuyHeap  *Heap
    SellHeap *Heap 
    OrderIDMap map[string]*Order 
    Mu sync.RWMutex
}

func NewOrderBook(symbol string) *OrderBook {
    return &OrderBook{
        BuyHeap: NewHeap(true),
        SellHeap: NewHeap(false),
        OrderIDMap: make(map[string]*Order),
    }
}