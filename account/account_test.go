package account

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/go-kit/kit/log"
)

type FakeRepo struct {
	makeError bool
}

func (f *FakeRepo) Store(*Account) error {
	if f.makeError {
		return errors.New("test error")
	}
	return nil
}
func (f *FakeRepo) Find(id string) (*Account, error) {
	if f.makeError {
		return nil, errors.New("test error")
	}
	return &Account{ID: id}, nil

}
func (f *FakeRepo) FindAll() []*Account {
	return []*Account{}
}

func TestAccountModel(t *testing.T) {

	acc := New()

	if acc.ID == "" {
		t.Error("Account did not generate ID")
		return
	}

	amount := acc.BalanceFor("USD")
	if amount > 0 {
		t.Errorf("Account balance map not initialized with zero values, got: %v", amount)
		return
	}

	acc.SetBalance("USD", 10)
	amount = acc.BalanceFor("USD")
	if amount != 10 {
		t.Errorf("expected %v got %v", 10, amount)
		return
	}

	if !acc.HasFunds("USD", 10) {
		t.Error("Account does not have required balance")
		return
	}

	acc.AppendBalance("USD", 1)
	amount = acc.BalanceFor("USD")
	if amount != 11 {
		t.Errorf("expected %v got %v", 11, amount)
	}

}

func TestAccountService(t *testing.T) {
	fr := &FakeRepo{}
	service := NewService(fr)

	account, err := service.CreateAccount()
	if err != nil {
		t.Error("Service cannot create account ", err)
	}
	if account == nil {
		t.Error("Service did not create an account, got nil")
	}

	account, err = service.GetAccount("123")
	if err != nil {
		t.Error("Service cannot get account ", err)
	}
	if account == nil {
		t.Error("Service did not get an account, got nil")
	}

	account, err = service.SetBalanceForAccount(&Account{ID: "123"}, Currency("USD"), 1)
	if err != nil {
		t.Error("Service cannot set balance for account ", err)
	}
	if account == nil {
		t.Error("Service did not set an balance, got nil")
	}

	// lets handle errors
	fr.makeError = true
	err = nil
	_, err = service.CreateAccount()
	if err == nil {
		t.Error("Service should yield error for creation an account, got nil ")
	}

	err = nil
	_, err = service.GetAccount("123")
	if err == nil {
		t.Error("Service should yield error for getting an account, got nil")
	}

	err = nil
	_, err = service.SetBalanceForAccount(&Account{}, Currency("USD"), 1)
	if err == nil {
		t.Error("Service should yield error for setting a balance, got nil")
	}

}

func TestAccountREST(t *testing.T) {
	fr := &FakeRepo{}
	service := NewService(fr)

	var logger = log.NewLogfmtLogger(os.Stderr)
	handler := MakeHandler(service, logger)

	rr := makeRequest(t, "GET", "/accounts", handler)

	rr = makeRequest(t, "POST", "/accounts", handler)
	var res accountsResponse
	if err := json.NewDecoder(rr.Body).Decode(&res); err != nil {
		t.Errorf("error decoding json %v", err)
		return
	}
	if res.Error != "" {
		t.Errorf("response with error %v", res.Error)
		return
	}
	t.Log(res)

	// test with errors
	fr.makeError = true

	rr = makeRequest(t, "POST", "/accounts", handler)
	res = accountsResponse{}
	if err := json.NewDecoder(rr.Body).Decode(&res); err != nil {
		t.Error(err)
		return
	}
	if res.Error == "" {
		t.Error("expected error, but got nil")
		return
	}

	listRes := listAccountsResponse{}
	rr = makeRequest(t, "GET", "/accounts", handler)
	if err := json.NewDecoder(rr.Body).Decode(&listRes); err != nil {
		t.Error(err)
		return
	}
	if len(listRes.Accounts) > 0 {
		t.Errorf("expected list to be empty, got %v", len(listRes.Accounts))
	}

}

func makeRequest(t *testing.T, method string, path string, handler http.Handler) *httptest.ResponseRecorder {
	req, err := http.NewRequest(method, path, nil)
	if err != nil {
		t.Error(err)
		t.Fail()
		return nil
	}
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Errorf("error from http, expected 200 got %v", rr.Code)
		t.Fail()
		return nil
	}
	return rr
}
