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
			VCPU:        2,
			Memory:      4096,
			Disk:        20,
			Public:      true,
			Tags:        []string{"general", "starter"},
		},
		"flav-medium": {
			ID:          "flav-medium",
			Name:        "medium",
			Description: "Medium compute instance",
			VCPU:        4,
			Memory:      8192,
			Disk:        40,
			Public:      true,
			Tags:        []string{"general", "balanced"},
		},
		"flav-large": {
			ID:          "flav-large",
			Name:        "large",
			Description: "Large compute instance",
			VCPU:        8,
			Memory:      16384,
			Disk:        80,
			Public:      true,
			Tags:        []string{"compute", "performance"},
		},
		"flav-gpu": {
			ID:          "flav-gpu",
			Name:        "gpu-large",
			Description: "GPU-enabled large instance",
			VCPU:        16,
			Memory:      32768,
			Disk:        200,
			Public:      true,
			Tags:        []string{"gpu", "ml", "compute"},
		},
		"flav-private": {
			ID:          "flav-private",
			Name:        "custom-private",
			Description: "Private flavor for specific projects",
			VCPU:        4,
			Memory:      8192,
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

				// Filter by tags (supports multiple)
				if tags := query["tag"]; len(tags) > 0 {
					for _, queryTag := range tags {
						hasTag := false
						for _, flavorTag := range flavor.Tags {
							if flavorTag == queryTag {
								hasTag = true
								break
							}
						}
						if !hasTag {
							include = false
							break
						}
					}
				}

				if include {
					filteredFlavors = append(filteredFlavors, flavor)
				}
			}

			w.WriteHeader(http.StatusOK)
			_ = json.NewEncoder(w).Encode(&flavors.FlavorListResponse{Flavors: filteredFlavors})
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
		if len(resp.Flavors) != len(mockFlavors) {
			t.Errorf("expected %d flavors, got %d", len(mockFlavors), len(resp.Flavors))
		}
	})

	t.Run("filter by name", func(t *testing.T) {
		opts := &flavors.ListFlavorsOptions{Name: "small"}
		resp, err := flavorsClient.List(ctx, opts)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(resp.Flavors) != 1 {
			t.Errorf("expected 1 flavor, got %d", len(resp.Flavors))
		}
		if len(resp.Flavors) > 0 && resp.Flavors[0].Name != "small" {
			t.Errorf("expected flavor name 'small', got '%s'", resp.Flavors[0].Name)
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
		if len(resp.Flavors) != 4 {
			t.Errorf("expected 4 public flavors, got %d", len(resp.Flavors))
		}
		for _, flavor := range resp.Flavors {
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
		if len(resp.Flavors) != 1 {
			t.Errorf("expected 1 private flavor, got %d", len(resp.Flavors))
		}
		if len(resp.Flavors) > 0 && resp.Flavors[0].Public {
			t.Errorf("expected private flavor, got public: %s", resp.Flavors[0].ID)
		}
	})

	t.Run("filter by tag", func(t *testing.T) {
		opts := &flavors.ListFlavorsOptions{Tags: []string{"gpu"}}
		resp, err := flavorsClient.List(ctx, opts)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(resp.Flavors) != 1 {
			t.Errorf("expected 1 GPU flavor, got %d", len(resp.Flavors))
		}
		if len(resp.Flavors) > 0 {
			hasGPUTag := false
			for _, tag := range resp.Flavors[0].Tags {
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
		if flavor.VCPU != 8 {
			t.Errorf("expected 8 VCPU, got %d", flavor.VCPU)
		}
		if flavor.Memory != 16384 {
			t.Errorf("expected 16384 Memory, got %d", flavor.Memory)
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

		// Find the smallest flavor (by VCPU)
		var smallest *flavors.Flavor
		for _, flavor := range allFlavors.Flavors {
			if smallest == nil || flavor.VCPU < smallest.VCPU {
				smallest = flavor
			}
		}

		if smallest == nil {
			t.Fatal("expected to find smallest flavor")
		}
		if smallest.Name != "small" {
			t.Errorf("expected smallest flavor to be 'small', got '%s'", smallest.Name)
		}

		// Find the largest flavor (by VCPU)
		var largest *flavors.Flavor
		for _, flavor := range allFlavors.Flavors {
			if largest == nil || flavor.VCPU > largest.VCPU {
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
		if flavorDetails.VCPU != largest.VCPU {
			t.Errorf("expected VCPU %d, got %d", largest.VCPU, flavorDetails.VCPU)
		}
	})
}
