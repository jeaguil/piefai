package main

import (
	"log"
	"net/http"
	"strings"

	config "github.com/jeaguil/piefai/_con"

	polygon "github.com/polygon-io/client-go/rest"
)

func NewInMemoryTickerStore() *InMemoryTickerStore {
	return &InMemoryTickerStore{map[string]string{}}
}

type InMemoryTickerStore struct {
	store map[string]string
}

func (i *InMemoryTickerStore) GetTicker(date string) string {
	month := strings.TrimPrefix(date, "date/")
	return i.store[month]
}

func (i *InMemoryTickerStore) ProcessTicker(date, transaction string) {
	i.store[date] = transaction
}

func NewInMemoryTransactionStore() *InMemoryTransactionStore {
	return &InMemoryTransactionStore{map[string]string{}}
}

func (i *InMemoryTransactionStore) GetJSONTransactions() []Transaction {
	var t []Transaction
	for ticker := range i.store {
		t = append(t, Transaction{ticker})
	}
	return t
}

func (i *InMemoryTransactionStore) ProcessTransaction(t Transaction) {
}

type InMemoryTransactionStore struct {
	store map[string]string
}

func main() {
	c := polygon.New(config.GetAPIKey())
	_ = c

	server := &TicketServer{NewInMemoryTickerStore()}
	// server := NewTransactionServr(NewInMemoryTransactionStore())
	log.Print("Local server on port :5000")
	log.Fatal(http.ListenAndServe(":5000", server))
}
