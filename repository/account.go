package repository

import (
	"encoding/json"
	"fmt"

	"github.com/MarinX/kit-payment/account"
	"github.com/boltdb/bolt"
)

const (
	accountBucket = "accounts"
)

type accountRepository struct {
	db *bolt.DB
}

func (a *accountRepository) Store(acc *account.Account) error {
	return a.db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(accountBucket))
		if err != nil {
			return err
		}
		buff, err := json.Marshal(acc)
		if err != nil {
			return err
		}
		return b.Put([]byte(acc.ID), buff)
	})
}

func (a *accountRepository) Find(id string) (*account.Account, error) {
	acc := new(account.Account)
	err := a.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(accountBucket))
		if b == nil {
			return nil
		}
		v := b.Get([]byte(id))
		if v == nil {
			return fmt.Errorf("%s account not found", id)
		}
		return json.Unmarshal(v, acc)
	})
	return acc, err
}
func (a *accountRepository) FindAll() []*account.Account {
	var accs []*account.Account
	a.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(accountBucket))
		if b == nil {
			return nil
		}
		c := b.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			tmp := &account.Account{}
			if err := json.Unmarshal(v, tmp); err != nil {
				return err
			}
			accs = append(accs, tmp)
		}
		return nil
	})
	return accs
}
