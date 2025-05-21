package main

import "math/rand"

type Account struct {
	ID         int
	FirstName  string
	LastName   string
	BankNumber int64
	Balance    int64
}

func NewAccount(firstname, lastname string) *Account {
	return &Account{
		ID:         rand.Intn(10000),
		FirstName:  firstname,
		LastName:   lastname,
		BankNumber: int64(rand.Intn(1000000)),
	}
}
