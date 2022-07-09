/*
TESTS FOR TickerServer;
server := &TicketServer{NewInMemoryTickerStore()}

Handles func for TickerStore:
	GetTicker
	ProcessTicker

Transaction Server uses the same func of these tests.
*/

package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

type StubTicketStore struct {
	transactions map[string]string
	altCalls     []string
}

// Gets history of transactions on a given date
func (s *StubTicketStore) GetTicker(date string) string {
	return s.transactions[date]
}

func (s *StubTicketStore) ProcessTicker(date, alt string) {
	s.altCalls = append(s.altCalls, date)
}

// Testing Store:
var store = StubTicketStore{
	map[string]string{
		"date/Missing": "data",
		"date/06":      "20",
		"date/07":      "10",
	},
	nil,
}

// Same server for all tests
var server = &TicketServer{&store}

// Response is the same for all tests
var response = httptest.NewRecorder()

func TestMissingData(t *testing.T) {
	t.Run("Returns 404 on missing data", func(t *testing.T) {
		// Suppose to return 404 since Miss is not in date/Miss
		request := newGetTickerRequest("Miss")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		got := response.Code
		want := http.StatusNotFound

		if got != want {
			t.Errorf("got status %d want %d", got, want)
		}
	})
}

// func TestTicketServer(t *testing.T) {
// 	server := &TicketServer{}

// 	t.Run("Returns map of data on a given date", func(t *testing.T) {
// 		request := newGetTickerRequest("06")
// 		response := httptest.NewRecorder()

// 		server.ServeHTTP(response, request)

// 		got := response.Body.String()
// 		want := "20"

// 		assertResponseBody(t, got, want)
// 	})
// }

func TestGetTicker(t *testing.T) {
	tests := []struct {
		name               string
		date               string
		expectedHTTPStatus int
		expectedResponse   string
	}{
		{
			name:               "Given M",
			date:               "06",
			expectedHTTPStatus: http.StatusOK,
			expectedResponse:   "20",
		},
		{
			name:               "Formated Date",
			date:               "06-22-24",
			expectedHTTPStatus: http.StatusNotFound,
			expectedResponse:   "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := newGetTickerRequest(tt.date)
			response := httptest.NewRecorder()

			server.ServeHTTP(response, request)

			assertStatus(t, response.Code, tt.expectedHTTPStatus)
			assertResponseBody(t, response.Body.String(), tt.expectedResponse)
		})
	}
}

func TestStoreTransactions(t *testing.T) {
	t.Run("it returns accepted on POST", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodPost, "/date/07", nil)

		server.ServeHTTP(response, request)

		assertStatus(t, response.Code, http.StatusAccepted)

		if len(store.altCalls) != 1 {
			t.Fatalf("got %d calls to Alt want %d", len(store.altCalls), 1)
		}
	})
}

func TestStoreCalls(t *testing.T) {
	t.Run("it processes transactions on POST", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodPost, "/date/08", nil)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatus(t, response.Code, http.StatusAccepted)

		want := "08"
		if store.altCalls[1] != want {
			t.Errorf("did not store correct transaction for date, got %q want %q", store.altCalls[1], want)
		}
	})
}

func assertStatus(t *testing.T, i1, i2 int) {
	if i1 != i2 {
		t.Errorf("unexpected response code, got %d want %d", i1, i2)
	}
}

func newGetTickerRequest(date string) *http.Request {
	request, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("date/%s", date), nil)
	return request
}

func assertResponseBody(t testing.TB, got, want string) {
	t.Helper()
	if got != want {
		t.Errorf("response body is wrong, got %q want %q", got, want)
	}
}

func TestProcessingAndRetrieving(t *testing.T) {
	store := NewInMemoryTickerStore()
	server := TicketServer{store}

	usr := Ticker{
		Symbol: "MSFT",
	}

	body, _ := json.Marshal(usr)

	year_month := "22-06"
	// POST
	req, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("/date/%s", year_month), bytes.NewBuffer(body))

	// Close response stream, once response is read.
	defer req.Body.Close()

	ctx, cancel := context.WithTimeout(req.Context(), 1*time.Millisecond)
	defer cancel()

	req = req.WithContext(ctx)

	server.ServeHTTP(response, req)

	response_integration := httptest.NewRecorder()

	// GET
	server.ServeHTTP(response_integration, newGetTickerRequest(year_month))

	assertStatus(t, response.Code, http.StatusAccepted)

	assertResponseBody(t, response_integration.Body.String(), usr.Symbol)
}
