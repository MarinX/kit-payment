package transaction

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/MarinX/kit-payment/account"
	"github.com/go-kit/kit/log"
)

type FakeRepoAccount struct {
	makeError bool
}

func (f *FakeRepoAccount) Store(*account.Account) error {
	if f.makeError {
		return errors.New("test error")
	}
	return nil
}
func (f *FakeRepoAccount) Find(id string) (*account.Account, error) {
	if f.makeError {
		return nil, errors.New("test error")
	}
	return &account.Account{ID: id}, nil

}
func (f *FakeRepoAccount) FindAll() []*account.Account {
	return []*account.Account{}
}

type FakeRepoTransaction struct {
	makeError bool
}

func (f *FakeRepoTransaction) Store(*Transaction) error {
	if f.makeError {
		return errors.New("test error")
	}
	return nil
}
func (f *FakeRepoTransaction) Find(id string) (*Transaction, error) {
	if f.makeError {
		return nil, errors.New("test error")
	}
	return &Transaction{ID: id, Status: StatusCreated}, nil
}
func (f *FakeRepoTransaction) FindAll() []*Transaction {
	return []*Transaction{}
}
func (f *FakeRepoTransaction) Delete(id string) error {
	if f.makeError {
		return errors.New("test error")
	}
	return nil
}

func TestTransactionModel(t *testing.T) {

	tx := New("123", "222", account.Currency("USD"), 10)
	if tx.ID != "" {
		t.Error("expectedd ID to be empty on init model")
		return
	}
	tx.Create()
	if tx.ID == "" {
		t.Error("expected ID to be populated after creation")
		return
	}

	if tx.Status != StatusCreated {
		t.Errorf("transaction status is wrong, want %v got %v", StatusCreated, tx.Status)
		return
	}

	// merkle tree hash
	h, err := tx.Hash()
	if err != nil {
		t.Errorf("error creating transaction hash %v", err)
		return
	}
	t.Logf("merkle tree hash content %v", h)

	tx.Commit()
	if tx.Status != StatusPending {
		t.Errorf("transaction status is wrong, want %v got %v", StatusPending, tx.Status)
	}
}

func TestTransactionService(t *testing.T) {
	tfr := &FakeRepoTransaction{}
	afr := &FakeRepoAccount{}
	var logger = log.NewLogfmtLogger(os.Stderr)
	service := NewService(tfr, afr, logger)

	tx, err := service.CreateTransaction("123", "222", account.Currency("USD"), 1)
	if err != nil {
		t.Errorf("transaction creation error %v", err)
		return
	}
	if tx.ID == "" {
		t.Errorf("expected ID after transaction creation, got nil")
		return
	}

	tx, err = service.CommitTransaction(tx.ID)
	if err != nil {
		t.Errorf("error commit transaction %v", err)
		return
	}
	if tx.Status != StatusPending {
		t.Errorf("transaction did not change status, want %v got %v", StatusPending, tx.Status)
		return
	}

	if _, err := service.GetTransaction("123"); err != nil {
		t.Errorf("error getting transaction %v", err)
		return
	}

	tfr.makeError = true

	if _, err := service.CreateTransaction("123", "222", account.Currency("USD"), 1); err == nil {
		t.Error("expected error for creation, got nil")
		return
	}

	if _, err := service.CommitTransaction("123"); err == nil {
		t.Error("expected error for commit, got nil")
		return
	}

	if _, err := service.GetTransaction("123"); err == nil {
		t.Error("expected error for getting transaction, got nil")
		return
	}

	trxs := service.Transactions()
	if len(trxs) > 0 {
		t.Errorf("expected 0 transactions got %v", len(trxs))
	}

}

func TestTransactionREST(t *testing.T) {
	tfr := &FakeRepoTransaction{}
	afr := &FakeRepoAccount{}
	var logger = log.NewLogfmtLogger(os.Stderr)
	service := NewService(tfr, afr, logger)
	handler := MakeHandler(service, logger)

	rr := makeRequest(t, "POST", "/transactions", handler)
	res := transactionsResponse{}
	if err := json.NewDecoder(rr.Body).Decode(&res); err != nil {
		t.Error(err)
		return
	}
	if res.Error == "" {
		t.Error("expected error, got nil")
		return
	}

	rr = makeRequest(t, "GET", "/transactions", handler)
	listRes := listTransactionsResponse{}
	if err := json.NewDecoder(rr.Body).Decode(&listRes); err != nil {
		t.Error(err)
		return
	}

	rr = makeRequest(t, "GET", "/transactions/123", handler)
	res = transactionsResponse{}
	if err := json.NewDecoder(rr.Body).Decode(&res); err != nil {
		t.Error(err)
		return
	}
	if res.Error != "" {
		t.Errorf("unexpected error for getting transaction %v", res.Error)
		return
	}

	rr = makeRequest(t, "PUT", "/transactions/123/commit", handler)
	res = transactionsResponse{}
	if err := json.NewDecoder(rr.Body).Decode(&res); err != nil {
		t.Error(err)
		return
	}
	if res.Error != "" {
		t.Errorf("unexpected error for commiting transaction %v", res.Error)
		return
	}

	rr = makeRequest(t, "GET", "/transactions/123/hash", handler)
	hashRes := hashTransactionsResponse{}
	if err := json.NewDecoder(rr.Body).Decode(&hashRes); err != nil {
		t.Error(err)
		return
	}
	if res.Error != "" {
		t.Errorf("unexpected error for getting hash transaction %v", res.Error)
		return
	}
	t.Log(hashRes.Hash)

}

func makeRequest(t *testing.T, method string, path string, handler http.Handler) *httptest.ResponseRecorder {
	req, err := http.NewRequest(method, path, nil)
	if err != nil {
		t.Error(err)
		return nil
	}
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Errorf("error from http, expected 200 got %v", rr.Code)
		return nil
	}
	return rr
}
