package transaction

import (
	"errors"

	"github.com/MarinX/kit-payment/account"
	"github.com/go-kit/kit/log"
)

// Service is the interface that provides transaction methods.
type Service interface {
	// CreateTransaction creates a raw transaction
	CreateTransaction(string, string, account.Currency, float64) (*Transaction, error)

	// CommitTransaction commits the transaction by ID
	CommitTransaction(string) (*Transaction, error)

	// Transactions lists all transactions
	Transactions() []*Transaction

	// GetTransaction returns transaction by IDD
	GetTransaction(string) (*Transaction, error)

	// Watch is a event for transaction update
	Watch()
}

type service struct {
	transactions Repository
	accounts     account.Repository
	onCreate     chan *Transaction
	onPending    chan *Transaction
	log          log.Logger
}

// NewService creates transaction service
func NewService(transactions Repository, accounts account.Repository, log log.Logger) Service {
	return &service{
		transactions: transactions,
		accounts:     accounts,
		log:          log,
		onCreate:     make(chan *Transaction, 250),
		onPending:    make(chan *Transaction, 250),
	}
}

func (s *service) CreateTransaction(from string, to string, currency account.Currency, amount float64) (*Transaction, error) {

	if _, err := s.accounts.Find(from); err != nil {
		return nil, err
	}

	if _, err := s.accounts.Find(to); err != nil {
		return nil, err
	}

	tx := New(from, to, currency, amount)
	tx.Create()
	err := s.transactions.Store(tx)
	s.onCreate <- tx
	return tx, err
}

func (s *service) CommitTransaction(id string) (*Transaction, error) {
	tx, err := s.transactions.Find(id)
	if err != nil {
		return nil, err
	}
	if tx.Status != StatusCreated {
		return nil, errors.New("unknown transaction")
	}
	err = tx.Commit()
	if err != nil {
		return nil, err
	}
	err = s.transactions.Store(tx)
	s.onPending <- tx
	return tx, err
}

func (s *service) Transactions() []*Transaction {
	return s.transactions.FindAll()
}

func (s *service) GetTransaction(id string) (*Transaction, error) {
	return s.transactions.Find(id)
}

func (s *service) Watch() {
	for {
		select {
		case <-s.onCreate:
			// we can notify 3rd party systems here for new created transaction
			break
		case tx := <-s.onPending:
			from, err := s.accounts.Find(tx.From)
			if err != nil {
				s.log.Log("cannot find account", "from", tx.From)
				return
			}
			to, err := s.accounts.Find(tx.To)
			if err != nil {
				s.log.Log("cannot find account", "to", tx.To)
			}
			if !from.HasFunds(tx.Currency, tx.Amount) {
				tx.Status = StatusInsufficientFunds
				err := s.transactions.Store(tx)
				s.checkError(err)
				return
			}

			// everything is fine, transfer the money
			from.AppendBalance(tx.Currency, -tx.Amount)
			err = s.accounts.Store(from)
			s.checkError(err)

			to.AppendBalance(tx.Currency, tx.Amount)
			err = s.accounts.Store(to)
			s.checkError(err)

			tx.Status = StatusOK
			err = s.transactions.Store(tx)
			s.checkError(err)

			break
		}
	}
}

func (s *service) checkError(err error) {
	if err != nil {
		s.log.Log("payment", "error", err)
	}
}
