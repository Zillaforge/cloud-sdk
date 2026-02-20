package vps

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	internalhttp "github.com/Zillaforge/cloud-sdk/internal/http"
	"github.com/Zillaforge/cloud-sdk/internal/waiter"
	lbmodels "github.com/Zillaforge/cloud-sdk/models/vps/loadbalancers"
	"github.com/Zillaforge/cloud-sdk/modules/vps/loadbalancers"
)

func lbResponse(status string) map[string]interface{} {
	return map[string]interface{}{
		"id":         "lb-test-1",
		"name":       "test-lb",
		"status":     status,
		"network_id": "net-1",
		"network":    map[string]string{"id": "net-1", "name": "test-net"},
		"project_id": "proj-1",
		"project":    map[string]string{"id": "proj-1", "name": "test-proj"},
		"user_id":    "user-1",
		"user":       map[string]string{"id": "user-1", "name": "test-user"},
	}
}

func TestWaitForLoadBalancerStatus_Success(t *testing.T) {
	tests := []struct {
		name         string
		targetStatus lbmodels.LoadBalancerStatus
		statusFlow   []string
	}{
		{
			name:         "wait for ACTIVE from PENDING_CREATE",
			targetStatus: lbmodels.LoadBalancerStatusActive,
			statusFlow:   []string{"PENDING_CREATE", "PENDING_CREATE", "ACTIVE"},
		},
		{
			name:         "already ACTIVE",
			targetStatus: lbmodels.LoadBalancerStatusActive,
			statusFlow:   []string{"ACTIVE"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			currentStatusIndex := 0

			mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				status := tt.statusFlow[currentStatusIndex]
				if currentStatusIndex < len(tt.statusFlow)-1 {
					currentStatusIndex++
				}
				w.Header().Set("Content-Type", "application/json")
				_ = json.NewEncoder(w).Encode(lbResponse(status))
			}))
			defer mockServer.Close()

			httpClient := internalhttp.NewClient(mockServer.URL, "test-token", &http.Client{}, nil)
			lbClient := loadbalancers.NewClient(httpClient, "proj-1")

			ctx := context.Background()
			err := WaitForLoadBalancerStatus(ctx, LoadBalancerWaiterConfig{
				Client:         lbClient,
				LoadBalancerID: "lb-test-1",
				TargetStatus:   tt.targetStatus,
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

func TestWaitForLoadBalancerStatus_ErrorState(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(lbResponse("ERROR"))
	}))
	defer mockServer.Close()

	httpClient := internalhttp.NewClient(mockServer.URL, "test-token", &http.Client{}, nil)
	lbClient := loadbalancers.NewClient(httpClient, "proj-1")

	ctx := context.Background()
	err := WaitForLoadBalancerStatus(ctx, LoadBalancerWaiterConfig{
		Client:         lbClient,
		LoadBalancerID: "lb-test-1",
		TargetStatus:   lbmodels.LoadBalancerStatusActive,
		WaiterOptions: []waiter.Option{
			waiter.WithInterval(10 * time.Millisecond),
			waiter.WithMaxWait(1 * time.Second),
		},
	})

	if err == nil {
		t.Fatal("expected error for ERROR state, got nil")
	}
}

func TestWaitForLoadBalancerStatus_ValidationErrors(t *testing.T) {
	httpClient := internalhttp.NewClient("https://api.example.com", "test-token", &http.Client{}, nil)
	lbClient := loadbalancers.NewClient(httpClient, "proj-1")

	tests := []struct {
		name string
		cfg  LoadBalancerWaiterConfig
	}{
		{name: "nil client", cfg: LoadBalancerWaiterConfig{Client: nil, LoadBalancerID: "lb-1", TargetStatus: lbmodels.LoadBalancerStatusActive}},
		{name: "empty ID", cfg: LoadBalancerWaiterConfig{Client: lbClient, LoadBalancerID: "", TargetStatus: lbmodels.LoadBalancerStatusActive}},
		{name: "empty status", cfg: LoadBalancerWaiterConfig{Client: lbClient, LoadBalancerID: "lb-1", TargetStatus: ""}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := WaitForLoadBalancerStatus(context.Background(), tt.cfg)
			if err == nil {
				t.Fatal("expected validation error, got nil")
			}
		})
	}
}

func TestWaitForLoadBalancerActive(t *testing.T) {
	currentStatusIndex := 0
	statusFlow := []string{"PENDING_CREATE", "ACTIVE"}

	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		status := statusFlow[currentStatusIndex]
		if currentStatusIndex < len(statusFlow)-1 {
			currentStatusIndex++
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(lbResponse(status))
	}))
	defer mockServer.Close()

	httpClient := internalhttp.NewClient(mockServer.URL, "test-token", &http.Client{}, nil)
	lbClient := loadbalancers.NewClient(httpClient, "proj-1")

	ctx := context.Background()
	err := WaitForLoadBalancerActive(ctx, lbClient, "lb-test-1",
		waiter.WithInterval(10*time.Millisecond),
		waiter.WithMaxWait(2*time.Second),
	)

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestWaitForLoadBalancerDeleted(t *testing.T) {
	callCount := 0

	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		callCount++
		if callCount <= 2 {
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(lbResponse("PENDING_DELETE"))
		} else {
			w.WriteHeader(http.StatusNotFound)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": "not found"})
		}
	}))
	defer mockServer.Close()

	httpClient := internalhttp.NewClient(mockServer.URL, "test-token", &http.Client{}, nil)
	lbClient := loadbalancers.NewClient(httpClient, "proj-1")

	ctx := context.Background()
	err := WaitForLoadBalancerDeleted(ctx, lbClient, "lb-test-1",
		waiter.WithInterval(10*time.Millisecond),
		waiter.WithMaxWait(2*time.Second),
	)

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}
