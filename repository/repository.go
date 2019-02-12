package repository

import (
	"time"

	"github.com/boltdb/bolt"
)

// Repository is our main holder for database access
type Repository struct {
	db *bolt.DB
}

// New creates new boltdb database
func New() (*Repository, error) {
	db, err := bolt.Open("data.db", 0600, &bolt.Options{Timeout: 5 * time.Second})
	return &Repository{db: db}, err
}

// Account returns account repository
func (r *Repository) Account() *accountRepository {
	return &accountRepository{db: r.db}
}

// Transaction returns transaction repository
func (r *Repository) Transaction() *transactionRepository {
	return &transactionRepository{db: r.db}
}

// Close the database
func (r *Repository) Close() error {
	return r.db.Close()
}
