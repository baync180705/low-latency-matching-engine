package dataType

import "container/list"

type Order struct {
    ID        string  `json:"id"`
    Symbol    string  `json:"symbol"`
    Side      string  `json:"side"`
    Type      string  `json:"type"`
    Price     int64   `json:"price,omitempty"`
    Quantity  int64   `json:"quantity"`
    Timestamp int64   `json:"timestamp"`
}

type OrderList struct {
    *list.List
} // Have declared OrderList as a doubly Linked List. I have done this because a doubly linked list stores both, the head and the tail pointer. This will enable us to levarage queue like property - FIFO in O(1)

type TradeRecord struct {
    TradeID   string `json:"trade_id"`
    Price     int64  `json:"price"`
    Quantity  int64  `json:"quantity"`
    Timestamp int64  `json:"timestamp"`
}
