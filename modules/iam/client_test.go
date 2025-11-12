package iam_test

import (
	"testing"

	internalhttp "github.com/Zillaforge/cloud-sdk/internal/http"
	"github.com/Zillaforge/cloud-sdk/modules/iam"
	"github.com/Zillaforge/cloud-sdk/modules/iam/projects"
	"github.com/Zillaforge/cloud-sdk/modules/iam/users"
)

func TestNewClient(t *testing.T) {
	// Create a mock HTTP client
	baseClient := &internalhttp.Client{}

	// Create IAM client
	client := iam.NewClient(baseClient)

	if client == nil {
		t.Fatal("NewClient() returned nil")
	}
}

func TestClient_Users(t *testing.T) {
	// Create a mock HTTP client
	baseClient := &internalhttp.Client{}

	// Create IAM client
	client := iam.NewClient(baseClient)

	// Get users client
	usersClient := client.Users()

	if usersClient == nil {
		t.Fatal("Users() returned nil")
	}

	// Verify it's the correct type
	if _, ok := interface{}(usersClient).(*users.Client); !ok {
		t.Errorf("Users() returned wrong type, expected *users.Client")
	}
}

func TestClient_Projects(t *testing.T) {
	// Create a mock HTTP client
	baseClient := &internalhttp.Client{}

	// Create IAM client
	client := iam.NewClient(baseClient)

	// Get projects client
	projectsClient := client.Projects()

	if projectsClient == nil {
		t.Fatal("Projects() returned nil")
	}

	// Verify it's the correct type
	if _, ok := interface{}(projectsClient).(*projects.Client); !ok {
		t.Errorf("Projects() returned wrong type, expected *projects.Client")
	}
}

func TestClient_FactoryMethods_ReturnNewInstances(t *testing.T) {
	// Create a mock HTTP client
	baseClient := &internalhttp.Client{}

	// Create IAM client
	client := iam.NewClient(baseClient)

	// Call Users() twice and verify they're different instances
	users1 := client.Users()
	users2 := client.Users()

	// Compare pointers - they should be different instances
	if users1 == users2 {
		t.Error("Users() should return new instances, but returned same instance")
	}

	// Call Projects() twice and verify they're different instances
	projects1 := client.Projects()
	projects2 := client.Projects()

	// Compare pointers - they should be different instances
	if projects1 == projects2 {
		t.Error("Projects() should return new instances, but returned same instance")
	}
}
