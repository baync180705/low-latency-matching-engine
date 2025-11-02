package dataType

import "container/heap"

type PriceLevel struct {
	Price int64 
	IsBuy bool 
}

type PriceHeap []*PriceLevel

func (ph PriceHeap) Len() int {
	return len(ph)
}

func (ph PriceHeap) Swap(i, j int) {
	ph[i], ph[j] = ph[j], ph[i]
}
func (ph PriceHeap) Less(i, j int) bool {
	isBuy := ph[i].IsBuy
	if isBuy {
		return ph[i].Price > ph[j].Price 
	} else {
		return ph[i].Price < ph[j].Price
	}
}

func (ph *PriceHeap) Push(x any) {
	*ph = append(*ph, x.(*PriceLevel))
}

func (ph *PriceHeap) Pop() any {
	old := *ph
	n := len(old)
	x := old[n-1]
	*ph = old[0 : n-1]
	return x
}

func (ph *PriceHeap) Remove(i int) any {
    return heap.Remove(ph, i)
}