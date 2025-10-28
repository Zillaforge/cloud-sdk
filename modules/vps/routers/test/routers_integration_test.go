package test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	cloudsdk "github.com/Zillaforge/cloud-sdk"
	"github.com/Zillaforge/cloud-sdk/models/vps/routers"
)

// TestContract_RouterLifecycle tests the complete router lifecycle:
// 1. List routers (should be empty initially)
// 2. Create a router
// 3. Get the router
// 4. List routers (should contain created router)
// 5. Update the router
// 6. Set router state (enable/disable)
// 7. Associate a network
// 8. List router networks
// 9. Disassociate the network
// 10. Delete the router
// 11. Verify deletion
func TestContract_RouterLifecycle(t *testing.T) {
	var createdRouterID string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		// Step 1: Initial list (empty)
		case r.Method == http.MethodGet && r.URL.Path == "/vps/api/v1/project/proj-1/routers" && createdRouterID == "":
			w.WriteHeader(http.StatusOK)
			_ = json.NewEncoder(w).Encode(routers.RouterListResponse{
				Routers: []routers.Router{},
				Total:   0,
			})

		// Step 2: Create router
		case r.Method == http.MethodPost && r.URL.Path == "/vps/api/v1/project/proj-1/routers":
			var req routers.RouterCreateRequest
			_ = json.NewDecoder(r.Body).Decode(&req)
			createdRouterID = "router-lifecycle-123"
			w.WriteHeader(http.StatusCreated)
			_ = json.NewEncoder(w).Encode(routers.Router{
				ID:          createdRouterID,
				Name:        req.Name,
				Description: req.Description,
				State:       true,
				Status:      "ACTIVE",
				ProjectID:   "proj-1",
			})

		// Step 3 & 4 & 7: Get router (needed before accessing sub-resources)
		case r.Method == http.MethodGet && strings.Contains(r.URL.Path, "/routers/"+createdRouterID) && !strings.Contains(r.URL.Path, "/networks"):
			w.WriteHeader(http.StatusOK)
			_ = json.NewEncoder(w).Encode(routers.Router{
				ID:          createdRouterID,
				Name:        "lifecycle-router",
				Description: "Test lifecycle",
				State:       true,
				Status:      "ACTIVE",
				ProjectID:   "proj-1",
			})

		// Step 4: List routers (after creation)
		case r.Method == http.MethodGet && r.URL.Path == "/vps/api/v1/project/proj-1/routers" && createdRouterID != "":
			w.WriteHeader(http.StatusOK)
			_ = json.NewEncoder(w).Encode(routers.RouterListResponse{
				Routers: []routers.Router{
					{
						ID:          createdRouterID,
						Name:        "lifecycle-router",
						Description: "Test lifecycle",
						State:       true,
						Status:      "ACTIVE",
						ProjectID:   "proj-1",
					},
				},
				Total: 1,
			})

		// Step 5: Update router
		case r.Method == http.MethodPut && strings.Contains(r.URL.Path, createdRouterID):
			var req routers.RouterUpdateRequest
			_ = json.NewDecoder(r.Body).Decode(&req)
			w.WriteHeader(http.StatusOK)
			_ = json.NewEncoder(w).Encode(routers.Router{
				ID:          createdRouterID,
				Name:        *req.Name,
				Description: *req.Description,
				State:       true,
				Status:      "ACTIVE",
				ProjectID:   "proj-1",
			})

		// Step 6: Set state
		case r.Method == http.MethodPost && strings.Contains(r.URL.Path, "/action"):
			w.WriteHeader(http.StatusNoContent)

		// Step 7: Associate network
		case r.Method == http.MethodPost && strings.Contains(r.URL.Path, "/networks/net-"):
			w.WriteHeader(http.StatusNoContent)

		// Step 8: List router networks (more specific - ends with /networks, not /networks/something)
		case r.Method == http.MethodGet && strings.HasSuffix(r.URL.Path, "/networks"):
			w.WriteHeader(http.StatusOK)
			_ = json.NewEncoder(w).Encode([]routers.RouterNetwork{
				{
					NetworkID:   "net-lifecycle-123",
					NetworkName: "test-network",
				},
			})

		// Step 9: Disassociate network
		case r.Method == http.MethodDelete && strings.Contains(r.URL.Path, "/networks/net-"):
			w.WriteHeader(http.StatusNoContent)

		// Step 10: Delete router
		case r.Method == http.MethodDelete && strings.Contains(r.URL.Path, createdRouterID):
			w.WriteHeader(http.StatusNoContent)

		default:
			t.Logf("Unhandled request: %s %s", r.Method, r.URL.Path)
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer server.Close()

	client := cloudsdk.NewClient(server.URL, "test-token")
	vpsClient := client.Project("proj-1").VPS()
	routerClient := vpsClient.Routers()
	ctx := context.Background()

	// Step 1: List routers (should be empty)
	t.Log("Step 1: List routers (initial)")
	initialList, err := routerClient.List(ctx, nil)
	if err != nil {
		t.Fatalf("Step 1 failed: %v", err)
	}
	if len(initialList.Routers) != 0 {
		t.Errorf("Step 1: expected 0 routers, got %d", len(initialList.Routers))
	}

	// Step 2: Create router
	t.Log("Step 2: Create router")
	createReq := &routers.RouterCreateRequest{
		Name:        "lifecycle-router",
		Description: "Test lifecycle",
	}
	created, err := routerClient.Create(ctx, createReq)
	if err != nil {
		t.Fatalf("Step 2 failed: %v", err)
	}
	if created.ID == "" {
		t.Fatal("Step 2: created router has empty ID")
	}
	createdRouterID = created.ID

	// Step 3: Get router
	t.Log("Step 3: Get router")
	fetched, err := routerClient.Get(ctx, createdRouterID)
	if err != nil {
		t.Fatalf("Step 3 failed: %v", err)
	}
	if fetched.ID != createdRouterID {
		t.Errorf("Step 3: expected ID %s, got %s", createdRouterID, fetched.ID)
	}

	// Step 4: List routers (should contain created router)
	t.Log("Step 4: List routers (after creation)")
	updatedList, err := routerClient.List(ctx, nil)
	if err != nil {
		t.Fatalf("Step 4 failed: %v", err)
	}
	if len(updatedList.Routers) != 1 {
		t.Errorf("Step 4: expected 1 router, got %d", len(updatedList.Routers))
	}

	// Step 5: Update router
	t.Log("Step 5: Update router")
	newName := "updated-router"
	newDesc := "Updated description"
	updateReq := &routers.RouterUpdateRequest{
		Name:        &newName,
		Description: &newDesc,
	}
	updated, err := routerClient.Update(ctx, createdRouterID, updateReq)
	if err != nil {
		t.Fatalf("Step 5 failed: %v", err)
	}
	if updated.Name != newName {
		t.Errorf("Step 5: expected name %s, got %s", newName, updated.Name)
	}

	// Step 6: Set state (disable)
	t.Log("Step 6: Set router state")
	stateReq := &routers.RouterSetStateRequest{
		State: false,
	}
	err = routerClient.SetState(ctx, createdRouterID, stateReq)
	if err != nil {
		t.Fatalf("Step 6 failed: %v", err)
	}

	// Step 7: Associate network (using router resource pattern)
	t.Log("Step 7: Associate network")
	// Note: This assumes Get returns a RouterResource with Networks() method
	routerResource, err := routerClient.Get(ctx, createdRouterID)
	if err != nil {
		t.Fatalf("Step 7 (get router) failed: %v", err)
	}
	err = routerResource.Networks().Associate(ctx, "net-lifecycle-123")
	if err != nil {
		t.Fatalf("Step 7 failed: %v", err)
	}

	// Step 8: List router networks
	t.Log("Step 8: List router networks")
	networks, err := routerResource.Networks().List(ctx)
	if err != nil {
		t.Fatalf("Step 8 failed: %v", err)
	}
	if len(networks) != 1 {
		t.Errorf("Step 8: expected 1 network, got %d", len(networks))
	}

	// Step 9: Disassociate network
	t.Log("Step 9: Disassociate network")
	err = routerResource.Networks().Disassociate(ctx, "net-lifecycle-123")
	if err != nil {
		t.Fatalf("Step 9 failed: %v", err)
	}

	// Step 10: Delete router
	t.Log("Step 10: Delete router")
	err = routerClient.Delete(ctx, createdRouterID)
	if err != nil {
		t.Fatalf("Step 10 failed: %v", err)
	}

	t.Log("Integration test completed successfully")
}
