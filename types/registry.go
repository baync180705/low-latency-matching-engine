package dataType

import (
	"sync"
)

type Regsitry struct {
	Books map[string]*OrderBook
	Mu sync.RWMutex
}

type OrderBook struct {
	BuyPriorityQueue  *PriorityQueue
    SellPriorityQueueSide *PriorityQueue 
    OrderIDMap map[string]*Order 
    Mu sync.RWMutex
}

type PriorityQueue struct {
	PriceHeap any
	TimeQueue map[int64]*OrderList
}