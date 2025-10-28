package flavors_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	cloudsdk "github.com/Zillaforge/cloud-sdk"
	"github.com/Zillaforge/cloud-sdk/models/vps/flavors"
)

// TestFlavorDiscovery verifies the complete flavor discovery workflow
func TestFlavorDiscovery(t *testing.T) {
	// Mock flavors database
	mockFlavors := map[string]*flavors.Flavor{
		"flav-small": {
			ID:          "flav-small",
			Name:        "small",
			Description: "Small compute instance",
			VCPUs:       2,
			RAM:         4096,
			Disk:        20,
			Public:      true,
			Tags:        []string{"general", "starter"},
		},
		"flav-medium": {
			ID:          "flav-medium",
			Name:        "medium",
			Description: "Medium compute instance",
			VCPUs:       4,
			RAM:         8192,
			Disk:        40,
			Public:      true,
			Tags:        []string{"general", "balanced"},
		},
		"flav-large": {
			ID:          "flav-large",
			Name:        "large",
			Description: "Large compute instance",
			VCPUs:       8,
			RAM:         16384,
			Disk:        80,
			Public:      true,
			Tags:        []string{"compute", "performance"},
		},
		"flav-gpu": {
			ID:          "flav-gpu",
			Name:        "gpu-large",
			Description: "GPU-enabled large instance",
			VCPUs:       16,
			RAM:         32768,
			Disk:        200,
			Public:      true,
			Tags:        []string{"gpu", "ml", "compute"},
		},
		"flav-private": {
			ID:          "flav-private",
			Name:        "custom-private",
			Description: "Private flavor for specific projects",
			VCPUs:       4,
			RAM:         8192,
			Disk:        40,
			Public:      false,
			Tags:        []string{"custom"},
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": "Method not allowed"})
			return
		}

		// Handle list flavors
		if r.URL.Path == "/vps/api/v1/project/proj-123/flavors" {
			query := r.URL.Query()
			var filteredFlavors []*flavors.Flavor

			// Apply filters
			for _, flavor := range mockFlavors {
				include := true

				// Filter by name
				if name := query.Get("name"); name != "" && flavor.Name != name {
					include = false
				}

				// Filter by public flag
				if publicStr := query.Get("public"); publicStr != "" {
					if (publicStr == "true" && !flavor.Public) || (publicStr == "false" && flavor.Public) {
						include = false
					}
				}

				// Filter by tag
				if tag := query.Get("tag"); tag != "" {
					hasTag := false
					for _, t := range flavor.Tags {
						if t == tag {
							hasTag = true
							break
						}
					}
					if !hasTag {
						include = false
					}
				}

				if include {
					filteredFlavors = append(filteredFlavors, flavor)
				}
			}

			w.WriteHeader(http.StatusOK)
			_ = json.NewEncoder(w).Encode(&flavors.FlavorListResponse{Items: filteredFlavors})
			return
		}

		// Handle get specific flavor
		if len(r.URL.Path) > len("/vps/api/v1/project/proj-123/flavors/") {
			flavorID := r.URL.Path[len("/vps/api/v1/project/proj-123/flavors/"):]
			if flavor, ok := mockFlavors[flavorID]; ok {
				w.WriteHeader(http.StatusOK)
				_ = json.NewEncoder(w).Encode(flavor)
				return
			}
			w.WriteHeader(http.StatusNotFound)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": "Flavor not found"})
			return
		}

		w.WriteHeader(http.StatusNotFound)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "Not found"})
	}))
	defer server.Close()

	client := cloudsdk.NewClient(server.URL, "test-token")
	vpsClient := client.Project("proj-123").VPS()
	flavorsClient := vpsClient.Flavors()
	ctx := context.Background()

	t.Run("list all flavors", func(t *testing.T) {
		resp, err := flavorsClient.List(ctx, nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(resp.Items) != len(mockFlavors) {
			t.Errorf("expected %d flavors, got %d", len(mockFlavors), len(resp.Items))
		}
	})

	t.Run("filter by name", func(t *testing.T) {
		opts := &flavors.ListFlavorsOptions{Name: "small"}
		resp, err := flavorsClient.List(ctx, opts)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(resp.Items) != 1 {
			t.Errorf("expected 1 flavor, got %d", len(resp.Items))
		}
		if len(resp.Items) > 0 && resp.Items[0].Name != "small" {
			t.Errorf("expected flavor name 'small', got '%s'", resp.Items[0].Name)
		}
	})

	t.Run("filter by public=true", func(t *testing.T) {
		publicFlag := true
		opts := &flavors.ListFlavorsOptions{Public: &publicFlag}
		resp, err := flavorsClient.List(ctx, opts)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		// Should return all public flavors (4 out of 5)
		if len(resp.Items) != 4 {
			t.Errorf("expected 4 public flavors, got %d", len(resp.Items))
		}
		for _, flavor := range resp.Items {
			if !flavor.Public {
				t.Errorf("expected public flavor, got private: %s", flavor.ID)
			}
		}
	})

	t.Run("filter by public=false", func(t *testing.T) {
		publicFlag := false
		opts := &flavors.ListFlavorsOptions{Public: &publicFlag}
		resp, err := flavorsClient.List(ctx, opts)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		// Should return only private flavors (1 out of 5)
		if len(resp.Items) != 1 {
			t.Errorf("expected 1 private flavor, got %d", len(resp.Items))
		}
		if len(resp.Items) > 0 && resp.Items[0].Public {
			t.Errorf("expected private flavor, got public: %s", resp.Items[0].ID)
		}
	})

	t.Run("filter by tag", func(t *testing.T) {
		opts := &flavors.ListFlavorsOptions{Tag: "gpu"}
		resp, err := flavorsClient.List(ctx, opts)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(resp.Items) != 1 {
			t.Errorf("expected 1 GPU flavor, got %d", len(resp.Items))
		}
		if len(resp.Items) > 0 {
			hasGPUTag := false
			for _, tag := range resp.Items[0].Tags {
				if tag == "gpu" {
					hasGPUTag = true
					break
				}
			}
			if !hasGPUTag {
				t.Error("expected flavor with 'gpu' tag")
			}
		}
	})

	t.Run("get specific flavor", func(t *testing.T) {
		flavor, err := flavorsClient.Get(ctx, "flav-large")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if flavor == nil {
			t.Fatal("expected flavor, got nil")
		}
		if flavor.ID != "flav-large" {
			t.Errorf("expected flavor ID 'flav-large', got '%s'", flavor.ID)
		}
		if flavor.VCPUs != 8 {
			t.Errorf("expected 8 VCPUs, got %d", flavor.VCPUs)
		}
		if flavor.RAM != 16384 {
			t.Errorf("expected 16384 RAM, got %d", flavor.RAM)
		}
	})

	t.Run("get nonexistent flavor", func(t *testing.T) {
		_, err := flavorsClient.Get(ctx, "nonexistent")
		if err == nil {
			t.Error("expected error for nonexistent flavor, got nil")
		}
	})

	t.Run("discover flavors by characteristics", func(t *testing.T) {
		// List all flavors
		allFlavors, err := flavorsClient.List(ctx, nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		// Find the smallest flavor (by VCPUs)
		var smallest *flavors.Flavor
		for _, flavor := range allFlavors.Items {
			if smallest == nil || flavor.VCPUs < smallest.VCPUs {
				smallest = flavor
			}
		}

		if smallest == nil {
			t.Fatal("expected to find smallest flavor")
		}
		if smallest.Name != "small" {
			t.Errorf("expected smallest flavor to be 'small', got '%s'", smallest.Name)
		}

		// Find the largest flavor (by VCPUs)
		var largest *flavors.Flavor
		for _, flavor := range allFlavors.Items {
			if largest == nil || flavor.VCPUs > largest.VCPUs {
				largest = flavor
			}
		}

		if largest == nil {
			t.Fatal("expected to find largest flavor")
		}
		if largest.Name != "gpu-large" {
			t.Errorf("expected largest flavor to be 'gpu-large', got '%s'", largest.Name)
		}

		// Verify we can retrieve details of the largest flavor
		flavorDetails, err := flavorsClient.Get(ctx, largest.ID)
		if err != nil {
			t.Fatalf("unexpected error getting flavor details: %v", err)
		}
		if flavorDetails.VCPUs != largest.VCPUs {
			t.Errorf("expected VCPUs %d, got %d", largest.VCPUs, flavorDetails.VCPUs)
		}
	})
}
