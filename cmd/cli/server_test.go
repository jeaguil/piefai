package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

type StubTicketStore struct {
	transactions map[string]string
}

func (s *StubTicketStore) GetTransactions(date string) string {
	return s.transactions[date]
}

func TestMissingData(t *testing.T) {

	store := StubTicketStore{
		map[string]string{
			"date/Missing": "data",
		},
	}

	server := &TicketServer{&store}

	t.Run("Returns 404 on missing data", func(t *testing.T) {
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

	// Current Store:
	store := StubTicketStore{
		map[string]string{
			"date/06": "20",
			"date/07": "10",
		},
	}

	server := &TicketServer{&store}

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
