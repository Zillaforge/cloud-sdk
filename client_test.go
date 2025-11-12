package cloudsdk

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	iamprojects "github.com/Zillaforge/cloud-sdk/models/iam/projects"
)

// mockLogger implements the Logger interface for testing
type mockLogger struct{}

func (m *mockLogger) Debug(_ string, _ ...interface{}) {}
func (m *mockLogger) Info(_ string, _ ...interface{})  {}
func (m *mockLogger) Warn(_ string, _ ...interface{})  {}
func (m *mockLogger) Error(_ string, _ ...interface{}) {}

// TestClient_IAM tests the IAM() method
func TestClient_IAM(t *testing.T) {
	client, err := New("https://api.example.com", "test-token")
	if err != nil {
		t.Fatalf("New() failed: %v", err)
	}

	iamClient := client.IAM()
	if iamClient == nil {
		t.Fatal("IAM() returned nil")
	}
}

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
	// Create mock IAM server
	iamServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if r.Method == http.MethodGet && r.URL.Path == "/iam/api/v1/project/proj-123" {
			w.WriteHeader(http.StatusOK)
			response := map[string]interface{}{
				"projectId":   "proj-123",
				"displayName": "Test Project",
			}
			_ = json.NewEncoder(w).Encode(response)
			return
		}
		w.WriteHeader(http.StatusNotFound)
	}))
	defer iamServer.Close()

	client, err := New(iamServer.URL, "test-token")
	if err != nil {
		t.Fatalf("New() failed: %v", err)
	}

	projectClient, err := client.Project(context.Background(), "proj-123")
	if err != nil {
		t.Fatalf("Project() failed: %v", err)
	}

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

// mockIAMServerForProjectTests creates a mock HTTP server for IAM API calls in Project tests
func mockIAMServerForProjectTests(t *testing.T, projects []iamprojects.ProjectMembership) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Logf("Mock server received request: %s %s", r.Method, r.URL.String())
		w.Header().Set("Content-Type", "application/json")

		// Handle GET /iam/api/v1/project/{projectID}
		if r.Method == http.MethodGet && r.URL.Path == "/iam/api/v1/project/proj-123" {
			t.Logf("Handling project get request for proj-123")
			w.WriteHeader(http.StatusOK)
			project := map[string]interface{}{
				"projectId":   "proj-123",
				"displayName": "Test Project",
				"extra": map[string]interface{}{
					"iservice": map[string]interface{}{
						"projectSysCode": "TEST-PROJ",
					},
				},
			}
			_ = json.NewEncoder(w).Encode(project)
			return
		}

		// Handle GET /iam/api/v1/project/{projectID} for non-existent project
		if r.Method == http.MethodGet && r.URL.Path == "/iam/api/v1/project/nonexistent" {
			t.Logf("Handling nonexistent project request")
			w.WriteHeader(http.StatusNotFound)
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"error": "Project not found",
			})
			return
		}

		// Handle GET /iam/api/v1/projects (list all projects)
		if r.Method == http.MethodGet && r.URL.Path == "/iam/api/v1/projects" {
			t.Logf("Handling projects list request")
			w.WriteHeader(http.StatusOK)
			response := map[string]interface{}{
				"projects": projects,
			}
			_ = json.NewEncoder(w).Encode(response)
			return
		}

		// Default: not found
		t.Logf("Unhandled request: %s %s", r.Method, r.URL.Path)
		w.WriteHeader(http.StatusNotFound)
	}))
}

// TestProject_WithValidProjectID tests Project() with a valid project ID
func TestProject_WithValidProjectID(t *testing.T) {
	// Setup mock IAM server
	iamServer := mockIAMServerForProjectTests(t, []iamprojects.ProjectMembership{})
	defer iamServer.Close()

	// Create client with mock IAM server URL
	client, err := New(iamServer.URL, "test-token")
	if err != nil {
		t.Fatalf("New() failed: %v", err)
	}

	// Test with valid project ID
	projectClient, err := client.Project(context.Background(), "proj-123")
	if err != nil {
		t.Fatalf("Project() failed: %v", err)
	}

	if projectClient == nil {
		t.Fatal("expected project client, got nil")
	}

	if projectClient.projectID != "proj-123" {
		t.Errorf("expected projectID 'proj-123', got '%s'", projectClient.projectID)
	}

	if projectClient.client != client {
		t.Error("expected project client to reference parent client")
	}
}

// TestProject_WithValidProjectCode tests Project() with a valid project code
func TestProject_WithValidProjectCode(t *testing.T) {
	// Setup mock projects with projectSysCode
	mockProjects := []iamprojects.ProjectMembership{
		{
			Project: &iamprojects.Project{
				ProjectID:   "proj-456",
				DisplayName: "Test Project 456",
				Extra: map[string]interface{}{
					"iservice": map[string]interface{}{
						"projectSysCode": "TEST-CODE",
					},
				},
			},
		},
	}

	iamServer := mockIAMServerForProjectTests(t, mockProjects)
	defer iamServer.Close()

	client, err := New(iamServer.URL, "test-token")
	if err != nil {
		t.Fatalf("New() failed: %v", err)
	}

	// Test with project code that should resolve to proj-456
	projectClient, err := client.Project(context.Background(), "TEST-CODE")
	if err != nil {
		t.Fatalf("Project() failed: %v", err)
	}

	if projectClient == nil {
		t.Fatal("expected project client, got nil")
	}

	if projectClient.projectID != "proj-456" {
		t.Errorf("expected projectID 'proj-456', got '%s'", projectClient.projectID)
	}
}

// TestProject_WithInvalidProjectID tests Project() with invalid project ID
func TestProject_WithInvalidProjectID(t *testing.T) {
	iamServer := mockIAMServerForProjectTests(t, []iamprojects.ProjectMembership{})
	defer iamServer.Close()

	client, newErr := New(iamServer.URL, "test-token")
	if newErr != nil {
		t.Fatalf("New() failed: %v", newErr)
	}

	// Test with invalid project ID
	_, err := client.Project(context.Background(), "nonexistent")
	if err == nil {
		t.Fatal("expected error for nonexistent project ID, got nil")
	}

	expectedError := "no project found with projectSysCode nonexistent"
	if !strings.Contains(err.Error(), expectedError) {
		t.Errorf("expected error containing '%s', got '%s'", expectedError, err.Error())
	}
}

// TestProject_WithInvalidProjectCode tests Project() with invalid project code
func TestProject_WithInvalidProjectCode(t *testing.T) {
	// Empty projects list - no project codes to match
	iamServer := mockIAMServerForProjectTests(t, []iamprojects.ProjectMembership{})
	defer iamServer.Close()

	client, newErr := New(iamServer.URL, "test-token")
	if newErr != nil {
		t.Fatalf("New() failed: %v", newErr)
	}

	// Test with invalid project code
	_, err := client.Project(context.Background(), "INVALID-CODE")
	if err == nil {
		t.Fatal("expected error for invalid project code, got nil")
	}

	expectedError := "no project found with projectSysCode INVALID-CODE"
	if !strings.Contains(err.Error(), expectedError) {
		t.Errorf("expected error containing '%s', got '%s'", expectedError, err.Error())
	}
}

// TestProject_WithMultipleMatchingCodes tests Project() when multiple projects have the same projectSysCode
func TestProject_WithMultipleMatchingCodes(t *testing.T) {
	// Setup mock projects with duplicate projectSysCode
	mockProjects := []iamprojects.ProjectMembership{
		{
			Project: &iamprojects.Project{
				ProjectID:   "proj-1",
				DisplayName: "Test Project 1",
				Extra: map[string]interface{}{
					"iservice": map[string]interface{}{
						"projectSysCode": "DUPLICATE-CODE",
					},
				},
			},
		},
		{
			Project: &iamprojects.Project{
				ProjectID:   "proj-2",
				DisplayName: "Test Project 2",
				Extra: map[string]interface{}{
					"iservice": map[string]interface{}{
						"projectSysCode": "DUPLICATE-CODE",
					},
				},
			},
		},
	}

	iamServer := mockIAMServerForProjectTests(t, mockProjects)
	defer iamServer.Close()

	client, newErr := New(iamServer.URL, "test-token")
	if newErr != nil {
		t.Fatalf("New() failed: %v", newErr)
	}

	// Test with duplicate project code
	_, err := client.Project(context.Background(), "DUPLICATE-CODE")
	if err == nil {
		t.Fatal("expected error for duplicate project code, got nil")
	}

	expectedError := "multiple projects found with projectSysCode DUPLICATE-CODE"
	if !strings.Contains(err.Error(), expectedError) {
		t.Errorf("expected error containing '%s', got '%s'", expectedError, err.Error())
	}
}

// TestProject_IAMServerError tests Project() when IAM server returns an error
func TestProject_IAMServerError(t *testing.T) {
	// Create a server that always returns 500
	iamServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"error": "Internal server error",
		})
	}))
	defer iamServer.Close()

	client, newErr := New(iamServer.URL, "test-token")
	if newErr != nil {
		t.Fatalf("New() failed: %v", newErr)
	}

	// Test with any project ID - should fail due to server error
	_, err := client.Project(context.Background(), "any-project")
	if err == nil {
		t.Fatal("expected error due to IAM server error, got nil")
	}
}

// TestProjectClient_ServiceClients tests the VPS and VRM methods
func TestProjectClient_ServiceClients(t *testing.T) {
	// Create mock IAM server
	iamServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if r.Method == http.MethodGet && r.URL.Path == "/iam/api/v1/project/proj-123" {
			w.WriteHeader(http.StatusOK)
			response := map[string]interface{}{
				"projectId":   "proj-123",
				"displayName": "Test Project",
			}
			_ = json.NewEncoder(w).Encode(response)
			return
		}
		w.WriteHeader(http.StatusNotFound)
	}))
	defer iamServer.Close()

	client, err := New(iamServer.URL, "test-token")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	projectClient, err := client.Project(context.Background(), "proj-123")
	if err != nil {
		t.Fatalf("Project() failed: %v", err)
	}
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

	// Also test VRM method
	vrmClient := projectClient.VRM()

	if vrmClient == nil {
		t.Fatal("expected VRM client, got nil")
	}

	// Test that VRM client is properly initialized
	// We can verify this by checking the Repositories accessor
	repositoriesClient := vrmClient.Repositories()
	if repositoriesClient == nil {
		t.Error("expected repositories client, got nil")
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
