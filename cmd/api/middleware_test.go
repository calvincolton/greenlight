package main

import (
	"expvar"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestRateLimitMiddleware(t *testing.T) {
	app := &application{}

	app.config.limiter.enabled = true
	app.config.limiter.rps = 2 // Allow 2 requests per second for testing
	app.config.limiter.burst = 2

	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	handler := app.rateLimit(testHandler)

	for i := 0; i < 4; i++ {
		req, err := http.NewRequest(http.MethodGet, "/v1/healthcheck", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()

		handler.ServeHTTP(rr, req)

		if i < 2 && rr.Code != http.StatusOK {
			t.Errorf("expected status 200 OK, got %v", rr.Code)
		}

		if i >= 2 && rr.Code != http.StatusTooManyRequests {
			t.Errorf("expected status 429 Too Many Requests, got %v", rr.Code)
		}
	}

	// Ensure some delay between subsequent tests
	time.Sleep(time.Second)
}

// WARNING: may not be deterministic
func TestMetricsMiddleware(t *testing.T) {
	// resetMetrics() // Reset metrics before each test

	app := &application{}

	// Create a simple test handler that just returns a 200 OK response
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// Wrap the test handler with the metrics middleware
	handler := app.metrics(testHandler)

	// Send a request
	req, err := http.NewRequest(http.MethodGet, "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	// Check expvar metrics
	if totalRequestsReceived := expvar.Get("total_requests_received").(*expvar.Int).Value(); totalRequestsReceived != 1 {
		t.Errorf("expected 1 total request received, got %d", totalRequestsReceived)
	}

	if totalResponsesSent := expvar.Get("total_responses_sent").(*expvar.Int).Value(); totalResponsesSent != 1 {
		t.Errorf("expected 1 total response sent, got %d", totalResponsesSent)
	}

	if totalProcessingTimeMicroseconds := expvar.Get("total_processing_time_μs").(*expvar.Int).Value(); totalProcessingTimeMicroseconds <= 0 {
		t.Errorf("expected positive processing time, got %d", totalProcessingTimeMicroseconds)
	}

	if totalResponsesSentByStatus := expvar.Get("total_responses_sent_by_status").(*expvar.Map).Get("200").(*expvar.Int).Value(); totalResponsesSentByStatus != 1 {
		t.Errorf("expected 1 total response sent by status 200, got %d", totalResponsesSentByStatus)
	}
}

func TestEnableCORSMiddleware(t *testing.T) {
	app := &application{}
	app.config.cors.trustedOrigins = []string{"https://trusted.com"}

	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// Wrap the test handler with the enableCORS middleware
	handler := app.enableCORS(testHandler)

	tests := []struct {
		name               string
		origin             string
		method             string
		expectedStatusCode int
		expectedHeaders    map[string]string
	}{
		{
			name:               "No Origin Header",
			origin:             "",
			method:             http.MethodGet,
			expectedStatusCode: http.StatusOK,
			expectedHeaders:    map[string]string{"Vary": "Origin, Access-Control-Request-Method"},
		},
		{
			name:               "Trusted Origin",
			origin:             "https://trusted.com",
			method:             http.MethodGet,
			expectedStatusCode: http.StatusOK,
			expectedHeaders:    map[string]string{"Access-Control-Allow-Origin": "https://trusted.com", "Vary": "Origin, Access-Control-Request-Method"},
		},
		{
			name:               "Trusted Origin with Preflight",
			origin:             "https://trusted.com",
			method:             http.MethodOptions,
			expectedStatusCode: http.StatusOK,
			expectedHeaders: map[string]string{
				"Access-Control-Allow-Origin":  "https://trusted.com",
				"Access-Control-Allow-Methods": "OPTIONS, PUT, PATCH, DELETE",
				"Access-Control-Allow-Headers": "Authorization, Content-Type",
				"Vary":                         "Origin, Access-Control-Request-Method",
			},
		},
		{
			name:               "Untrusted Origin",
			origin:             "https://untrusted.com",
			method:             http.MethodGet,
			expectedStatusCode: http.StatusOK,
			expectedHeaders:    map[string]string{"Vary": "Origin, Access-Control-Request-Method"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest(tt.method, "/", nil)
			if err != nil {
				t.Fatal(err)
			}

			req.Header.Set("Origin", tt.origin)
			if tt.method == http.MethodOptions {
				req.Header.Set("Access-Control-Request-Method", "PUT")
			}

			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			if rr.Code != tt.expectedStatusCode {
				t.Errorf("expected status %d, got %d", tt.expectedStatusCode, rr.Code)
			}

			for key, expectedValue := range tt.expectedHeaders {
				if value := rr.Header().Get(key); value != expectedValue {
					t.Errorf("expected header %s to be %s, got %s", key, expectedValue, value)
				}
			}
		})
	}
}
