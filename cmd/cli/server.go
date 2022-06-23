package main

import (
	"fmt"
	"net/http"
	"strings"
)

// TransactionStore stores information for transactions on a given date.
type TransactionStore interface {
	GetTransactions(date string) string
}

// TicketServer is an HTTP interface for transaction information.
type TicketServer struct {
	store TransactionStore
}

func (t *TicketServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// date/06-22-22
	date := strings.TrimPrefix(r.URL.Path, "/date/")

	transaction := t.store.GetTransactions(date)
	if len(transaction) == 0 {
		w.WriteHeader(http.StatusNotFound)
	}

	fmt.Fprint(w, transaction)
}