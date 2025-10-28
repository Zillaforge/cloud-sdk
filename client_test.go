package cloudsdk

import (
	"net/http"
	"testing"
	"time"
)

// mockLogger implements the Logger interface for testing
type mockLogger struct{}

func (m *mockLogger) Debug(_ string, _ ...interface{}) {}
func (m *mockLogger) Info(_ string, _ ...interface{})  {}
func (m *mockLogger) Warn(_ string, _ ...interface{})  {}
func (m *mockLogger) Error(_ string, _ ...interface{}) {}

// TestNew tests the New function with various scenarios
func TestNew(t *testing.T) {
	tests := []struct {
		name        string
		baseURL     string
		token       string
		expectError bool
		errorMsg    string
	}{
		{
			name:        "valid https URL and token",
			baseURL:     "https://api.example.com",
			token:       "test-token",
			expectError: false,
		},
		{
			name:        "valid http URL and token",
			baseURL:     "http://localhost:8080",
			token:       "test-token",
			expectError: false,
		},
		{
			name:        "URL without scheme",
			baseURL:     "api.example.com",
			token:       "test-token",
			expectError: true,
			errorMsg:    "base URL must include scheme",
		},
		{
			name:        "invalid URL",
			baseURL:     "ht!tp://invalid",
			token:       "test-token",
			expectError: true,
			errorMsg:    "invalid base URL",
		},
		{
			name:        "empty token",
			baseURL:     "https://api.example.com",
			token:       "",
			expectError: true,
			errorMsg:    "token cannot be empty",
		},
		{
			name:        "both invalid",
			baseURL:     "",
			token:       "",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := New(tt.baseURL, tt.token)

			if tt.expectError {
				if err == nil {
					t.Error("expected error, got nil")
				}
				if tt.errorMsg != "" && err != nil {
					if !containsString(err.Error(), tt.errorMsg) {
						t.Errorf("expected error containing %q, got %q", tt.errorMsg, err.Error())
					}
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if client == nil {
				t.Fatal("expected client, got nil")
			}

			if client.baseURL != tt.baseURL {
				t.Errorf("expected baseURL %q, got %q", tt.baseURL, client.baseURL)
			}

			if client.token != tt.token {
				t.Errorf("expected token %q, got %q", tt.token, client.token)
			}

			// Check default timeout
			if client.httpClient.Timeout != 30*time.Second {
				t.Errorf("expected default timeout 30s, got %v", client.httpClient.Timeout)
			}
		})
	}
}

// TestNewWithOptions tests the New function with various options
func TestNewWithOptions(t *testing.T) {
	t.Run("with logger option", func(t *testing.T) {
		logger := &mockLogger{}
		client, err := New("https://api.example.com", "test-token", WithLogger(logger))

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if client.logger != logger {
			t.Error("expected custom logger to be set")
		}
	})

	t.Run("with timeout option", func(t *testing.T) {
		customTimeout := 60 * time.Second
		client, err := New("https://api.example.com", "test-token", WithTimeout(customTimeout))

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if client.httpClient.Timeout != customTimeout {
			t.Errorf("expected timeout %v, got %v", customTimeout, client.httpClient.Timeout)
		}
	})

	t.Run("with custom HTTP client", func(t *testing.T) {
		customClient := &http.Client{
			Timeout: 45 * time.Second,
		}
		client, err := New("https://api.example.com", "test-token", WithHTTPClient(customClient))

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if client.httpClient != customClient {
			t.Error("expected custom HTTP client to be set")
		}

		if client.httpClient.Timeout != 45*time.Second {
			t.Errorf("expected timeout 45s, got %v", client.httpClient.Timeout)
		}
	})

	t.Run("with multiple options", func(t *testing.T) {
		logger := &mockLogger{}
		customTimeout := 90 * time.Second

		client, err := New("https://api.example.com", "test-token",
			WithLogger(logger),
			WithTimeout(customTimeout),
		)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if client.logger != logger {
			t.Error("expected custom logger to be set")
		}

		if client.httpClient.Timeout != customTimeout {
			t.Errorf("expected timeout %v, got %v", customTimeout, client.httpClient.Timeout)
		}
	})
}

// TestNewClient tests the convenience function
func TestNewClient(t *testing.T) {
	t.Run("valid parameters", func(t *testing.T) {
		client := NewClient("https://api.example.com", "test-token")

		if client == nil {
			t.Fatal("expected client, got nil")
		}

		if client.baseURL != "https://api.example.com" {
			t.Errorf("expected baseURL %q, got %q", "https://api.example.com", client.baseURL)
		}

		if client.token != "test-token" {
			t.Errorf("expected token %q, got %q", "test-token", client.token)
		}
	})

	t.Run("invalid parameters return nil", func(t *testing.T) {
		client := NewClient("", "")

		if client != nil {
			t.Error("expected nil client for invalid parameters")
		}
	})
}

// TestClient_Accessors tests the accessor methods
func TestClient_Accessors(t *testing.T) {
	client, err := New("https://api.example.com", "test-token")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	t.Run("BaseURL", func(t *testing.T) {
		if client.BaseURL() != "https://api.example.com" {
			t.Errorf("expected baseURL %q, got %q", "https://api.example.com", client.BaseURL())
		}
	})

	t.Run("HTTPClient", func(t *testing.T) {
		httpClient := client.HTTPClient()
		if httpClient == nil {
			t.Error("expected HTTP client, got nil")
			return
		}
		if httpClient.Timeout != 30*time.Second {
			t.Errorf("expected timeout 30s, got %v", httpClient.Timeout)
		}
	})
}

// TestClient_Project tests the Project method
func TestClient_Project(t *testing.T) {
	client, err := New("https://api.example.com", "test-token")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	projectClient := client.Project("proj-123")

	if projectClient == nil {
		t.Fatal("expected project client, got nil")
	}

	if projectClient.projectID != "proj-123" {
		t.Errorf("expected projectID %q, got %q", "proj-123", projectClient.projectID)
	}

	if projectClient.client != client {
		t.Error("expected project client to reference parent client")
	}
}

// TestProjectClient_VPS tests the VPS method
func TestProjectClient_VPS(t *testing.T) {
	client, err := New("https://api.example.com", "test-token")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	projectClient := client.Project("proj-123")
	vpsClient := projectClient.VPS()

	if vpsClient == nil {
		t.Fatal("expected VPS client, got nil")
	}

	// Test that VPS client is properly initialized
	// We can verify this by checking the Networks accessor
	networksClient := vpsClient.Networks()
	if networksClient == nil {
		t.Error("expected networks client, got nil")
	}
}

// TestClient_Integration tests the full client creation and service access flow
func TestClient_Integration(t *testing.T) {
	// Create client
	client, err := New("https://api.example.com", "test-token")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	// Access project-scoped VPS service
	vpsClient := client.Project("proj-123").VPS()
	if vpsClient == nil {
		t.Fatal("expected VPS client, got nil")
	}

	// Verify VPS client can access sub-services
	networksClient := vpsClient.Networks()
	if networksClient == nil {
		t.Error("expected networks client from VPS client")
	}
}

// TestClient_MultipleProjects tests accessing multiple projects
func TestClient_MultipleProjects(t *testing.T) {
	client, err := New("https://api.example.com", "test-token")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	// Create project clients for different projects
	proj1 := client.Project("proj-1")
	proj2 := client.Project("proj-2")

	if proj1.projectID == proj2.projectID {
		t.Error("expected different project IDs")
	}

	// Both should reference the same parent client
	if proj1.client != proj2.client {
		t.Error("expected both project clients to reference the same parent client")
	}

	// Each should create independent VPS clients
	vps1 := proj1.VPS()
	vps2 := proj2.VPS()

	if vps1 == vps2 {
		t.Error("expected different VPS client instances for different projects")
	}
}

// containsString checks if a string contains a substring
func containsString(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && containsSubstring(s, substr))
}

func containsSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
