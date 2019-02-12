package transaction

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/gorilla/mux"

	kitlog "github.com/go-kit/kit/log"
	kithttp "github.com/go-kit/kit/transport/http"
)

// MakeHandler returns a handler for the transaction service.
func MakeHandler(ts Service, logger kitlog.Logger) http.Handler {
	r := mux.NewRouter()

	opts := []kithttp.ServerOption{
		kithttp.ServerErrorLogger(logger),
	}

	transactionsHandler := kithttp.NewServer(
		makeTransactionsEndpoint(ts),
		decodeTransactionsRequest,
		encodeResponse,
		opts...,
	)

	transactionsListHandler := kithttp.NewServer(
		makeListTransactionsEndpoint(ts),
		decodeListTransactionsRequest,
		encodeResponse,
		opts...,
	)

	transactionsGetHandler := kithttp.NewServer(
		makeGetTransactionsEndpoint(ts),
		decodeGetTransactionsRequest,
		encodeResponse,
		opts...,
	)

	transactionsCommitHandler := kithttp.NewServer(
		makeCommitTransactionsEndpoint(ts),
		decodeCommitTransactionsRequest,
		encodeResponse,
		opts...,
	)

	transactionsHashHandler := kithttp.NewServer(
		makeHashTransactionsEndpoint(ts),
		decodeHashTransactionsRequest,
		encodeResponse,
		opts...,
	)

	r.Handle("/transactions", transactionsHandler).Methods("POST")
	r.Handle("/transactions", transactionsListHandler).Methods("GET")
	r.Handle("/transactions/{id}", transactionsGetHandler).Methods("GET")
	r.Handle("/transactions/{id}/commit", transactionsCommitHandler).Methods("PUT")
	r.Handle("/transactions/{id}/hash", transactionsHashHandler).Methods("GET")

	return r
}

func decodeTransactionsRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var body transactionsRequest
	if r.Body != nil {
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			return nil, err
		}
	}
	return body, nil
}

func decodeGetTransactionsRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		return nil, errors.New("bad request")
	}
	return getTransactionsRequest{
		ID: id,
	}, nil
}

func decodeListTransactionsRequest(_ context.Context, r *http.Request) (interface{}, error) {
	return listTransactionsRequest{}, nil
}

func decodeCommitTransactionsRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		return nil, errors.New("bad request")
	}
	return commitTransactionsRequest{
		ID: id,
	}, nil
}

func decodeHashTransactionsRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		return nil, errors.New("bad request")
	}
	return hashTransactionsRequest{
		ID: id,
	}, nil
}

func encodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}
