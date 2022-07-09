package main

import (
	"log"
	"net/http"
)

func main() {

	// server := &TicketServer{NewInMemoryTickerStore()}
	server := NewTransactionServer(NewInMemoryTransactionStore())
	log.Print("Local server on port :5000")
	log.Fatal(http.ListenAndServe(":5000", server))
}
