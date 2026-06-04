package data

import (
	"context"
	"crypto/rand"
	"database/sql"
	"fmt"
	"time"
)

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{db: db}
}

func (s *Store) ListEntities() ([]Entity, error) {
	query := `SELECT id, name, wallet_balance FROM entities`

	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var entities []Entity
	for rows.Next() {
		var e Entity
		if err := rows.Scan(&e.ID, &e.Name, &e.WalletBalance); err != nil {
			return nil, err
		}
		entities = append(entities, e)
	}

	return entities, nil
}

func GetEntity(s *Store, id string) *Entity {
	query := `SELECT id, name, wallet_balance FROM entities WHERE id = ?`
	row := s.db.QueryRow(query, id)
	var e Entity
	if err := row.Scan(&e.ID, &e.Name, &e.WalletBalance); err != nil {
		return nil
	}
	return &e
}

// ProcessTip simulates receiving funds from an external mobile wallet
func ProcessTip(s *Store, entityID string, amount float64) (*Transaction, error) {
	if amount <= 0 {
		return nil, ErrInvalidAmount
	}
	// Use a background context for execution control
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 1. Begin the ACID Transaction
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	// Ensure the transaction is rolled back if any error occurs
	defer tx.Rollback()

	// 2. Lock the entity's row for update
	var exists bool
	err = tx.QueryRowContext(ctx, `SELECT EXISTS(SELECT 1 FROM entities WHERE id = ?)`, entityID).Scan(&exists)
	if err != nil || !exists {
		return nil, ErrEntityNotFound
	}

	// 3. Update the entity's wallet balance
	_, err = tx.ExecContext(ctx, `UPDATE entities SET wallet_balance = wallet_balance + ? WHERE id = ?`, amount, entityID)
	if err != nil {
		return nil, err
	}

	// 4. Log the transaction
	newTx := Transaction{
		ID:        generateID(),
		EntityID:  entityID,
		Amount:    amount,
		Type:      "tip",
		Timestamp: time.Now(),
	}
	_, err = tx.ExecContext(ctx, `INSERT INTO transactions (id, entity_id, amount, type, timestamp) VALUES (?, ?, ?, ?, ?)`,
		newTx.ID, newTx.EntityID, newTx.Amount, newTx.Type, newTx.Timestamp)
	if err != nil {
		return nil, err
	}

	// 5. Commit the transaction
	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return &newTx, nil
}

// ProcessWithdrawal simulates sending funds to an external mobile wallet
func ProcessWithdrawal(s *Store, entityID string, amount float64) (*Transaction, error) {
	const minThreshold = 1000.0 // Entities can only withdraw if amount >= 1000F CFA

	if amount < minThreshold {
		return nil, ErrBelowThreshold
	}
	// Use a background context for execution control
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	// 1. Begin the ACID Transaction
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	// Ensure the transaction is rolled back if any error occurs
	defer tx.Rollback()

	var currentBalance float64
	err = tx.QueryRowContext(ctx, `SELECT wallet_balance FROM entities WHERE id = ? FOR UPDATE`, entityID).Scan(&currentBalance)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrEntityNotFound
		}
		return nil, err
	}

	if currentBalance < amount {
		return nil, ErrInsufficientFunds
	}

	// 2. Update the entity's wallet balance
	_, err = tx.ExecContext(ctx, `UPDATE entities SET wallet_balance = wallet_balance - ? WHERE id = ?`, amount, entityID)
	if err != nil {
		return nil, err
	}

	// 3. Log the transaction
	newTx := Transaction{
		ID:        generateID(),
		EntityID:  entityID,
		Amount:    amount,
		Type:      "withdrawal",
		Timestamp: time.Now(),
	}

	_, err = tx.ExecContext(ctx, `INSERT INTO transactions (id, entity_id, amount, type, timestamp) VALUES (?, ?, ?, ?, ?)`,
		newTx.ID, newTx.EntityID, newTx.Amount, newTx.Type, newTx.Timestamp)
	if err != nil {
		return nil, err
	}

	// 4. Commit the transaction
	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return &newTx, nil
}

// Basic cryptographically secure ID generator helper
func generateID() string {
	b := make([]byte, 4)
	_, _ = rand.Read(b)
	return fmt.Sprintf("tx_%x", b)
}
