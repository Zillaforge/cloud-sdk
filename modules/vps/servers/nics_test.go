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

func TestNICsClient_List(t *testing.T) {
	mockResponse := []*servers.ServerNIC{
		{
			ID:        "nic-1",
			NetworkID: "net-1",
			FixedIPs:  []string{"10.0.0.10"},
			MACAddr:   "fa:16:3e:00:00:01",
			SGIDs:     []string{"sg-1"},
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
	nicsClient := &NICsClient{
		baseClient: baseClient,
		projectID:  "proj-123",
		serverID:   "svr-1",
	}

	result, err := nicsClient.List(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(result) != len(mockResponse) {
		t.Errorf("expected %d NICs, got %d", len(mockResponse), len(result))
	}
}

func TestNICsClient_Add(t *testing.T) {
	mockResponse := &servers.ServerNIC{
		ID:        "nic-new",
		NetworkID: "net-1",
		FixedIPs:  []string{"10.0.0.20"},
		MACAddr:   "fa:16:3e:00:00:02",
		SGIDs:     []string{"sg-1"},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("expected POST, got %s", r.Method)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(mockResponse)
	}))
	defer server.Close()

	baseClient := internalhttp.NewClient(server.URL, "test-token", &http.Client{Timeout: 5 * time.Second}, nil)
	nicsClient := &NICsClient{
		baseClient: baseClient,
		projectID:  "proj-123",
		serverID:   "svr-1",
	}

	req := &servers.ServerNICCreateRequest{
		NetworkID: "net-1",
		SGIDs:     []string{"sg-1"},
	}

	result, err := nicsClient.Add(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.ID != mockResponse.ID {
		t.Errorf("expected NIC ID %s, got %s", mockResponse.ID, result.ID)
	}
}

func TestNICsClient_Update(t *testing.T) {
	mockResponse := &servers.ServerNIC{
		ID:        "nic-1",
		NetworkID: "net-1",
		FixedIPs:  []string{"10.0.0.10"},
		MACAddr:   "fa:16:3e:00:00:01",
		SGIDs:     []string{"sg-1", "sg-2"},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "PUT" {
			t.Errorf("expected PUT, got %s", r.Method)
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(mockResponse)
	}))
	defer server.Close()

	baseClient := internalhttp.NewClient(server.URL, "test-token", &http.Client{Timeout: 5 * time.Second}, nil)
	nicsClient := &NICsClient{
		baseClient: baseClient,
		projectID:  "proj-123",
		serverID:   "svr-1",
	}

	req := &servers.ServerNICUpdateRequest{
		SGIDs: []string{"sg-1", "sg-2"},
	}

	result, err := nicsClient.Update(context.Background(), "nic-1", req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(result.SGIDs) != 2 {
		t.Errorf("expected 2 security groups, got %d", len(result.SGIDs))
	}
}

func TestNICsClient_Delete(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "DELETE" {
			t.Errorf("expected DELETE, got %s", r.Method)
		}

		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	baseClient := internalhttp.NewClient(server.URL, "test-token", &http.Client{Timeout: 5 * time.Second}, nil)
	nicsClient := &NICsClient{
		baseClient: baseClient,
		projectID:  "proj-123",
		serverID:   "svr-1",
	}

	err := nicsClient.Delete(context.Background(), "nic-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestNICsClient_AssociateFloatingIP(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("expected POST, got %s", r.Method)
		}

		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	baseClient := internalhttp.NewClient(server.URL, "test-token", &http.Client{Timeout: 5 * time.Second}, nil)
	nicsClient := &NICsClient{
		baseClient: baseClient,
		projectID:  "proj-123",
		serverID:   "svr-1",
	}

	req := &servers.FloatingIPAssociateRequest{
		FloatingIPID: "fip-1",
	}

	err := nicsClient.AssociateFloatingIP(context.Background(), "nic-1", req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
