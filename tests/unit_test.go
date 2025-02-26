package tests

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mauriciomartinezc/real-estate-api-gateway/middlewares"
)

func TestRateLimiter(t *testing.T) {
	// Define the base handler
	baseHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// Wrap the handler with the RateLimiterMiddleware
	limitedHandler := middlewares.RateLimiterMiddleware(100, 100)(baseHandler)

	// Simulate 101 requests
	for i := 0; i < 101; i++ {
		req, _ := http.NewRequest("GET", "/", nil)
		rr := httptest.NewRecorder()
		limitedHandler.ServeHTTP(rr, req)

		if i < 100 && rr.Code != http.StatusOK {
			t.Errorf("Expected status OK; got %v on request %d", rr.Code, i)
		} else if i >= 100 && rr.Code != http.StatusTooManyRequests {
			t.Errorf("Expected status Too Many Requests; got %v on request %d", rr.Code, i)
		}
	}
}

func TestSecurityHeaders(t *testing.T) {
	// Define the base handler
	baseHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// Wrap the handler with the SecurityHeadersMiddleware
	securityHandler := middlewares.SecurityHeadersMiddleware(baseHandler)

	// Create a test request and response
	req, _ := http.NewRequest("GET", "/", nil)
	rr := httptest.NewRecorder()
	securityHandler.ServeHTTP(rr, req)

	// Verify security headers
	expectedHeaders := map[string]string{
		"X-Content-Type-Options": "nosniff",
		"X-Frame-Options":        "DENY",
		"X-XSS-Protection":       "1; mode=block",
	}

	for header, expectedValue := range expectedHeaders {
		if got := rr.Header().Get(header); got != expectedValue {
			t.Errorf("Expected %s header to be %q; got %q", header, expectedValue, got)
		}
	}
}
