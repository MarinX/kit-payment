package account

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/gorilla/mux"

	kitlog "github.com/go-kit/kit/log"
	kithttp "github.com/go-kit/kit/transport/http"
)

// MakeHandler returns a handler for the account service.
func MakeHandler(as Service, logger kitlog.Logger) http.Handler {
	r := mux.NewRouter()

	opts := []kithttp.ServerOption{
		kithttp.ServerErrorLogger(logger),
	}

	accountsHandler := kithttp.NewServer(
		makeAccountsEndpoint(as),
		decodeAccountsRequest,
		encodeResponse,
		opts...,
	)

	accountsListHandler := kithttp.NewServer(
		makeListAccountsEndpoint(as),
		decodeListAccountsRequest,
		encodeResponse,
		opts...,
	)

	accountsBalanceHandler := kithttp.NewServer(
		makeAccountsBalanceEndpoint(as),
		decodeAccountsBalanceRequest,
		encodeResponse,
		opts...,
	)

	r.Handle("/accounts", accountsHandler).Methods("POST")
	r.Handle("/accounts", accountsListHandler).Methods("GET")
	r.Handle("/accounts/{id}/balances", accountsBalanceHandler).Methods("POST")

	return r
}

func decodeAccountsRequest(_ context.Context, r *http.Request) (interface{}, error) {
	return accountsRequest{}, nil
}

func decodeListAccountsRequest(_ context.Context, r *http.Request) (interface{}, error) {
	return listAccountsRequest{}, nil
}

func decodeAccountsBalanceRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		return nil, errors.New("bad request")
	}

	var body accountsBalanceRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		return nil, err
	}

	body.AccountID = id
	return body, nil
}

func encodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}
