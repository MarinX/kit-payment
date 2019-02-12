package account

import (
	"context"
	"errors"
	"strings"

	"github.com/go-kit/kit/endpoint"
)

type accountsRequest struct{}

type accountsResponse struct {
	Account *Account `json:"account"`
	Error   string   `json:"error,omitempty"`
}

func makeAccountsEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		res := accountsResponse{}
		account, err := s.CreateAccount()
		if err != nil {
			res.Error = err.Error()
		}
		res.Account = account
		return res, nil
	}
}

type listAccountsRequest struct{}

type listAccountsResponse struct {
	Accounts []*Account `json:"accounts"`
}

func makeListAccountsEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		return listAccountsResponse{Accounts: s.Accounts()}, nil
	}
}

type accountsBalanceRequest struct {
	Currency  string  `json:"currency"`
	Amount    float64 `json:"amount"`
	AccountID string  `json:"-"`
}

type accountsBalanceResponse struct {
	Account *Account `json:"account"`
	Error   string   `json:"error,omitempty"`
}

func makeAccountsBalanceEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(accountsBalanceRequest)
		res := accountsBalanceResponse{}

		if req.Amount <= 0 {
			res.Error = errors.New("invalid balance set").Error()
			return res, nil
		}
		if req.Currency == "" {
			res.Error = errors.New("invalid currency set").Error()
			return res, nil
		}

		account, err := s.GetAccount(req.AccountID)
		if err != nil {
			res.Error = err.Error()
			return res, nil
		}

		account, err = s.SetBalanceForAccount(account, Currency(strings.ToUpper(req.Currency)), req.Amount)
		if err != nil {
			res.Error = err.Error()
		}
		res.Account = account
		return res, nil
	}
}
