package http

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Zillaforge/cloud-sdk/internal/types"
)

func TestClient_Do_Success(t *testing.T) {
	// Create test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify authorization header
		if r.Header.Get("Authorization") != "Bearer test-token" {
			t.Errorf("expected Authorization header 'Bearer test-token', got '%s'", r.Header.Get("Authorization"))
		}

		// Return success response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	}))
	defer server.Close()

	client := NewClient(server.URL, "test-token", &http.Client{Timeout: 5 * time.Second}, nil)

	var result map[string]string
	err := client.Do(context.Background(), &Request{
		Method: "GET",
		Path:   "/test",
	}, &result)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result["status"] != "ok" {
		t.Errorf("expected status 'ok', got '%s'", result["status"])
	}
}

func TestClient_Do_Timeout(t *testing.T) {
	// Create test server that delays response
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		time.Sleep(200 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	// Create client with short timeout
	httpClient := &http.Client{Timeout: 50 * time.Millisecond}
	client := NewClient(server.URL, "test-token", httpClient, nil)

	err := client.Do(context.Background(), &Request{
		Method: "GET",
		Path:   "/test",
	}, nil)

	if err == nil {
		t.Fatal("expected timeout error, got nil")
	}

	var sdkErr *types.SDKError
	if !errors.As(err, &sdkErr) {
		t.Fatalf("expected SDKError, got %T", err)
	}

	if sdkErr.Meta["category"] != "timeout" {
		t.Errorf("expected timeout error category, got %v", sdkErr.Meta["category"])
	}
}

func TestClient_Do_ContextCancellation(t *testing.T) {
	// Create test server that delays response
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		time.Sleep(200 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := NewClient(server.URL, "test-token", &http.Client{Timeout: 5 * time.Second}, nil)

	// Create cancelable context
	ctx, cancel := context.WithCancel(context.Background())

	// Cancel immediately
	cancel()

	err := client.Do(ctx, &Request{
		Method: "GET",
		Path:   "/test",
	}, nil)

	if err == nil {
		t.Fatal("expected canceled error, got nil")
	}

	var sdkErr *types.SDKError
	if !errors.As(err, &sdkErr) {
		t.Fatalf("expected SDKError, got %T", err)
	}

	if sdkErr.Meta["category"] != "canceled" {
		t.Errorf("expected canceled error category, got %v", sdkErr.Meta["category"])
	}
}

func TestClient_Do_HTTPError(t *testing.T) {
	tests := []struct {
		name            string
		statusCode      int
		responseBody    interface{}
		expectErrorCode int
		expectMessage   string
	}{
		{
			name:       "structured error response",
			statusCode: 400,
			responseBody: ErrorResponse{
				ErrorCode: 1001,
				Message:   "Invalid request",
				Meta:      map[string]interface{}{"field": "name"},
			},
			expectErrorCode: 1001,
			expectMessage:   "Invalid request",
		},
		{
			name:            "unstructured error response",
			statusCode:      500,
			responseBody:    "Internal Server Error",
			expectErrorCode: 0,
			expectMessage:   "HTTP 500",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tt.statusCode)

				if errResp, ok := tt.responseBody.(ErrorResponse); ok {
					_ = json.NewEncoder(w).Encode(errResp)
				} else {
					_, _ = w.Write([]byte(tt.responseBody.(string)))
				}
			}))
			defer server.Close()

			client := NewClient(server.URL, "test-token", &http.Client{Timeout: 5 * time.Second}, nil)

			err := client.Do(context.Background(), &Request{
				Method: "GET",
				Path:   "/test",
			}, nil)

			if err == nil {
				t.Fatal("expected error, got nil")
			}

			var sdkErr *types.SDKError
			if !errors.As(err, &sdkErr) {
				t.Fatalf("expected SDKError, got %T", err)
			}

			if sdkErr.StatusCode != tt.statusCode {
				t.Errorf("expected status code %d, got %d", tt.statusCode, sdkErr.StatusCode)
			}

			if sdkErr.ErrorCode != tt.expectErrorCode {
				t.Errorf("expected error code %d, got %d", tt.expectErrorCode, sdkErr.ErrorCode)
			}

			if sdkErr.Message != tt.expectMessage {
				t.Errorf("expected message '%s', got '%s'", tt.expectMessage, sdkErr.Message)
			}
		})
	}
}

func TestClient_Do_Retry(t *testing.T) {
	attemptCount := 0

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		attemptCount++

		// Fail first 2 attempts, succeed on 3rd
		if attemptCount < 3 {
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		}

		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	}))
	defer server.Close()

	client := NewClient(server.URL, "test-token", &http.Client{Timeout: 5 * time.Second}, nil)

	var result map[string]string
	err := client.Do(context.Background(), &Request{
		Method: "GET",
		Path:   "/test",
	}, &result)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if attemptCount != 3 {
		t.Errorf("expected 3 attempts, got %d", attemptCount)
	}

	if result["status"] != "ok" {
		t.Errorf("expected status 'ok', got '%s'", result["status"])
	}
}

func TestClient_Do_NoRetryForPost(t *testing.T) {
	attemptCount := 0

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		attemptCount++
		w.WriteHeader(http.StatusServiceUnavailable)
	}))
	defer server.Close()

	client := NewClient(server.URL, "test-token", &http.Client{Timeout: 5 * time.Second}, nil)

	err := client.Do(context.Background(), &Request{
		Method: "POST",
		Path:   "/test",
		Body:   map[string]string{"key": "value"},
	}, nil)

	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if attemptCount != 1 {
		t.Errorf("expected 1 attempt (no retry for POST), got %d", attemptCount)
	}
}
