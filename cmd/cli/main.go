package main

import (
	"log"
	"net/http"

	config "github.com/jeaguil/piefai/_con"

	polygon "github.com/polygon-io/client-go/rest"
)

func NewInMemoryTransactionStore() *InMemoryTransactionStore {
	return &InMemoryTransactionStore{map[string]string{}, nil}
}

type InMemoryTransactionStore struct {
	store    map[string]string
	altCalls []string
}

func (i *InMemoryTransactionStore) GetTransactions(transaction string) string {
	return i.store[transaction]
}

func (i *InMemoryTransactionStore) ProcessTransaction(transaction string) {
	i.altCalls = append(i.altCalls, transaction)
}

func main() {
	c := polygon.New(config.GetAPIKey())
	_ = c

	server := &TicketServer{NewInMemoryTransactionStore()}
	log.Print("Local server on port :5000")
	log.Fatal(http.ListenAndServe(":5000", server))
}
