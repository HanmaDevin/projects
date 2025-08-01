package main

import (
	"flag"
	"fmt"
	"log"
)

func seedAccount(store Storage, firstName, lastName, pw string) *Account {
	acc, err := NewAccount(firstName, lastName, pw)
	if err != nil {
		log.Fatal(err)
	}

	if err := store.CreateAccount(acc); err != nil {
		log.Fatal(err)
	}

	return acc
}

func seedAccounts(s Storage) {
	seedAccount(s, "anthony", "GG", "blajfsd")
}

func main() {
	seed := flag.Bool("seed", false, "seed the database")
	flag.Parse()

	store, err := NewPostgresStore()
	if err != nil {
		log.Fatal(err)
	}

	if err := store.Init(); err != nil {
		log.Fatal(err)
	}

	if *seed {
		fmt.Println("seeding database")
		seedAccounts(store)
	}

	server := NewAPIServer("localhost:3000", store)
	server.Run()
}
