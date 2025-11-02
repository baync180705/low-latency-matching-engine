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

func NewRegistry() *Regsitry {
    return &Regsitry{
        Books: make(map[string]*OrderBook), 
    }
}

func NewOrderBook(symbol string) *OrderBook {
    return &OrderBook{
        BuyHeap: NewHeap(true),
        SellHeap: NewHeap(false),
        OrderIDMap: make(map[string]*Order),
    }
}