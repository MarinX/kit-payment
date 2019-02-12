package repository

import (
	"os"
	"testing"

	"github.com/MarinX/kit-payment/transaction"

	"github.com/MarinX/kit-payment/account"
)

func openRepo(t *testing.T) *Repository {
	repo, err := New()
	if err != nil {
		t.Error(err)
		t.Fail()
		return nil
	}
	return repo
}
func closeRepo(t *testing.T, r *Repository) {
	if err := r.Close(); err != nil {
		t.Error(err)
		t.Fail()
	}
	os.Remove("data.db")
}

func TestAccountRepository(t *testing.T) {
	repo := openRepo(t)
	defer closeRepo(t, repo)

	accRepo := repo.Account()
	tmpAcc := account.New()

	if err := accRepo.Store(tmpAcc); err != nil {
		t.Errorf("error storing account %v", err)
		return
	}

	expected, err := accRepo.Find(tmpAcc.ID)
	if err != nil {
		t.Errorf("error getting account %v", err)
		return
	}
	if tmpAcc.ID != expected.ID {
		t.Errorf("account ids does not match, want %v got %v", tmpAcc.ID, expected.ID)
	}

	tmpAcc.SetBalance(account.Currency("USD"), 1)
	if err := accRepo.Store(tmpAcc); err != nil {
		t.Errorf("error updating account %v", err)
		return
	}

	allAcc := accRepo.FindAll()

	if len(allAcc) != 1 {
		t.Errorf("invalid number of accounts")
	}
}

func TestTransactionRepository(t *testing.T) {
	repo := openRepo(t)
	defer closeRepo(t, repo)

	txRepo := repo.Transaction()
	tmpTx := transaction.New("123", "222", account.Currency("USD"), 10)

	if err := txRepo.Store(tmpTx); err == nil {
		t.Errorf("missing ID should yield error,got nil")
		return
	}

	tmpTx.Create()

	if err := txRepo.Store(tmpTx); err != nil {
		t.Errorf("error creating transaction %v", err)
		return
	}

	expected, err := txRepo.Find(tmpTx.ID)
	if err != nil {
		t.Errorf("error finding transaction %v", err)
		return
	}
	if tmpTx.ID != expected.ID {
		t.Errorf("transaction ids does not match, want %v got %v", tmpTx.ID, expected.ID)
		return
	}

	allTxs := txRepo.FindAll()

	if len(allTxs) != 1 {
		t.Errorf("invalid number of transactions")
	}

	if err := txRepo.Delete(tmpTx.ID); err != nil {
		t.Errorf("error deleting transaction %v", err)
		return
	}

	allTxs = txRepo.FindAll()
	if len(allTxs) != 0 {
		t.Errorf("invalid number of transactions")
	}
}
