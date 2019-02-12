package transaction

import (
	"context"
	"errors"

	"github.com/MarinX/kit-payment/account"

	"github.com/go-kit/kit/endpoint"
)

type transactionsRequest struct {
	From     string           `json:"from"`
	To       string           `json:"to"`
	Amount   float64          `json:"amount"`
	Currency account.Currency `json:"currency"`
}

type transactionsResponse struct {
	Transaction *Transaction `json:"transaction"`
	Error       string       `json:"error,omitempty"`
}

func makeTransactionsEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(transactionsRequest)
		res := transactionsResponse{}
		if len(req.From) == 0 || len(req.To) == 0 {
			res.Error = errors.New("missing address").Error()
			return res, nil
		}

		if len(req.Currency) == 0 {
			res.Error = errors.New("missing currency").Error()
			return res, nil
		}

		if req.Amount <= 0 {
			res.Error = errors.New("invalid amount").Error()
			return res, nil
		}

		tx, err := s.CreateTransaction(req.From, req.To, req.Currency, req.Amount)
		if err != nil {
			res.Error = err.Error()
		}
		res.Transaction = tx
		return res, nil
	}
}

type listTransactionsRequest struct{}

type listTransactionsResponse struct {
	Transactions []*Transaction `json:"transactions"`
}

func makeListTransactionsEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		return listTransactionsResponse{
			Transactions: s.Transactions(),
		}, nil
	}
}

type getTransactionsRequest struct {
	ID string
}

func makeGetTransactionsEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(getTransactionsRequest)
		res := transactionsResponse{}

		if req.ID == "" {
			res.Error = errors.New("missing required ID").Error()
			return res, nil
		}
		trx, err := s.GetTransaction(req.ID)
		if err != nil {
			res.Error = err.Error()
		}
		res.Transaction = trx

		return res, nil
	}
}

type commitTransactionsRequest struct {
	ID string
}

func makeCommitTransactionsEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(commitTransactionsRequest)
		res := transactionsResponse{}
		if req.ID == "" {
			res.Error = errors.New("missing required ID").Error()
			return res, nil
		}

		trx, err := s.CommitTransaction(req.ID)
		if err != nil {
			res.Error = err.Error()
			return res, nil
		}
		res.Transaction = trx
		return res, nil
	}
}

type hashTransactionsRequest struct {
	ID string
}
type hashTransactionsResponse struct {
	Hash  string `json:"hash"`
	Error string `json:"error,omitempty"`
}

func makeHashTransactionsEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(hashTransactionsRequest)
		res := hashTransactionsResponse{}

		if req.ID == "" {
			res.Error = errors.New("missing required ID").Error()
			return res, nil
		}

		trx, err := s.GetTransaction(req.ID)
		if err != nil {
			res.Error = err.Error()
			return res, nil
		}

		hash, err := trx.Hash()
		if err != nil {
			res.Error = err.Error()
		}
		res.Hash = hash
		return res, nil
	}
}
