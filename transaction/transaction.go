package transaction

import (
	"crypto/sha256"
	"encoding/hex"

	"github.com/MarinX/kit-payment/account"
	"github.com/cbergoon/merkletree"
	uuid "github.com/satori/go.uuid"
)

// TransactionStatus is our status handler
type TransactionStatus string

const (
	// StatusOK if the transaction passes verifications
	StatusOK TransactionStatus = "ok"

	// StatusInsufficientFunds if the account owner does not have enough balance
	StatusInsufficientFunds TransactionStatus = "insufficient_funds"

	// StatusErr for unknown errors that can happen
	StatusErr TransactionStatus = "unknown_error"

	// StatusPending is when transaction is ready to be processed
	StatusPending TransactionStatus = "pending"

	// StatusCreated is when transaction is created and ready for commit
	StatusCreated TransactionStatus = "created"
)

// Transaction represents transaction between 2 accounts
type Transaction struct {
	ID       string            `json:"id"`
	From     string            `json:"from"`
	To       string            `json:"to"`
	Status   TransactionStatus `json:"status"`
	Amount   float64           `json:"amount"`
	Currency account.Currency  `json:"currency"`
}

// Repository provides access a transaction store.
type Repository interface {
	Store(*Transaction) error
	Find(id string) (*Transaction, error)
	FindAll() []*Transaction
	Delete(string) error
}

// New creates transaction between 2 accounts
func New(from string, to string, currency account.Currency, amount float64) *Transaction {
	return &Transaction{
		From:     from,
		To:       to,
		Currency: currency,
		Amount:   amount,
	}
}

// Create creates new transaction with generated ID
func (t *Transaction) Create() {
	t.ID = uuid.Must(uuid.NewV4()).String()
	t.Status = StatusCreated
}

// Commit commits the transaction, ready to be processed
func (t *Transaction) Commit() error {
	t.Status = StatusPending
	return nil
}

//CalculateHash hashes the values of a transaction ID
func (t Transaction) CalculateHash() ([]byte, error) {
	h := sha256.New()
	if _, err := h.Write([]byte(t.ID)); err != nil {
		return nil, err
	}

	return h.Sum(nil), nil
}

// Hash is string representation of calculated hash
func (t *Transaction) Hash() (string, error) {
	calculated, err := t.CalculateHash()
	return hex.EncodeToString(calculated), err
}

// Equals checks if the content of transaction is equal to another content transaction
func (t Transaction) Equals(other merkletree.Content) (bool, error) {
	return t.ID == other.(Transaction).ID, nil
}
