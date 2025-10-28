package servers

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	internalhttp "github.com/Zillaforge/cloud-sdk/internal/http"
	"github.com/Zillaforge/cloud-sdk/models/vps/servers"
)

func TestVolumesClient_List(t *testing.T) {
	mockResponse := []*servers.VolumeAttachment{
		{
			VolumeID: "vol-1",
			ServerID: "svr-1",
			Device:   "/dev/vdb",
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Errorf("expected GET, got %s", r.Method)
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(mockResponse)
	}))
	defer server.Close()

	baseClient := internalhttp.NewClient(server.URL, "test-token", &http.Client{Timeout: 5 * time.Second}, nil)
	volumesClient := &VolumesClient{
		baseClient: baseClient,
		projectID:  "proj-123",
		serverID:   "svr-1",
	}

	result, err := volumesClient.List(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(result) != len(mockResponse) {
		t.Errorf("expected %d volumes, got %d", len(mockResponse), len(result))
	}

	if result[0].VolumeID != mockResponse[0].VolumeID {
		t.Errorf("expected volume ID %s, got %s", mockResponse[0].VolumeID, result[0].VolumeID)
	}
}

func TestVolumesClient_Attach(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("expected POST, got %s", r.Method)
		}

		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	baseClient := internalhttp.NewClient(server.URL, "test-token", &http.Client{Timeout: 5 * time.Second}, nil)
	volumesClient := &VolumesClient{
		baseClient: baseClient,
		projectID:  "proj-123",
		serverID:   "svr-1",
	}

	err := volumesClient.Attach(context.Background(), "vol-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestVolumesClient_Detach(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "DELETE" {
			t.Errorf("expected DELETE, got %s", r.Method)
		}

		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	baseClient := internalhttp.NewClient(server.URL, "test-token", &http.Client{Timeout: 5 * time.Second}, nil)
	volumesClient := &VolumesClient{
		baseClient: baseClient,
		projectID:  "proj-123",
		serverID:   "svr-1",
	}

	err := volumesClient.Detach(context.Background(), "vol-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
