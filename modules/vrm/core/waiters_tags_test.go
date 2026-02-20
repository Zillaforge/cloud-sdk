package vrm

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	internalhttp "github.com/Zillaforge/cloud-sdk/internal/http"
	"github.com/Zillaforge/cloud-sdk/internal/waiter"
	commonmodels "github.com/Zillaforge/cloud-sdk/models/vrm/common"
	"github.com/Zillaforge/cloud-sdk/modules/vrm/tags"
)

// TestWaitForTagStatus_Success verifies waiting for a tag to reach target status.
func TestWaitForTagStatus_Success(t *testing.T) {
	tests := []struct {
		name         string
		targetStatus commonmodels.TagStatus
		statusFlow   []string // Sequence of statuses returned by mock tag
	}{
		{
			name:         "wait for ACTIVE from CREATING",
			targetStatus: commonmodels.TagStatusActive,
			statusFlow:   []string{"creating", "creating", "active"},
		},
		{
			name:         "wait for AVAILABLE from SAVING",
			targetStatus: commonmodels.TagStatusAvailable,
			statusFlow:   []string{"saving", "saving", "available"},
		},
		{
			name:         "already at target status",
			targetStatus: commonmodels.TagStatusActive,
			statusFlow:   []string{"active"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			currentStatusIndex := 0

			// Create mock server
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != "GET" {
					t.Errorf("expected GET request, got %s", r.Method)
					w.WriteHeader(http.StatusMethodNotAllowed)
					return
				}

				// Return tag with current status in flow
				status := tt.statusFlow[currentStatusIndex]
				if currentStatusIndex < len(tt.statusFlow)-1 {
					currentStatusIndex++
				}

				tag := commonmodels.Tag{
					ID:     "tag-123",
					Name:   "test-tag",
					Status: commonmodels.TagStatus(status),
				}

				w.Header().Set("Content-Type", "application/json")
				if err := json.NewEncoder(w).Encode(tag); err != nil {
					t.Fatalf("failed to encode response: %v", err)
				}
			}))
			defer server.Close()

			// Create client
			baseClient := internalhttp.NewClient(server.URL, "token", &http.Client{}, nil)
			client := tags.NewClient(baseClient, "project-123", "/api/v1/project/project-123")

			// Test waiter
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()

			err := WaitForTagStatus(ctx, TagWaiterConfig{
				Client:       client,
				TagID:        "tag-123",
				TargetStatus: tt.targetStatus,
				WaiterOptions: []waiter.Option{
					waiter.WithInterval(10 * time.Millisecond), // Fast for testing
					waiter.WithMaxWait(1 * time.Second),
				},
			})

			if err != nil {
				t.Errorf("WaitForTagStatus() error = %v", err)
			}
		})
	}
}

// TestWaitForTagStatus_Error verifies error handling in tag status waiting.
func TestWaitForTagStatus_Error(t *testing.T) {
	tests := []struct {
		name         string
		targetStatus commonmodels.TagStatus
		statusFlow   []string
		expectError  bool
	}{
		{
			name:         "tag enters ERROR state",
			targetStatus: commonmodels.TagStatusActive,
			statusFlow:   []string{"creating", "error"},
			expectError:  true,
		},
		{
			name:         "invalid client",
			targetStatus: commonmodels.TagStatusActive,
			statusFlow:   []string{},
			expectError:  true,
		},
		{
			name:         "empty tag ID",
			targetStatus: commonmodels.TagStatusActive,
			statusFlow:   []string{},
			expectError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "invalid client" {
				err := WaitForTagStatus(context.Background(), TagWaiterConfig{
					Client:       nil,
					TagID:        "tag-123",
					TargetStatus: tt.targetStatus,
				})
				if !tt.expectError || err == nil {
					t.Errorf("expected error for invalid client, got %v", err)
				}
				return
			}

			if tt.name == "empty tag ID" {
				server := httptest.NewServer(http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {}))
				defer server.Close()

				baseClient := internalhttp.NewClient(server.URL, "token", &http.Client{}, nil)
				client := tags.NewClient(baseClient, "project-123", "/api/v1/project/project-123")

				err := WaitForTagStatus(context.Background(), TagWaiterConfig{
					Client:       client,
					TagID:        "",
					TargetStatus: tt.targetStatus,
				})
				if !tt.expectError || err == nil {
					t.Errorf("expected error for empty tag ID, got %v", err)
				}
				return
			}

			currentStatusIndex := 0

			// Create mock server
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				status := tt.statusFlow[currentStatusIndex]
				if currentStatusIndex < len(tt.statusFlow)-1 {
					currentStatusIndex++
				}

				tag := commonmodels.Tag{
					ID:     "tag-123",
					Name:   "test-tag",
					Status: commonmodels.TagStatus(status),
				}

				w.Header().Set("Content-Type", "application/json")
				if err := json.NewEncoder(w).Encode(tag); err != nil {
					t.Fatalf("failed to encode response: %v", err)
				}
			}))
			defer server.Close()

			// Create client
			baseClient := internalhttp.NewClient(server.URL, "token", &http.Client{}, nil)
			client := tags.NewClient(baseClient, "project-123", "/api/v1/project/project-123")

			// Test waiter
			ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
			defer cancel()

			err := WaitForTagStatus(ctx, TagWaiterConfig{
				Client:       client,
				TagID:        "tag-123",
				TargetStatus: tt.targetStatus,
				WaiterOptions: []waiter.Option{
					waiter.WithInterval(10 * time.Millisecond),
					waiter.WithMaxWait(500 * time.Millisecond),
				},
			})

			if tt.expectError && err == nil {
				t.Errorf("expected error, but got none")
			}
			if !tt.expectError && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}

// Test convenience functions
func TestConvenienceFunctions(t *testing.T) {
	// Test WaitForTagActive
	t.Run("WaitForTagActive", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			tag := commonmodels.Tag{
				ID:     "tag-123",
				Name:   "test-tag",
				Status: commonmodels.TagStatusActive,
			}
			w.Header().Set("Content-Type", "application/json")
			if err := json.NewEncoder(w).Encode(tag); err != nil {
				t.Fatalf("failed to encode response: %v", err)
			}
		}))
		defer server.Close()

		baseClient := internalhttp.NewClient(server.URL, "token", &http.Client{}, nil)
		client := tags.NewClient(baseClient, "project-123", "/api/v1/project/project-123")

		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		defer cancel()

		err := WaitForTagActive(ctx, client, "tag-123",
			waiter.WithInterval(10*time.Millisecond),
			waiter.WithMaxWait(500*time.Millisecond),
		)
		if err != nil {
			t.Errorf("WaitForTagActive() error = %v", err)
		}
	})

	// Test WaitForTagAvailable
	t.Run("WaitForTagAvailable", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			tag := commonmodels.Tag{
				ID:     "tag-123",
				Name:   "test-tag",
				Status: commonmodels.TagStatusAvailable,
			}
			w.Header().Set("Content-Type", "application/json")
			if err := json.NewEncoder(w).Encode(tag); err != nil {
				t.Fatalf("failed to encode response: %v", err)
			}
		}))
		defer server.Close()

		baseClient := internalhttp.NewClient(server.URL, "token", &http.Client{}, nil)
		client := tags.NewClient(baseClient, "project-123", "/api/v1/project/project-123")

		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		defer cancel()

		err := WaitForTagAvailable(ctx, client, "tag-123",
			waiter.WithInterval(10*time.Millisecond),
			waiter.WithMaxWait(500*time.Millisecond),
		)
		if err != nil {
			t.Errorf("WaitForTagAvailable() error = %v", err)
		}
	})
}
