package users_test

import (
	"encoding/json"
	"testing"

	"github.com/Zillaforge/cloud-sdk/models/iam/users"
)

func TestUser_JSONParsing_ValidResponse(t *testing.T) {
	jsonData := `{
		"userId": "4990ccdb-a9b1-49e5-91df-67c921601d81",
		"account": "system@ci.asus.com",
		"displayName": "system",
		"description": "This is system account",
		"extra": {},
		"namespace": "ci.asus.com",
		"email": "system@ci.asus.com",
		"frozen": false,
		"mfa": false,
		"createdAt": "2025-11-11T15:18:36Z",
		"updatedAt": "2025-11-11T15:32:32Z",
		"lastLoginAt": "2025-11-11T15:32:32Z"
	}`

	var user users.User
	err := json.Unmarshal([]byte(jsonData), &user)
	if err != nil {
		t.Fatalf("Failed to unmarshal user: %v", err)
	}

	// Verify all fields
	if user.UserID != "4990ccdb-a9b1-49e5-91df-67c921601d81" {
		t.Errorf("UserID = %v, want %v", user.UserID, "4990ccdb-a9b1-49e5-91df-67c921601d81")
	}
	if user.Account != "system@ci.asus.com" {
		t.Errorf("Account = %v, want %v", user.Account, "system@ci.asus.com")
	}
	if user.DisplayName != "system" {
		t.Errorf("DisplayName = %v, want %v", user.DisplayName, "system")
	}
	if user.Description != "This is system account" {
		t.Errorf("Description = %v, want %v", user.Description, "This is system account")
	}
	if user.Namespace != "ci.asus.com" {
		t.Errorf("Namespace = %v, want %v", user.Namespace, "ci.asus.com")
	}
	if user.Email != "system@ci.asus.com" {
		t.Errorf("Email = %v, want %v", user.Email, "system@ci.asus.com")
	}
	if user.Frozen != false {
		t.Errorf("Frozen = %v, want %v", user.Frozen, false)
	}
	if user.MFA != false {
		t.Errorf("MFA = %v, want %v", user.MFA, false)
	}
	if user.CreatedAt != "2025-11-11T15:18:36Z" {
		t.Errorf("CreatedAt = %v, want %v", user.CreatedAt, "2025-11-11T15:18:36Z")
	}
	if user.UpdatedAt != "2025-11-11T15:32:32Z" {
		t.Errorf("UpdatedAt = %v, want %v", user.UpdatedAt, "2025-11-11T15:32:32Z")
	}
	if user.LastLoginAt != "2025-11-11T15:32:32Z" {
		t.Errorf("LastLoginAt = %v, want %v", user.LastLoginAt, "2025-11-11T15:32:32Z")
	}
}

func TestUser_JSONParsing_WithUnknownFields(t *testing.T) {
	// Test forward compatibility - unknown fields should be ignored
	jsonData := `{
		"userId": "test-id",
		"account": "test@example.com",
		"displayName": "Test User",
		"description": "",
		"extra": {},
		"namespace": "test.com",
		"email": "test@example.com",
		"frozen": false,
		"mfa": true,
		"createdAt": "2025-01-01T00:00:00Z",
		"updatedAt": "2025-01-01T00:00:00Z",
		"lastLoginAt": "2025-01-01T00:00:00Z",
		"unknownField": "should be ignored",
		"anotherUnknownField": 12345
	}`

	var user users.User
	err := json.Unmarshal([]byte(jsonData), &user)
	if err != nil {
		t.Fatalf("Failed to unmarshal user with unknown fields: %v", err)
	}

	// Verify known fields are still parsed correctly
	if user.UserID != "test-id" {
		t.Errorf("UserID = %v, want %v", user.UserID, "test-id")
	}
	if user.Account != "test@example.com" {
		t.Errorf("Account = %v, want %v", user.Account, "test@example.com")
	}
	if user.MFA != true {
		t.Errorf("MFA = %v, want %v", user.MFA, true)
	}
}

func TestUser_JSONParsing_WithNestedExtra(t *testing.T) {
	// Test with complex nested extra metadata
	jsonData := `{
		"userId": "test-id",
		"account": "test@example.com",
		"displayName": "Test User",
		"description": "",
		"extra": {
			"department": "Engineering",
			"metadata": {
				"level": 5,
				"tags": ["senior", "backend"]
			}
		},
		"namespace": "test.com",
		"email": "test@example.com",
		"frozen": false,
		"mfa": false,
		"createdAt": "2025-01-01T00:00:00Z",
		"updatedAt": "2025-01-01T00:00:00Z",
		"lastLoginAt": "2025-01-01T00:00:00Z"
	}`

	var user users.User
	err := json.Unmarshal([]byte(jsonData), &user)
	if err != nil {
		t.Fatalf("Failed to unmarshal user with nested extra: %v", err)
	}

	// Verify extra is parsed as map
	if user.Extra == nil {
		t.Fatal("Extra should not be nil")
	}
	if len(user.Extra) == 0 {
		t.Error("Extra should contain data")
	}
}

func TestGetUserResponse_JSONParsing(t *testing.T) {
	jsonData := `{
		"userId": "test-id",
		"account": "test@example.com",
		"displayName": "Test User",
		"description": "Test description",
		"extra": {},
		"namespace": "test.com",
		"email": "test@example.com",
		"frozen": false,
		"mfa": true,
		"createdAt": "2025-01-01T00:00:00Z",
		"updatedAt": "2025-01-01T00:00:00Z",
		"lastLoginAt": "2025-01-01T00:00:00Z"
	}`

	var response users.GetUserResponse
	err := json.Unmarshal([]byte(jsonData), &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal GetUserResponse: %v", err)
	}

	// Verify embedded User fields are accessible
	if response.UserID != "test-id" {
		t.Errorf("UserID = %v, want %v", response.UserID, "test-id")
	}
	if response.Account != "test@example.com" {
		t.Errorf("Account = %v, want %v", response.Account, "test@example.com")
	}
	if response.DisplayName != "Test User" {
		t.Errorf("DisplayName = %v, want %v", response.DisplayName, "Test User")
	}
}
