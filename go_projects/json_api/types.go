package main

import (
	"math/rand"
	"time"
)

type CreateAccountRequest struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

type TransferRequest struct {
	RecipientBank int64 `json:"bank_number"`
	Amount        int64 `json:"amount"`
}

type Account struct {
	ID         int       `json:"id"`
	FirstName  string    `json:"first_name"`
	LastName   string    `json:"last_name"`
	BankNumber int64     `json:"bank_number"`
	Balance    int64     `json:"balance"`
	CreatedAt  time.Time `json:"created_at"`
}

func NewAccount(firstname, lastname string) *Account {
	return &Account{
		FirstName:  firstname,
		LastName:   lastname,
		BankNumber: int64(rand.Intn(1000000)),
		CreatedAt:  time.Now().UTC(),
	}
}
