package account

import (
	uuid "github.com/satori/go.uuid"
)

// Currency represents the key for global currency
type Currency string

// Account represents id holding multiple currencies
type Account struct {
	ID       string               `json:"id"`
	Balances map[Currency]float64 `json:"balances,omitempty"`
}

// Repository provides access a account store.
type Repository interface {
	Store(*Account) error
	Find(id string) (*Account, error)
	FindAll() []*Account
}

// New creates account with id
func New() *Account {
	return &Account{
		ID:       uuid.Must(uuid.NewV4()).String(),
		Balances: make(map[Currency]float64),
	}
}

// BalanceFor returns amount for given currency
func (a *Account) BalanceFor(currency Currency) float64 {
	return a.Balances[currency]
}

// SetBalance hard reset balance for given currency
func (a *Account) SetBalance(currency Currency, amount float64) {
	if a.Balances == nil {
		a.Balances = make(map[Currency]float64)
	}
	a.Balances[currency] = amount
}

// AppendBalance adds or removes from account balance for given currency
func (a *Account) AppendBalance(currency Currency, amount float64) {
	if a.Balances == nil {
		a.Balances = make(map[Currency]float64)
	}
	a.Balances[currency] += amount
}

// HasFunds checks if the account has enough amount for given currency
func (a *Account) HasFunds(currency Currency, amount float64) bool {
	if a.Balances == nil {
		a.Balances = make(map[Currency]float64)
	}
	return a.Balances[currency] >= amount
}
