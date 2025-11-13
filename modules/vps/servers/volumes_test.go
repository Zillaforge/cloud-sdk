package servers

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	internalhttp "github.com/Zillaforge/cloud-sdk/internal/http"
	"github.com/Zillaforge/cloud-sdk/models/vps/common"
	"github.com/Zillaforge/cloud-sdk/models/vps/servers"
	"github.com/Zillaforge/cloud-sdk/models/vps/volumes"
)

func TestVolumesClient_List(t *testing.T) {
	mockResponse := &servers.ServerVolumesResponse{
		Disks: []*servers.ServerVolume{
			{
				System:   true,
				VolumeID: "0a92879e-27ff-4d39-a156-84ab63f3d581",
				Device:   "/dev/vda",
				Volume:   nil,
			},
			{
				System:   false,
				VolumeID: "857ca66b-0906-4902-98b1-a2a917fd8321",
				Device:   "/dev/vdb",
				Volume: &volumes.Volume{
					ID:           "857ca66b-0906-4902-98b1-a2a917fd8321",
					Name:         "test-vol",
					Description:  "",
					Size:         1,
					Type:         "__DEFAULT__",
					Status:       "attaching",
					StatusReason: "",
					Attachments:  []common.IDName{},
					ProjectID:    "91457b61-0b92-4aa8-b136-b03d88f04946",
					UserID:       "4990ccdb-a9b1-49e5-91df-67c921601d81",
					Namespace:    "public",
					CreatedAt:    &time.Time{},
					UpdatedAt:    &time.Time{},
				},
			},
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

	if len(result) != len(mockResponse.Disks) {
		t.Errorf("expected %d disks, got %d", len(mockResponse.Disks), len(result))
	}

	if result[0].VolumeID != mockResponse.Disks[0].VolumeID {
		t.Errorf("expected volume ID %s, got %s", mockResponse.Disks[0].VolumeID, result[0].VolumeID)
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
