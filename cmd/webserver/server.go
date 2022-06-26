package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

type Ticker struct {
	Symbol string `json:"symbol"`
}

// TickerStore stores information for transactions on a given date.
type TickerStore interface {
	GetTicker(date string) string
	ProcessTicker(date, transaction string)
}

// TicketServer is an HTTP interface for transaction information.
type TicketServer struct {
	store TickerStore
}

type Transaction struct {
	Ticker string
}

// TransactionStore stores information all all transactions made.
type TransactionStore interface {
	GetJSONTransactions() []Transaction
	ProcessTransaction(t Transaction)
}

type TransactionServer struct {
	store TransactionStore
	http.Handler
}

func NewTransactionServr(store TransactionStore) *TransactionServer {
	t := new(TransactionServer)

	t.store = store

	router := http.NewServeMux()
	router.Handle("/transactions", http.HandlerFunc(t.transactionHandler))
	router.Handle("/settle", http.HandlerFunc(t.settleHandler))

	t.Handler = router

	return t
}

func (t *TransactionServer) transactionHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	json.NewEncoder(w).Encode(t.store.GetJSONTransactions())
}

func (t *TransactionServer) settleHandler(w http.ResponseWriter, r *http.Request) {
	req_body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatalln(err)
	}
	strbdy := string(req_body)
	var tr_set Transaction
	tr_set.Ticker = strbdy

	t.store.ProcessTransaction(tr_set)
	w.WriteHeader(http.StatusOK)
}

func (t *TicketServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// date/06-22-22
	date := strings.TrimPrefix(r.URL.Path, "/date/")
	switch r.Method {
	case http.MethodPost:
		if r.Body != nil {
			body, err := ioutil.ReadAll(r.Body)

			if err != nil {
				log.Fatalln(err)
			}

			sp := string(body)
			var uss Ticker
			json.Unmarshal([]byte(sp), &uss)

			t.store.ProcessTicker(date, uss.Symbol)
			return
		}
		t.store.ProcessTicker(date, "")
		w.WriteHeader(http.StatusAccepted)
		return
	case http.MethodGet:
		transaction := t.store.GetTicker(date)
		if len(transaction) == 0 {
			w.WriteHeader(http.StatusNotFound)
		}

		fmt.Fprint(w, transaction)
	}

}
