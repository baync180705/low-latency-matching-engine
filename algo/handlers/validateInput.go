package handlers

import (
    "errors"
    "strings"
	types "github.com/baync180705/low-latency-matching-engine/types"
)

func ValidateInput(input types.OrderInput) error {
    if input.Quantity <= 0 {
        return errors.New("invalid order: quantity must be positive")
    }
    if strings.TrimSpace(input.Symbol) == "" {
        return errors.New("invalid order: symbol is required")
    }
    if input.Side != "BUY" && input.Side != "SELL" {
        return errors.New("invalid order: side must be 'BUY' or 'SELL'")
    }
    if input.Type != "LIMIT" && input.Type != "MARKET" {
        return errors.New("invalid order: type must be 'LIMIT' or 'MARKET'")
    }
    if input.Type == "LIMIT" {
        if input.Price <= 0 {
            return errors.New("invalid limit order: price must be positive for LIMIT orders")
        }
    }
    return nil
}