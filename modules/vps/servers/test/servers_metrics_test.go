package servers_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	cloudsdk "github.com/Zillaforge/cloud-sdk"
	"github.com/Zillaforge/cloud-sdk/models/vps/servers"
)

// TestServersMetrics_Success verifies successful metrics retrieval
func TestServersMetrics_Success(t *testing.T) {
	mockResponse := []map[string]interface{}{
		{
			"name": "cpu",
			"measures": []map[string]interface{}{
				{"granularity": 3600, "timestamp": 1704067200, "value": 45.5},
				{"granularity": 3600, "timestamp": 1704070800, "value": 50.2},
			},
		},
	}

	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/vps/api/v1/project/proj-123/servers/svr-1/metric" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		// Verify query parameters
		if r.URL.Query().Get("type") != "cpu" {
			t.Errorf("expected type=cpu")
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(mockResponse)
	}))
	defer mockServer.Close()

	client, err := cloudsdk.New(mockServer.URL, "test-token")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}
	vpsClient := client.Project("proj-123").VPS()

	req := &servers.ServerMetricsRequest{
		Type:        "cpu",
		Start:       1704067200,
		Granularity: 3600,
	}

	metrics, err := vpsClient.Servers().Metrics(context.Background(), "svr-1", req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(*metrics) == 0 {
		t.Fatal("expected metrics data, got empty response")
	}

	firstMetric := (*metrics)[0]
	if firstMetric.Name != "cpu" {
		t.Errorf("expected metric name 'cpu', got '%s'", firstMetric.Name)
	}
	if len(firstMetric.Measures) != 2 {
		t.Errorf("expected 2 data points, got %d", len(firstMetric.Measures))
	}
}

// TestServersMetrics_Errors verifies error handling for metrics retrieval
func TestServersMetrics_Errors(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
	}{
		{"invalid range - 400", http.StatusBadRequest},
		{"unauthorized - 401", http.StatusUnauthorized},
		{"not found - 404", http.StatusNotFound},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tt.statusCode)
				_ = json.NewEncoder(w).Encode(map[string]interface{}{
					"message":    "Error",
					"error_code": tt.statusCode,
				})
			}))
			defer mockServer.Close()

			client, err := cloudsdk.New(mockServer.URL, "test-token")
			if err != nil {
				t.Fatalf("failed to create client: %v", err)
			}
			vpsClient := client.Project("proj-123").VPS()

			req := &servers.ServerMetricsRequest{Type: "cpu"}
			_, err = vpsClient.Servers().Metrics(context.Background(), "svr-1", req)
			if err == nil {
				t.Fatal("expected error, got nil")
			}
		})
	}
}
