package servers_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	cloudsdk "github.com/Zillaforge/cloud-sdk"
)

// TestServerVolumes_List_Success verifies successful volume listing
func TestServerVolumes_List_Success(t *testing.T) {
	mockResponse := []map[string]interface{}{
		{
			"volume_id": "vol-1",
			"server_id": "svr-1",
			"device":    "/dev/vdb",
		},
		{
			"volume_id": "vol-2",
			"server_id": "svr-1",
			"device":    "/dev/vdc",
		},
	}

	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" && r.URL.Path == "/vps/api/v1/project/proj-123/servers/svr-1/volumes" {
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(mockResponse)
			return
		}
		if r.Method == "GET" && r.URL.Path == "/vps/api/v1/project/proj-123/servers/svr-1" {
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"id":         "svr-1",
				"name":       "test-server",
				"status":     "ACTIVE",
				"flavor_id":  "flv-1",
				"image_id":   "img-1",
				"project_id": "proj-123",
			})
			return
		}
	}))
	defer mockServer.Close()

	client, err := cloudsdk.New(mockServer.URL, "test-token")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}
	vpsClient := client.Project("proj-123").VPS()

	serverResource, err := vpsClient.Servers().Get(context.Background(), "svr-1")
	if err != nil {
		t.Fatalf("failed to get server: %v", err)
	}

	volumes, err := serverResource.Volumes().List(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(volumes) != 2 {
		t.Errorf("expected 2 volumes, got %d", len(volumes))
	}
	if volumes[0].VolumeID != "vol-1" {
		t.Errorf("expected volume ID 'vol-1', got '%s'", volumes[0].VolumeID)
	}
}

// TestServerVolumes_Attach_Success verifies successful volume attachment
func TestServerVolumes_Attach_Success(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" && r.URL.Path == "/vps/api/v1/project/proj-123/servers/svr-1/volumes/vol-1" {
			w.WriteHeader(http.StatusOK)
			return
		}
		if r.Method == "GET" && r.URL.Path == "/vps/api/v1/project/proj-123/servers/svr-1" {
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"id":         "svr-1",
				"name":       "test-server",
				"status":     "ACTIVE",
				"flavor_id":  "flv-1",
				"image_id":   "img-1",
				"project_id": "proj-123",
			})
			return
		}
	}))
	defer mockServer.Close()

	client, err := cloudsdk.New(mockServer.URL, "test-token")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}
	vpsClient := client.Project("proj-123").VPS()

	serverResource, err := vpsClient.Servers().Get(context.Background(), "svr-1")
	if err != nil {
		t.Fatalf("failed to get server: %v", err)
	}

	err = serverResource.Volumes().Attach(context.Background(), "vol-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

// TestServerVolumes_Detach_Success verifies successful volume detachment
func TestServerVolumes_Detach_Success(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "DELETE" && r.URL.Path == "/vps/api/v1/project/proj-123/servers/svr-1/volumes/vol-1" {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		if r.Method == "GET" && r.URL.Path == "/vps/api/v1/project/proj-123/servers/svr-1" {
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"id":         "svr-1",
				"name":       "test-server",
				"status":     "ACTIVE",
				"flavor_id":  "flv-1",
				"image_id":   "img-1",
				"project_id": "proj-123",
			})
			return
		}
	}))
	defer mockServer.Close()

	client, err := cloudsdk.New(mockServer.URL, "test-token")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}
	vpsClient := client.Project("proj-123").VPS()

	serverResource, err := vpsClient.Servers().Get(context.Background(), "svr-1")
	if err != nil {
		t.Fatalf("failed to get server: %v", err)
	}

	err = serverResource.Volumes().Detach(context.Background(), "vol-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
