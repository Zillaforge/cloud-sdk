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

// TestServerNICs_List_Success verifies successful NIC listing
func TestServerNICs_List_Success(t *testing.T) {
	mockResponse := map[string]interface{}{
		"nics": []map[string]interface{}{
			{
				"id":         "nic-1",
				"network_id": "net-1",
				"addresses":  []string{"10.0.0.10"},
				"mac":        "fa:16:3e:00:00:01",
				"sg_ids":     []string{"sg-1"},
			},
			{
				"id":         "nic-2",
				"network_id": "net-2",
				"addresses":  []string{"10.0.1.10"},
				"mac":        "fa:16:3e:00:00:02",
				"sg_ids":     []string{"sg-1", "sg-2"},
			},
		},
	}

	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Errorf("expected GET, got %s", r.Method)
		}

		// Handle both GET server and GET NICs
		if r.URL.Path == "/vps/api/v1/project/proj-123/servers/svr-1" {
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

		if r.URL.Path == "/vps/api/v1/project/proj-123/servers/svr-1/nics" {
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(mockResponse)
			return
		}

		t.Errorf("unexpected path: %s", r.URL.Path)
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

	response, err := serverResource.NICs().List(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(response.NICs) != 2 {
		t.Errorf("expected 2 NICs, got %d", len(response.NICs))
	}
}

// TestServerNICs_Add_Success verifies successful NIC addition
func TestServerNICs_Add_Success(t *testing.T) {
	mockResponse := map[string]interface{}{
		"id":          "nic-new",
		"network_id":  "net-1",
		"fixed_ips":   []string{"10.0.0.20"},
		"mac_address": "fa:16:3e:00:00:03",
		"sg_ids":      []string{"sg-1"},
	}

	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" && r.URL.Path == "/vps/api/v1/project/proj-123/servers/svr-1/nics" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)
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
		t.Errorf("unexpected request: %s %s", r.Method, r.URL.Path)
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

	req := &servers.ServerNICCreateRequest{
		NetworkID: "net-1",
		SGIDs:     []string{"sg-1"},
	}

	nic, err := serverResource.NICs().Add(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if nic.ID != "nic-new" {
		t.Errorf("expected NIC ID 'nic-new', got '%s'", nic.ID)
	}
}

// TestServerNICs_Update_Success verifies successful NIC update
func TestServerNICs_Update_Success(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "PUT" && r.URL.Path == "/vps/api/v1/project/proj-123/servers/svr-1/nics/nic-1" {
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"id":          "nic-1",
				"network_id":  "net-1",
				"fixed_ips":   []string{"10.0.0.10"},
				"mac_address": "fa:16:3e:00:00:01",
				"sg_ids":      []string{"sg-1", "sg-2"},
			})
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

	req := &servers.ServerNICUpdateRequest{
		SGIDs: []string{"sg-1", "sg-2"},
	}

	_, err = serverResource.NICs().Update(context.Background(), "nic-1", req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

// TestServerNICs_Delete_Success verifies successful NIC deletion
func TestServerNICs_Delete_Success(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "DELETE" && r.URL.Path == "/vps/api/v1/project/proj-123/servers/svr-1/nics/nic-1" {
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

	err = serverResource.NICs().Delete(context.Background(), "nic-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

// TestServerNICs_AssociateFloatingIP_Success verifies successful floating IP association
func TestServerNICs_AssociateFloatingIP_Success(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" && r.URL.Path == "/vps/api/v1/project/proj-123/servers/svr-1/nics/nic-1/floatingip" {
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

	req := &servers.ServerNICAssociateFloatingIPRequest{
		FIPID: "fip-1",
	}

	_, err = serverResource.NICs().AssociateFloatingIP(context.Background(), "nic-1", req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
