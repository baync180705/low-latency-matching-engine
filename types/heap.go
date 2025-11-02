package dataType

type Heap struct {
	PriceHeap []int64 //Slice of Price
	TimeQueue map[int64]*OrderList
	isBuy bool
}

func (h Heap) Len() int {
	return len(h.PriceHeap)
}

func (h Heap) Swap(i, j int) {
	h.PriceHeap[i], h.PriceHeap[j] = h.PriceHeap[j], h.PriceHeap[i]
}
func (h Heap) Less(i, j int) bool {
	isBuy := h.isBuy
	if isBuy {
		return h.PriceHeap[i] > h.PriceHeap[j]
	} else {
		return h.PriceHeap[i] < h.PriceHeap[j]
	}
}

func (h *Heap) Push(x any) {
    h.PriceHeap = append(h.PriceHeap, x.(int64))
}

func (h *Heap) Pop() any {
    old := h.PriceHeap
    n := len(old)
    x := old[n-1]
    h.PriceHeap = old[0 : n-1]
    return x
}

func (h *Heap) Remove(i int) any {
    if i < 0 || i >= len(h.PriceHeap) {
        return nil
    }
    removed := h.PriceHeap[i]
    h.PriceHeap = append(h.PriceHeap[:i], h.PriceHeap[i+1:]...)
    return removed
}

func NewHeap(isBuy bool) *Heap {
    newHeap := &Heap{
        PriceHeap: []int64{},
        TimeQueue: make(map[int64]*OrderList),
        isBuy: isBuy,
    }
    return newHeap
}