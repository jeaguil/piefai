package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

type StubTicketStore struct {
	transactions map[string]string
	altCalls     []string
}

// Gets history of transactions on a given date
func (s *StubTicketStore) GetTransactions(date string) string {
	return s.transactions[date]
}

func (s *StubTicketStore) ProcessTransaction(alt string) {
	s.altCalls = append(s.altCalls, alt)
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
		request := newGetTransactionsRequest("Miss")
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
// 		request := newGetTransactionsRequest("06")
// 		response := httptest.NewRecorder()

// 		server.ServeHTTP(response, request)

// 		got := response.Body.String()
// 		want := "20"

// 		assertResponseBody(t, got, want)
// 	})
// }

func TestGETTransactions(t *testing.T) {
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
			request := newGetTransactionsRequest(tt.date)
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

func newGetTransactionsRequest(date string) *http.Request {
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
	store := NewInMemoryTransactionStore()
	server := TicketServer{store}

	// POST
	req, _ := http.NewRequest(http.MethodPost, "/date/010", nil)
	server.ServeHTTP(response, req)

	// GET
	server.ServeHTTP(httptest.NewRecorder(), newGetTransactionsRequest("010"))

	assertStatus(t, response.Code, http.StatusAccepted)

	assertResponseBody(t, httptest.NewRecorder().Body.String(), "")
}
