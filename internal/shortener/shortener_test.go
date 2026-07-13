package shortener

import (
	"context"
	"database/sql"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"url-shortener/internal/db"
)

// MockDB implements the DBQuerier interface for testing purposes
type MockDB struct {
	CreateURLFunc    func(ctx context.Context, arg db.CreateURLParams) (sql.Result, error)
	GetURLByCodeFunc func(ctx context.Context, shortCode string) (db.Url, error)
	RecordClickFunc  func(ctx context.Context, arg db.RecordClickParams) error
}

func (m *MockDB) CreateURL(ctx context.Context, arg db.CreateURLParams) (sql.Result, error) {
	return m.CreateURLFunc(ctx, arg)
}

func (m *MockDB) GetURLByCode(ctx context.Context, shortCode string) (db.Url, error) {
	return m.GetURLByCodeFunc(ctx, shortCode)
}

func (m *MockDB) RecordClick(ctx context.Context, arg db.RecordClickParams) error {
	return m.RecordClickFunc(ctx, arg)
}

func TestShortenHandler(t *testing.T) {
	// Define the test matrix table structure
	tests := []struct {
		name           string
		requestBody    string
		mockBehavior   func(ctx context.Context, arg db.CreateURLParams) (sql.Result, error) // ✅ Fixed signature
		expectedStatus int
		expectedBody   string
	}{
		{
			name:        "Success - URL Shortened Cleanly",
			requestBody: `{"long_url": "https://prontobroadband.com"}`,
			mockBehavior: func(ctx context.Context, arg db.CreateURLParams) (sql.Result, error) { // ✅ Fixed signature
				return nil, nil // Return no database errors
			},
			expectedStatus: http.StatusCreated,
			expectedBody:   `"short_code"`,
		},
		{
			name:           "Failure - Broken JSON Request Payload",
			requestBody:    `{"long_url": broken-json-syntax`,
			mockBehavior:   nil, // Database won't even be called
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `"error"`,
		},
		{
			name:        "Failure - Database Storage Error Triggered",
			requestBody: `{"long_url": "https://failing-database-link.com"}`,
			mockBehavior: func(ctx context.Context, arg db.CreateURLParams) (sql.Result, error) { // ✅ Fixed signature
				return nil, errors.New("mysql connection timed out")
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `"Unable to save link"`,
		},
	}

	// Iterate through the matrix table cases sequentially
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// 1. Initialize our mock with the current case behavioral closure
			mock := &MockDB{
				CreateURLFunc: tc.mockBehavior,
			}
			handler := NewHandler(mock)

			// 2. Build the fake HTTP Request and ResponseRecorder tracking blocks
			req := httptest.NewRequest(http.MethodPost, "/api/v1/shorten", strings.NewReader(tc.requestBody))
			req.Header.Set("Content-Type", "application/json")
			rr := httptest.NewRecorder()

			// 3. Execute the target handler directly
			handler.Shorten(rr, req)

			// 4. Assert and verify the HTTP status responses
			if rr.Code != tc.expectedStatus {
				t.Errorf("%s: expected status code %d, got %d", tc.name, tc.expectedStatus, rr.Code)
			}

			// 5. Assert and check that the correct response payload data strings match up
			if !strings.Contains(rr.Body.String(), tc.expectedBody) {
				t.Errorf("%s: expected response to contain %q, got %q", tc.name, tc.expectedBody, rr.Body.String())
			}
		})
	}
}
