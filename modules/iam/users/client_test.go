package users_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	internalhttp "github.com/Zillaforge/cloud-sdk/internal/http"
	"github.com/Zillaforge/cloud-sdk/models/iam/users"
	usersClient "github.com/Zillaforge/cloud-sdk/modules/iam/users"
)

func TestClient_Get_Success(t *testing.T) {
	// Mock HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request
		if r.Method != http.MethodGet {
			t.Errorf("Expected GET request, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/user" {
			t.Errorf("Expected path /api/v1/user, got %s", r.URL.Path)
		}

		// Verify authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader != "Bearer test-token" {
			t.Errorf("Expected Authorization: Bearer test-token, got %s", authHeader)
		}

		// Return mock user response
		response := users.User{
			UserID:      "test-user-id",
			Account:     "test@example.com",
			DisplayName: "Test User",
			Description: "Test description",
			Extra:       map[string]interface{}{},
			Namespace:   "test.com",
			Email:       "test@example.com",
			Frozen:      false,
			MFA:         true,
			CreatedAt:   "2025-01-01T00:00:00Z",
			UpdatedAt:   "2025-01-01T00:00:00Z",
			LastLoginAt: "2025-01-01T00:00:00Z",
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(response); err != nil {
			t.Fatalf("Failed to encode response: %v", err)
		}
	}))
	defer server.Close()

	// Create HTTP client with timeout
	httpClient := &http.Client{
		Timeout: 5 * time.Second,
	}

	// Create client
	baseClient := internalhttp.NewClient(server.URL, "test-token", httpClient, nil)
	client := usersClient.NewClient(baseClient, "/api/v1/")

	// Call Get with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	user, err := client.Get(ctx)

	// Verify result
	if err != nil {
		t.Fatalf("Get() returned error: %v", err)
	}
	if user == nil {
		t.Fatal("Get() returned nil user")
	}
	if user.UserID != "test-user-id" {
		t.Errorf("UserID = %v, want %v", user.UserID, "test-user-id")
	}
	if user.DisplayName != "Test User" {
		t.Errorf("DisplayName = %v, want %v", user.DisplayName, "Test User")
	}
	if user.MFA != true {
		t.Errorf("MFA = %v, want %v", user.MFA, true)
	}
}

func TestClient_Get_Forbidden(t *testing.T) {
	// Mock HTTP server returning 403
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusForbidden)
		if err := json.NewEncoder(w).Encode(map[string]interface{}{
			"errorCode": 403,
			"message":   "Invalid or expired token",
		}); err != nil {
			t.Fatalf("Failed to encode response: %v", err)
		}
	}))
	defer server.Close()

	// Create HTTP client with timeout
	httpClient := &http.Client{
		Timeout: 5 * time.Second,
	}

	// Create client
	baseClient := internalhttp.NewClient(server.URL, "invalid-token", httpClient, nil)
	client := usersClient.NewClient(baseClient, "/api/v1/")

	// Call Get
	ctx := context.Background()
	user, err := client.Get(ctx)

	// Verify error
	if err == nil {
		t.Fatal("Get() should return error for 403 response")
	}
	if user != nil {
		t.Error("Get() should return nil user on error")
	}
}

func TestClient_Get_ContextTimeout(t *testing.T) {
	// Mock HTTP server with delay
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		// Delay to trigger timeout
		time.Sleep(200 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	// Create HTTP client with long timeout (so context timeout triggers first)
	httpClient := &http.Client{
		Timeout: 10 * time.Second,
	}

	// Create client
	baseClient := internalhttp.NewClient(server.URL, "test-token", httpClient, nil)
	client := usersClient.NewClient(baseClient, "/api/v1/")

	// Call Get with short timeout
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	user, err := client.Get(ctx)

	// Verify timeout error
	if err == nil {
		t.Fatal("Get() should return error for context timeout")
	}
	if user != nil {
		t.Error("Get() should return nil user on timeout")
	}
}

func TestClient_Get_ContextCanceled(t *testing.T) {
	// Mock HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		// Delay to allow cancellation
		time.Sleep(100 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	// Create HTTP client with long timeout (so context cancellation triggers first)
	httpClient := &http.Client{
		Timeout: 10 * time.Second,
	}

	// Create client
	baseClient := internalhttp.NewClient(server.URL, "test-token", httpClient, nil)
	client := usersClient.NewClient(baseClient, "/api/v1/")

	// Call Get with cancelable context
	ctx, cancel := context.WithCancel(context.Background())

	// Cancel context after 50ms
	go func() {
		time.Sleep(50 * time.Millisecond)
		cancel()
	}()

	user, err := client.Get(ctx)

	// Verify cancellation error
	if err == nil {
		t.Fatal("Get() should return error for context cancellation")
	}
	if user != nil {
		t.Error("Get() should return nil user on cancellation")
	}
}
