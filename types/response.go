package dataType

type OrderBookLevel struct {
	Price    int64 `json:"price"`
	Quantity int64 `json:"quantity"`
}

type OrderBookResponse struct {
	Symbol    string            `json:"symbol"`
	Timestamp int64             `json:"timestamp"`
	Bids      []OrderBookLevel  `json:"bids"`
	Asks      []OrderBookLevel  `json:"asks"`
}

type OrderResponse struct {
	OrderID          string         `json:"order_id,omitempty"`
	Status           string         `json:"status,omitempty"`
	Message          string         `json:"message,omitempty"`
	FilledQuantity   int64          `json:"filled_quantity,omitempty"`
	RemainingQuantity int64         `json:"remaining_quantity,omitempty"`
	Trades           []*TradeRecord `json:"trades,omitempty"`
	Error            string         `json:"error,omitempty"`
}

type TradeRecord struct {
    TradeID   string `json:"trade_id"`
    Price     int64  `json:"price"`
    Quantity  int64  `json:"quantity"`
    Timestamp int64  `json:"timestamp"`
}

type StatusResponse struct {
	OrderID        string `json:"order_id"`
	Symbol         string `json:"symbol"`
	Side           string `json:"side"`
	Type           string `json:"type"`
	Price          int64  `json:"price"`
	Quantity       int64  `json:"quantity"`        
	FilledQuantity int64  `json:"filled_quantity"`  
	Status         string `json:"status"`           
	Timestamp      int64  `json:"timestamp"`
}

