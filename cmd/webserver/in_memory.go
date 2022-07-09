package main

/*
TICKER STORE
*/

type InMemoryTickerStore struct {
	store map[string]string
}

func NewInMemoryTickerStore() *InMemoryTickerStore {
	return &InMemoryTickerStore{map[string]string{}}
}

/*
TRANSACTION STORE
*/

type InMemoryTransactionStore struct {
	store map[string]string
}

func NewInMemoryTransactionStore() *InMemoryTransactionStore {
	return &InMemoryTransactionStore{map[string]string{}}
}
