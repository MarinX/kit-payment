package repository

import (
	"encoding/json"
	"fmt"

	"github.com/MarinX/kit-payment/transaction"
	"github.com/boltdb/bolt"
)

const (
	transactionBucket = "transactions"
)

type transactionRepository struct {
	db *bolt.DB
}

func (a *transactionRepository) Store(trx *transaction.Transaction) error {
	return a.db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(transactionBucket))
		if err != nil {
			return err
		}
		buff, err := json.Marshal(trx)
		if err != nil {
			return err
		}
		return b.Put([]byte(trx.ID), buff)
	})
}

func (a *transactionRepository) Find(id string) (*transaction.Transaction, error) {
	trx := new(transaction.Transaction)
	err := a.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(transactionBucket))
		if b == nil {
			return nil
		}
		v := b.Get([]byte(id))
		if v == nil {
			return fmt.Errorf("%s transaction not found", id)
		}
		return json.Unmarshal(v, trx)
	})
	return trx, err
}
func (a *transactionRepository) FindAll() []*transaction.Transaction {
	var txs []*transaction.Transaction
	a.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(transactionBucket))
		if b == nil {
			return nil
		}
		c := b.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			tmp := &transaction.Transaction{}
			if err := json.Unmarshal(v, tmp); err != nil {
				return err
			}
			txs = append(txs, tmp)
		}
		return nil
	})
	return txs
}

func (a *transactionRepository) Delete(id string) error {
	return a.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(transactionBucket))
		return b.Delete([]byte(id))
	})
}
