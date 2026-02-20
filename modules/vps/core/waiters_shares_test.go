package vps

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	internalhttp "github.com/Zillaforge/cloud-sdk/internal/http"
	"github.com/Zillaforge/cloud-sdk/internal/waiter"
	sharesmodels "github.com/Zillaforge/cloud-sdk/models/vps/shares"
	"github.com/Zillaforge/cloud-sdk/modules/vps/shares"
)

// shareResponse builds a standard share mock response payload.
func shareResponse(status string) map[string]interface{} {
	return map[string]interface{}{
		"id":            "share-test-1",
		"name":          "test-share",
		"description":   "test share",
		"size":          10,
		"path":          "10.0.0.1:/share-test-1",
		"status":        status,
		"status_reason": "",
		"network":       map[string]interface{}{"id": "net-1", "name": "test-net"},
		"network_id":    "net-1",
		"project":       map[string]interface{}{"id": "proj-1", "name": "test-project"},
		"project_id":    "proj-1",
		"user":          map[string]interface{}{"id": "user-1", "name": "test-user"},
		"user_id":       "user-1",
		"namespace":     "test-namespace",
		"createdAt":     "2025-10-28T00:00:00Z",
		"updatedAt":     "2025-10-28T00:00:00Z",
	}
}

// TestWaitForShareStatus_Success verifies waiting for a share to reach target status.
func TestWaitForShareStatus_Success(t *testing.T) {
	tests := []struct {
		name         string
		targetStatus sharesmodels.ShareStatus
		statusFlow   []string // Sequence of statuses returned by mock server
	}{
		{
			name:         "wait for AVAILABLE from CREATING",
			targetStatus: sharesmodels.ShareStatusAvailable,
			statusFlow:   []string{"creating", "creating", "available"},
		},
		{
			name:         "wait for DELETED from DELETING",
			targetStatus: sharesmodels.ShareStatusDeleted,
			statusFlow:   []string{"deleting", "deleted"},
		},
		{
			name:         "already at target status",
			targetStatus: sharesmodels.ShareStatusAvailable,
			statusFlow:   []string{"available"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			currentStatusIndex := 0

			mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != "GET" {
					t.Errorf("expected GET, got %s", r.Method)
				}

				status := tt.statusFlow[currentStatusIndex]
				if currentStatusIndex < len(tt.statusFlow)-1 {
					currentStatusIndex++
				}

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				if err := json.NewEncoder(w).Encode(shareResponse(status)); err != nil {
					t.Fatalf("failed to encode response: %v", err)
				}
			}))
			defer mockServer.Close()

			httpClient := internalhttp.NewClient(mockServer.URL, "test-token", &http.Client{}, nil)
			shareClient := shares.NewClient(httpClient, "proj-1")

			ctx := context.Background()
			err := WaitForShareStatus(ctx, ShareWaiterConfig{
				Client:       shareClient,
				ShareID:      "share-test-1",
				TargetStatus: tt.targetStatus,
				WaiterOptions: []waiter.Option{
					waiter.WithInterval(10 * time.Millisecond),
					waiter.WithMaxWait(2 * time.Second),
				},
			})

			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}

// TestWaitForShareStatus_ErrorState verifies handling when a share enters ERROR state.
func TestWaitForShareStatus_ErrorState(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(shareResponse("error")); err != nil {
			t.Fatalf("failed to encode response: %v", err)
		}
	}))
	defer mockServer.Close()

	httpClient := internalhttp.NewClient(mockServer.URL, "test-token", &http.Client{}, nil)
	shareClient := shares.NewClient(httpClient, "proj-1")

	ctx := context.Background()
	err := WaitForShareStatus(ctx, ShareWaiterConfig{
		Client:       shareClient,
		ShareID:      "share-test-1",
		TargetStatus: sharesmodels.ShareStatusAvailable,
		WaiterOptions: []waiter.Option{
			waiter.WithInterval(10 * time.Millisecond),
			waiter.WithMaxWait(1 * time.Second),
		},
	})

	if err == nil {
		t.Fatal("expected error when share enters ERROR state, got nil")
	}

	expectedMsg := "share entered ERROR state"
	if err.Error()[:len(expectedMsg)] != expectedMsg {
		t.Errorf("expected error message to start with '%s', got '%s'", expectedMsg, err.Error())
	}
}

// TestWaitForShareStatus_Timeout verifies timeout handling.
func TestWaitForShareStatus_Timeout(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		// Always return creating status
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(shareResponse("creating")); err != nil {
			t.Fatalf("failed to encode response: %v", err)
		}
	}))
	defer mockServer.Close()

	httpClient := &http.Client{Timeout: 30 * time.Second}
	baseClient := internalhttp.NewClient(mockServer.URL, "test-token", httpClient, nil)
	shareClient := shares.NewClient(baseClient, "proj-1")

	ctx := context.Background()
	err := WaitForShareStatus(ctx, ShareWaiterConfig{
		Client:       shareClient,
		ShareID:      "share-test-1",
		TargetStatus: sharesmodels.ShareStatusAvailable,
		WaiterOptions: []waiter.Option{
			waiter.WithInterval(50 * time.Millisecond),
			waiter.WithMaxWait(200 * time.Millisecond),
		},
	})

	if err == nil {
		t.Fatal("expected timeout error, got nil")
	}

	if err != waiter.ErrWaitTimeout {
		t.Errorf("expected ErrWaitTimeout, got: %v", err)
	}
}

// TestWaitForShareStatus_ContextCancellation verifies context cancellation.
func TestWaitForShareStatus_ContextCancellation(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		// Always return creating status
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(shareResponse("creating")); err != nil {
			t.Fatalf("failed to encode response: %v", err)
		}
	}))
	defer mockServer.Close()

	httpClient := internalhttp.NewClient(mockServer.URL, "test-token", &http.Client{}, nil)
	shareClient := shares.NewClient(httpClient, "proj-1")

	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		time.Sleep(50 * time.Millisecond)
		cancel()
	}()

	err := WaitForShareStatus(ctx, ShareWaiterConfig{
		Client:       shareClient,
		ShareID:      "share-test-1",
		TargetStatus: sharesmodels.ShareStatusAvailable,
		WaiterOptions: []waiter.Option{
			waiter.WithInterval(10 * time.Millisecond),
			waiter.WithMaxWait(5 * time.Second),
		},
	})

	if err == nil {
		t.Fatal("expected context cancellation error, got nil")
	}

	if !errors.Is(err, context.Canceled) {
		t.Errorf("expected context.Canceled, got: %v", err)
	}
}

// TestWaitForShareStatus_ValidationErrors verifies input validation.
func TestWaitForShareStatus_ValidationErrors(t *testing.T) {
	httpClient := internalhttp.NewClient("http://example.com", "test-token", &http.Client{}, nil)
	shareClient := shares.NewClient(httpClient, "proj-1")

	tests := []struct {
		name        string
		cfg         ShareWaiterConfig
		expectedErr string
	}{
		{
			name: "missing client",
			cfg: ShareWaiterConfig{
				Client:       nil,
				ShareID:      "share-1",
				TargetStatus: sharesmodels.ShareStatusAvailable,
			},
			expectedErr: "share client is required",
		},
		{
			name: "missing share ID",
			cfg: ShareWaiterConfig{
				Client:       shareClient,
				ShareID:      "",
				TargetStatus: sharesmodels.ShareStatusAvailable,
			},
			expectedErr: "share ID is required",
		},
		{
			name: "missing target status",
			cfg: ShareWaiterConfig{
				Client:       shareClient,
				ShareID:      "share-1",
				TargetStatus: "",
			},
			expectedErr: "target status is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			err := WaitForShareStatus(ctx, tt.cfg)

			if err == nil {
				t.Fatal("expected error, got nil")
			}

			if err.Error() != tt.expectedErr {
				t.Errorf("expected error '%s', got '%s'", tt.expectedErr, err.Error())
			}
		})
	}
}

// TestWaitForShareAvailable verifies the convenience function for waiting AVAILABLE status.
func TestWaitForShareAvailable(t *testing.T) {
	currentStatusIndex := 0
	statusFlow := []string{"creating", "available"}

	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Errorf("expected GET, got %s", r.Method)
		}

		status := statusFlow[currentStatusIndex]
		if currentStatusIndex < len(statusFlow)-1 {
			currentStatusIndex++
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(shareResponse(status)); err != nil {
			t.Fatalf("failed to encode response: %v", err)
		}
	}))
	defer mockServer.Close()

	httpClient := internalhttp.NewClient(mockServer.URL, "test-token", &http.Client{}, nil)
	shareClient := shares.NewClient(httpClient, "proj-1")

	ctx := context.Background()
	err := WaitForShareAvailable(ctx, shareClient, "share-test-1",
		waiter.WithInterval(10*time.Millisecond),
		waiter.WithMaxWait(2*time.Second),
	)

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}
