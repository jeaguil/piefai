package main

import (
	"log"
	"net/http"

	config "github.com/jeaguil/piefai/_con"

	polygon "github.com/polygon-io/client-go/rest"
)

type InMemoryTransactionStore struct{}

func (i *InMemoryTransactionStore) GetTransactions(name string) string {
	return "dsanklfnslaknfksa"
}

func main() {
	c := polygon.New(config.GetAPIKey())
	_ = c

	server := &TicketServer{&InMemoryTransactionStore{}}
	log.Fatal(http.ListenAndServe(":5000", server))
}
