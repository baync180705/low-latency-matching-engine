package dataType

type OrderInput struct {
    Symbol   string `json:"symbol"`
    Side     string `json:"side"`    
    Type     string `json:"type"`    
    Price    int64  `json:"price,omitempty"` 
    Quantity int64  `json:"quantity"`
}