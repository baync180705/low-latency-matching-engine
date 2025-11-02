package dataType

import (
	"sync"
)

type Regsitry struct {
	Books map[string]*OrderBook
	Mu sync.RWMutex
}

type OrderBook struct {
	BuyPriorityQueue  *Heap
    SellPriorityQueueSide *Heap 
    OrderIDMap map[string]*Order 
    Mu sync.RWMutex
}

type Heap struct {
	PriceHeap *PriceHeap
	TimeQueue map[int64]*OrderList
}