package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/MarinX/kit-payment/transaction"

	"github.com/MarinX/kit-payment/account"

	"github.com/MarinX/kit-payment/repository"
	"github.com/go-kit/kit/log"
)

func main() {
	var (
		httpAddr = flag.String("http.addr", ":8080", "HTTP listen address")
	)
	flag.Parse()

	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = log.With(logger, "ts", log.DefaultTimestampUTC)
		logger = log.With(logger, "caller", log.DefaultCaller)
	}

	repo, err := repository.New()
	if err != nil {
		logger.Log("exit", err)
		return
	}
	defer repo.Close()

	var (
		accountRepo     = repo.Account()
		transactionRepo = repo.Transaction()
	)

	var (
		as = account.NewService(accountRepo)
		ts = transaction.NewService(transactionRepo, accountRepo, logger)
	)

	httpLogger := log.With(logger, "component", "http")
	mux := http.NewServeMux()
	mux.Handle("/accounts", account.MakeHandler(as, httpLogger))
	mux.Handle("/accounts/", account.MakeHandler(as, httpLogger))
	mux.Handle("/transactions", transaction.MakeHandler(ts, httpLogger))
	mux.Handle("/transactions/", transaction.MakeHandler(ts, httpLogger))

	go ts.Watch()

	errs := make(chan error, 2)
	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errs <- fmt.Errorf("%s", <-c)
	}()

	go func() {
		logger.Log("transport", "HTTP", "addr", *httpAddr)
		errs <- http.ListenAndServe(*httpAddr, mux)
	}()

	logger.Log("exit", <-errs)

}
