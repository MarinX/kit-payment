# Payment service with go-kit
[![Build Status](https://travis-ci.org/MarinX/kit-payment.svg?branch=master)](https://travis-ci.org/MarinX/kit-payment)
[![GoDoc](https://godoc.org/github.com/MarinX/kit-payment?status.svg)](https://godoc.org/github.com/MarinX/kit-payment)
[![License MIT](https://img.shields.io/badge/license-MIT-brightgreen.svg?style=flat)](LICENSE)

## Getting started

### Building
Clone the repo
```sh
git clone https://github.com/MarinX/kit-payment
```
Fetch dependencies
```sh
go get ./...
```
Build as any other go project
```sh
go build
```

### Usage
```sh
Usage of ./kit-payment:
  -http.addr string
        HTTP listen address (default ":8080")
```

### Storage
Kit-payment is using embedded key/value database called [boltdb](https://github.com/boltdb/bolt).

## Endpoints

### Accounts
#### Creating account
```sh
curl -H "Content-Type: application/json" -X POST http://localhost:8080/accounts
```

#### Listing accounts
```sh
curl -H "Content-Type: application/json" -X GET http://localhost:8080/accounts
```

#### Adding balance to account
Example addding $100 to account id `3479d3a8-42c4-4b40-8a1d-1b1661d7b6ef`
```sh
curl -d '{"currency":"USD", "amount":100}' -H "Content-Type: application/json" -X POST http://localhost:8080/accounts/3479d3a8-42c4-4b40-8a1d-1b1661d7b6ef/balances
```

### Transactions
#### List Transactions
```sh
curl -H "Content-Type: application/json" -X GET http://localhost:8080/transactions
```

#### Creating Transaction
Example of moving $50 from account `3479d3a8-42c4-4b40-8a1d-1b1661d7b6ef` to account `06e39e77-776a-4694-bc59-fea69bc8afd8`
```sh
curl -d '{"from":"3479d3a8-42c4-4b40-8a1d-1b1661d7b6ef", "to":"06e39e77-776a-4694-bc59-fea69bc8afd8", "currency":"USD", "amount":50}' -H "Content-Type: application/json" -X POST http://localhost:8080/transactions
```
If success, it will return a created transaction object
```sh
{"transaction":{"id":"fecf39a1-c4f2-4706-8eca-bc71f310eeb6","from":"3479d3a8-42c4-4b40-8a1d-1b1661d7b6ef","to":"06e39e77-776a-4694-bc59-fea69bc8afd8","status":"created","amount":50,"currency":"USD"}}
```

#### Get Transaction
Example of getting single transaction by id `fecf39a1-c4f2-4706-8eca-bc71f310eeb6`
```sh
curl -H "Content-Type: application/json" -X GET http://localhost:8080/transactions/fecf39a1-c4f2-4706-8eca-bc71f310eeb6
```

#### Commit Transaction
Once the transaction is created, you need to commit.
Example is commiting created transaction `fecf39a1-c4f2-4706-8eca-bc71f310eeb6`
```sh
curl -H "Content-Type: application/json" -X PUT http://localhost:8080/transactions/fecf39a1-c4f2-4706-8eca-bc71f310eeb6/commit
```
Now you can check the status with `Get Transaction` method.
If account has enough balance, you will see the change on amount when listing accounts.

#### Transaction verification
It provides a interface for [merkle tree](https://github.com/cbergoon/merkletree) so we can check if all transactions are verified.
Example of checking our last transaction `fecf39a1-c4f2-4706-8eca-bc71f310eeb6`
```sh
curl -H "Content-Type: application/json" -X GET http://localhost:8080/transactions/fecf39a1-c4f2-4706-8eca-bc71f310eeb6/hash
```
will return the hash for verification
```sh
{"hash":"fdf227bade5496e59824a4c9ef59ec992c61b4521fe0003c5be4d79cec3c885c"}
```

## Tests
Nothing fancy, just run
```sh
go test -v ./...
```

## Roadmap
- Support currency conversion
- Support user accounts
- Integrate merkle tree so we can verify transactions and extend (mining?)

<hr/>
PR's welcome :)