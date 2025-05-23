package main

import (
	"math/rand"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type loginRequest struct {
	Number   int64  `json:"number"`
	Password string `json:"password"`
}

type CreateAccountRequest struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Password  string `json:"-"`
}

type TransferRequest struct {
	RecipientBank int64 `json:"bank_number"`
	Amount        int64 `json:"amount"`
}

type Account struct {
	ID                int       `json:"id"`
	FirstName         string    `json:"first_name"`
	LastName          string    `json:"last_name"`
	BankNumber        int64     `json:"bank_number"`
	EncryptedPassword string    `json:"-"`
	Balance           int64     `json:"balance"`
	CreatedAt         time.Time `json:"created_at"`
}

func NewAccount(firstname, lastname, password string) (*Account, error) {
	encpw, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	return &Account{
		FirstName:         firstname,
		LastName:          lastname,
		EncryptedPassword: string(encpw),
		BankNumber:        int64(rand.Intn(1000000)),
		CreatedAt:         time.Now().UTC(),
	}, nil
}
