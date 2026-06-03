package data

import (
	"errors"
	"time"
)

// Define clear domain errors
var (
	ErrEntityNotFound    = errors.New("entity not found")
	ErrBelowThreshold    = errors.New("withdrawal amount is below the minimum threshold")
	ErrInsufficientFunds = errors.New("insufficient funds in entity wallet")
	ErrInvalidAmount     = errors.New("amount must be greater than zero")
)

// An Entity represents a sports club with a centralized wallet
type Entity struct {
	ID            string  `json:"id"`
	Name          string  `json:"name"`
	WalletBalance float64 `json:"wallet_balance"`
}

// A Transaction logs tip history
type Transaction struct {
	ID        string    `json:"id"`
	EntityID  string    `json:"entity_id"`
	Amount    float64   `json:"amount"`
	Type      string    `json:"type"` // "tip" or "withdrawal"
	Timestamp time.Time `json:"timestamp"`
}
