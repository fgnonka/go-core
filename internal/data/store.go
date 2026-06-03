package data

import (
	"crypto/rand"
	"fmt"
	"sync"
	"time"
)

type Store struct {
	mu           sync.RWMutex
	entities     map[string]*Entity
	transactions []Transaction
}

func NewStore() *Store {
	// Seed with dummy entities for our CRUD simulations
	s := &Store{
		entities: make(map[string]*Entity),
	}
	s.entities["asec"] = &Entity{ID: "asec", Name: "Asec Mimosas FC", WalletBalance: 0.0}
	s.entities["jca"] = &Entity{ID: "jca", Name: "JCA Kings", WalletBalance: 500.0} // Pre-funded to test withdrawal threshold
	return s
}

func ListEntities(s *Store) []*Entity {
	s.mu.RLock()
	defer s.mu.RUnlock()

	entities := make([]*Entity, 0, len(s.entities))
	for _, entity := range s.entities {
		entities = append(entities, entity)
	}
	return entities
}

func GetEntity(s *Store, id string) *Entity {
	s.mu.RLock()
	defer s.mu.RUnlock()
	entity, exists := s.entities[id]
	if !exists {
		return nil
	}
	return entity
}

// ProcessTip simulates receiving funds from an external mobile wallet
func ProcessTip(s *Store, entityID string, amount float64) (*Transaction, error) {
	if amount <= 0 {
		return nil, ErrInvalidAmount
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	entity, exists := s.entities[entityID]
	if !exists {
		return nil, ErrEntityNotFound
	}

	// Update the entity's wallet balance
	entity.WalletBalance += amount

	// Log the transaction
	tx := Transaction{
		ID:        generateID(),
		EntityID:  entityID,
		Amount:    amount,
		Type:      "tip",
		Timestamp: time.Now(),
	}
	s.transactions = append(s.transactions, tx)

	return &tx, nil
}

// ProcessWithdrawal simulates sending funds to an external mobile wallet
func ProcessWithdrawal(s *Store, entityID string, amount float64) (*Transaction, error) {
	const minThreshold = 1000.0 // Teams can only withdraw if amount >= 1000F CFA

	if amount < minThreshold {
		return nil, ErrBelowThreshold
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	entity, exists := s.entities[entityID]
	if !exists {
		return nil, ErrEntityNotFound
	}

	if entity.WalletBalance < amount {
		return nil, ErrInsufficientFunds
	}

	// Update the entity's wallet balance
	entity.WalletBalance -= amount

	// Log the transaction
	tx := Transaction{
		ID:        generateID(),
		EntityID:  entityID,
		Amount:    amount,
		Type:      "withdrawal",
		Timestamp: time.Now(),
	}
	s.transactions = append(s.transactions, tx)
	return &tx, nil
}

// Basic cryptographically secure ID generator helper
func generateID() string {
	b := make([]byte, 4)
	_, _ = rand.Read(b)
	return fmt.Sprintf("tx_%x", b)
}
