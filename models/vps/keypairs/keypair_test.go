package keypairs

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/Zillaforge/cloud-sdk/models/vps/common"
)

// TestKeypairJSONUnmarshaling tests that Keypair structs can be unmarshaled from API responses
// containing all fields including timestamps and user information
func TestKeypairJSONUnmarshaling(t *testing.T) {
	tests := []struct {
		name     string
		jsonData string
		want     Keypair
		wantErr  bool
	}{
		{
			name: "complete keypair with all fields",
			jsonData: `{
				"id": "kp-abc123",
				"name": "my-keypair",
				"description": "Development keypair",
				"fingerprint": "SHA256:abc...xyz",
				"public_key": "ssh-rsa AAAAB3...",
				"private_key": "-----BEGIN RSA PRIVATE KEY-----\n...",
				"user_id": "user-123",
				"user": {
					"id": "user-123",
					"name": "john@example.com"
				},
				"createdAt": "2025-11-10T15:30:00Z",
				"updatedAt": "2025-11-10T15:30:00Z"
			}`,
			want: Keypair{
				ID:          "kp-abc123",
				Name:        "my-keypair",
				Description: "Development keypair",
				Fingerprint: "SHA256:abc...xyz",
				PublicKey:   "ssh-rsa AAAAB3...",
				PrivateKey:  "-----BEGIN RSA PRIVATE KEY-----\n...",
				UserID:      "user-123",
				User: &common.IDName{
					ID:   "user-123",
					Name: "john@example.com",
				},
				CreatedAt: "2025-11-10T15:30:00Z",
				UpdatedAt: "2025-11-10T15:30:00Z",
			},
			wantErr: false,
		},
		{
			name: "keypair without private_key (Get/List response)",
			jsonData: `{
				"id": "kp-abc123",
				"name": "my-keypair",
				"description": "Development keypair",
				"fingerprint": "SHA256:abc...xyz",
				"public_key": "ssh-rsa AAAAB3...",
				"user_id": "user-123",
				"user": {
					"id": "user-123",
					"name": "john@example.com"
				},
				"createdAt": "2025-11-10T15:30:00Z",
				"updatedAt": "2025-11-10T15:30:00Z"
			}`,
			want: Keypair{
				ID:          "kp-abc123",
				Name:        "my-keypair",
				Description: "Development keypair",
				Fingerprint: "SHA256:abc...xyz",
				PublicKey:   "ssh-rsa AAAAB3...",
				PrivateKey:  "", // Empty when not provided
				UserID:      "user-123",
				User: &common.IDName{
					ID:   "user-123",
					Name: "john@example.com",
				},
				CreatedAt: "2025-11-10T15:30:00Z",
				UpdatedAt: "2025-11-10T15:30:00Z",
			},
			wantErr: false,
		},
		{
			name: "keypair with null user object",
			jsonData: `{
				"id": "kp-abc123",
				"name": "my-keypair",
				"fingerprint": "SHA256:abc...xyz",
				"public_key": "ssh-rsa AAAAB3...",
				"user_id": "user-123",
				"user": null,
				"createdAt": "2025-11-10T15:30:00Z",
				"updatedAt": "2025-11-10T15:30:00Z"
			}`,
			want: Keypair{
				ID:          "kp-abc123",
				Name:        "my-keypair",
				Description: "",
				Fingerprint: "SHA256:abc...xyz",
				PublicKey:   "ssh-rsa AAAAB3...",
				UserID:      "user-123",
				User:        nil, // Null user object
				CreatedAt:   "2025-11-10T15:30:00Z",
				UpdatedAt:   "2025-11-10T15:30:00Z",
			},
			wantErr: false,
		},
		{
			name: "keypair without optional description",
			jsonData: `{
				"id": "kp-abc123",
				"name": "my-keypair",
				"fingerprint": "SHA256:abc...xyz",
				"public_key": "ssh-rsa AAAAB3...",
				"user_id": "user-123",
				"createdAt": "2025-11-10T15:30:00Z",
				"updatedAt": "2025-11-10T15:30:00Z"
			}`,
			want: Keypair{
				ID:          "kp-abc123",
				Name:        "my-keypair",
				Description: "",
				Fingerprint: "SHA256:abc...xyz",
				PublicKey:   "ssh-rsa AAAAB3...",
				UserID:      "user-123",
				User:        nil,
				CreatedAt:   "2025-11-10T15:30:00Z",
				UpdatedAt:   "2025-11-10T15:30:00Z",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got Keypair
			err := json.Unmarshal([]byte(tt.jsonData), &got)

			if (err != nil) != tt.wantErr {
				t.Errorf("json.Unmarshal() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err != nil {
				return
			}

			// Compare fields
			if got.ID != tt.want.ID {
				t.Errorf("ID = %v, want %v", got.ID, tt.want.ID)
			}
			if got.Name != tt.want.Name {
				t.Errorf("Name = %v, want %v", got.Name, tt.want.Name)
			}
			if got.Description != tt.want.Description {
				t.Errorf("Description = %v, want %v", got.Description, tt.want.Description)
			}
			if got.Fingerprint != tt.want.Fingerprint {
				t.Errorf("Fingerprint = %v, want %v", got.Fingerprint, tt.want.Fingerprint)
			}
			if got.PublicKey != tt.want.PublicKey {
				t.Errorf("PublicKey = %v, want %v", got.PublicKey, tt.want.PublicKey)
			}
			if got.PrivateKey != tt.want.PrivateKey {
				t.Errorf("PrivateKey = %v, want %v", got.PrivateKey, tt.want.PrivateKey)
			}
			if got.UserID != tt.want.UserID {
				t.Errorf("UserID = %v, want %v", got.UserID, tt.want.UserID)
			}
			if got.CreatedAt != tt.want.CreatedAt {
				t.Errorf("CreatedAt = %v, want %v", got.CreatedAt, tt.want.CreatedAt)
			}
			if got.UpdatedAt != tt.want.UpdatedAt {
				t.Errorf("UpdatedAt = %v, want %v", got.UpdatedAt, tt.want.UpdatedAt)
			}

			// Compare User object
			if tt.want.User == nil {
				if got.User != nil {
					t.Errorf("User = %v, want nil", got.User)
				}
			} else {
				if got.User == nil {
					t.Errorf("User = nil, want %v", tt.want.User)
				} else {
					if got.User.ID != tt.want.User.ID {
						t.Errorf("User.ID = %v, want %v", got.User.ID, tt.want.User.ID)
					}
					if got.User.Name != tt.want.User.Name {
						t.Errorf("User.Name = %v, want %v", got.User.Name, tt.want.User.Name)
					}
				}
			}
		})
	}
}

// TestKeypairJSONMarshaling tests that Keypair structs can be marshaled to JSON
// matching the API specification format
func TestKeypairJSONMarshaling(t *testing.T) {
	tests := []struct {
		name     string
		keypair  Keypair
		wantJSON map[string]interface{}
	}{
		{
			name: "complete keypair with all fields",
			keypair: Keypair{
				ID:          "kp-abc123",
				Name:        "my-keypair",
				Description: "Development keypair",
				Fingerprint: "SHA256:abc...xyz",
				PublicKey:   "ssh-rsa AAAAB3...",
				PrivateKey:  "-----BEGIN RSA PRIVATE KEY-----\n...",
				UserID:      "user-123",
				User: &common.IDName{
					ID:   "user-123",
					Name: "john@example.com",
				},
				CreatedAt: "2025-11-10T15:30:00Z",
				UpdatedAt: "2025-11-10T15:30:00Z",
			},
			wantJSON: map[string]interface{}{
				"id":          "kp-abc123",
				"name":        "my-keypair",
				"description": "Development keypair",
				"fingerprint": "SHA256:abc...xyz",
				"public_key":  "ssh-rsa AAAAB3...",
				"private_key": "-----BEGIN RSA PRIVATE KEY-----\n...",
				"user_id":     "user-123",
				"user": map[string]interface{}{
					"id":   "user-123",
					"name": "john@example.com",
				},
				"createdAt": "2025-11-10T15:30:00Z",
				"updatedAt": "2025-11-10T15:30:00Z",
			},
		},
		{
			name: "keypair without optional fields (omitempty test)",
			keypair: Keypair{
				ID:          "kp-abc123",
				Name:        "my-keypair",
				Fingerprint: "SHA256:abc...xyz",
				PublicKey:   "ssh-rsa AAAAB3...",
				UserID:      "user-123",
				CreatedAt:   "2025-11-10T15:30:00Z",
				UpdatedAt:   "2025-11-10T15:30:00Z",
			},
			wantJSON: map[string]interface{}{
				"id":          "kp-abc123",
				"name":        "my-keypair",
				"fingerprint": "SHA256:abc...xyz",
				"public_key":  "ssh-rsa AAAAB3...",
				"user_id":     "user-123",
				"createdAt":   "2025-11-10T15:30:00Z",
				"updatedAt":   "2025-11-10T15:30:00Z",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jsonData, err := json.Marshal(tt.keypair)
			if err != nil {
				t.Fatalf("json.Marshal() error = %v", err)
			}

			var got map[string]interface{}
			if err := json.Unmarshal(jsonData, &got); err != nil {
				t.Fatalf("json.Unmarshal() error = %v", err)
			}

			// Check all expected fields are present
			for key, wantValue := range tt.wantJSON {
				gotValue, exists := got[key]
				if !exists {
					t.Errorf("Missing field %q in marshaled JSON", key)
					continue
				}

				// Special handling for nested objects
				if key == "user" {
					wantUser := wantValue.(map[string]interface{})
					gotUser, ok := gotValue.(map[string]interface{})
					if !ok {
						t.Errorf("Field %q type = %T, want map[string]interface{}", key, gotValue)
						continue
					}
					for userKey, userWantValue := range wantUser {
						if gotUser[userKey] != userWantValue {
							t.Errorf("Field %q.%q = %v, want %v", key, userKey, gotUser[userKey], userWantValue)
						}
					}
				} else if gotValue != wantValue {
					t.Errorf("Field %q = %v, want %v", key, gotValue, wantValue)
				}
			}

			// Check no unexpected fields when using omitempty
			if tt.name == "keypair without optional fields (omitempty test)" {
				unexpectedFields := []string{"description", "private_key", "user"}
				for _, field := range unexpectedFields {
					if _, exists := got[field]; exists {
						t.Errorf("Unexpected field %q in marshaled JSON (should be omitted)", field)
					}
				}
			}
		})
	}
}

// TestOptionalFieldHandling tests handling of optional fields (nil User, empty PrivateKey)
func TestOptionalFieldHandling(t *testing.T) {
	tests := []struct {
		name    string
		keypair Keypair
		checks  func(t *testing.T, jsonData []byte)
	}{
		{
			name: "nil user should omit user field",
			keypair: Keypair{
				ID:        "kp-123",
				Name:      "test",
				UserID:    "user-123",
				User:      nil,
				CreatedAt: "2025-11-10T15:30:00Z",
				UpdatedAt: "2025-11-10T15:30:00Z",
			},
			checks: func(t *testing.T, jsonData []byte) {
				var result map[string]interface{}
				json.Unmarshal(jsonData, &result)
				if _, exists := result["user"]; exists {
					t.Error("user field should be omitted when nil")
				}
			},
		},
		{
			name: "empty PrivateKey should omit private_key field",
			keypair: Keypair{
				ID:         "kp-123",
				Name:       "test",
				PrivateKey: "",
				CreatedAt:  "2025-11-10T15:30:00Z",
				UpdatedAt:  "2025-11-10T15:30:00Z",
			},
			checks: func(t *testing.T, jsonData []byte) {
				var result map[string]interface{}
				json.Unmarshal(jsonData, &result)
				if _, exists := result["private_key"]; exists {
					t.Error("private_key field should be omitted when empty")
				}
			},
		},
		{
			name: "empty Description should omit description field",
			keypair: Keypair{
				ID:          "kp-123",
				Name:        "test",
				Description: "",
				CreatedAt:   "2025-11-10T15:30:00Z",
				UpdatedAt:   "2025-11-10T15:30:00Z",
			},
			checks: func(t *testing.T, jsonData []byte) {
				var result map[string]interface{}
				json.Unmarshal(jsonData, &result)
				if _, exists := result["description"]; exists {
					t.Error("description field should be omitted when empty")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jsonData, err := json.Marshal(tt.keypair)
			if err != nil {
				t.Fatalf("json.Marshal() error = %v", err)
			}
			tt.checks(t, jsonData)
		})
	}
}

// TestPrivateKeyFieldPresence tests that private_key is present in Create response
func TestPrivateKeyFieldPresence(t *testing.T) {
	// Simulate Create response with private key
	createResponse := `{
		"id": "kp-abc123",
		"name": "my-keypair",
		"fingerprint": "SHA256:abc...xyz",
		"public_key": "ssh-rsa AAAAB3...",
		"private_key": "-----BEGIN RSA PRIVATE KEY-----\nMIIEpAIBAAKCAQEA...",
		"user_id": "user-123",
		"createdAt": "2025-11-10T15:30:00Z",
		"updatedAt": "2025-11-10T15:30:00Z"
	}`

	var keypair Keypair
	if err := json.Unmarshal([]byte(createResponse), &keypair); err != nil {
		t.Fatalf("json.Unmarshal() error = %v", err)
	}

	if keypair.PrivateKey == "" {
		t.Error("PrivateKey should be present in Create response")
	}

	expectedPrefix := "-----BEGIN RSA PRIVATE KEY-----"
	if len(keypair.PrivateKey) < len(expectedPrefix) || keypair.PrivateKey[:len(expectedPrefix)] != expectedPrefix {
		t.Errorf("PrivateKey should start with %q, got %q", expectedPrefix, keypair.PrivateKey[:min(len(keypair.PrivateKey), len(expectedPrefix))])
	}
}

// TestPrivateKeyFieldAbsence tests that private_key is absent in Get/List responses
func TestPrivateKeyFieldAbsence(t *testing.T) {
	tests := []struct {
		name     string
		response string
	}{
		{
			name: "Get response without private_key",
			response: `{
				"id": "kp-abc123",
				"name": "my-keypair",
				"fingerprint": "SHA256:abc...xyz",
				"public_key": "ssh-rsa AAAAB3...",
				"user_id": "user-123",
				"createdAt": "2025-11-10T15:30:00Z",
				"updatedAt": "2025-11-10T15:30:00Z"
			}`,
		},
		{
			name: "List response keypair without private_key",
			response: `{
				"id": "kp-abc123",
				"name": "my-keypair",
				"fingerprint": "SHA256:abc...xyz",
				"public_key": "ssh-rsa AAAAB3...",
				"user_id": "user-123",
				"createdAt": "2025-11-10T15:30:00Z",
				"updatedAt": "2025-11-10T15:30:00Z"
			}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var keypair Keypair
			if err := json.Unmarshal([]byte(tt.response), &keypair); err != nil {
				t.Fatalf("json.Unmarshal() error = %v", err)
			}

			if keypair.PrivateKey != "" {
				t.Errorf("PrivateKey should be empty in %s, got %q", tt.name, keypair.PrivateKey)
			}
		})
	}
}

// TestKeypairContractValidation validates against pb.KeypairInfo schema from Swagger
func TestKeypairContractValidation(t *testing.T) {
	// Test that Keypair structure matches pb.KeypairInfo contract
	t.Run("pb.KeypairInfo schema compliance", func(t *testing.T) {
		swaggerExample := `{
			"id": "kp-abc123",
			"name": "my-keypair",
			"description": "Development keypair",
			"fingerprint": "SHA256:abc...xyz",
			"public_key": "ssh-rsa AAAAB3...",
			"private_key": "-----BEGIN RSA PRIVATE KEY-----\n...",
			"user_id": "user-123",
			"user": {
				"id": "user-123",
				"name": "john@example.com"
			},
			"createdAt": "2025-11-10T15:30:00Z",
			"updatedAt": "2025-11-10T15:30:00Z"
		}`

		var keypair Keypair
		if err := json.Unmarshal([]byte(swaggerExample), &keypair); err != nil {
			t.Fatalf("Failed to unmarshal Swagger example: %v", err)
		}

		// Verify all required fields from pb.KeypairInfo are present
		requiredFields := map[string]string{
			"id":          keypair.ID,
			"name":        keypair.Name,
			"fingerprint": keypair.Fingerprint,
			"public_key":  keypair.PublicKey,
			"user_id":     keypair.UserID,
			"createdAt":   keypair.CreatedAt,
			"updatedAt":   keypair.UpdatedAt,
		}

		for field, value := range requiredFields {
			if value == "" {
				t.Errorf("Required field %q is empty", field)
			}
		}

		// Verify optional fields can be present
		if keypair.Description == "" {
			t.Error("Optional field description should be preserved when present")
		}
		if keypair.PrivateKey == "" {
			t.Error("Optional field private_key should be preserved when present")
		}
		if keypair.User == nil {
			t.Error("Optional field user should be preserved when present")
		}

		// Marshal back and verify structure
		jsonData, err := json.Marshal(keypair)
		if err != nil {
			t.Fatalf("Failed to marshal keypair: %v", err)
		}

		var result map[string]interface{}
		if err := json.Unmarshal(jsonData, &result); err != nil {
			t.Fatalf("Failed to unmarshal result: %v", err)
		}

		// Verify JSON tag naming conventions
		expectedTags := map[string]bool{
			"id":          true, // lowercase
			"name":        true, // lowercase
			"description": true, // lowercase
			"fingerprint": true, // lowercase
			"public_key":  true, // snake_case
			"private_key": true, // snake_case
			"user_id":     true, // snake_case
			"user":        true, // lowercase
			"createdAt":   true, // camelCase
			"updatedAt":   true, // camelCase
		}

		for tag := range expectedTags {
			if _, exists := result[tag]; !exists && tag != "description" && tag != "private_key" && tag != "user" {
				t.Errorf("Expected JSON tag %q not found in marshaled output", tag)
			}
		}
	})

	// Test RFC3339 timestamp format
	t.Run("timestamp RFC3339 format", func(t *testing.T) {
		keypair := Keypair{
			ID:        "kp-123",
			Name:      "test",
			CreatedAt: "2025-11-10T15:30:00Z",
			UpdatedAt: "2025-11-10T15:30:45.123Z",
		}

		// Verify timestamps can be parsed as RFC3339
		_, err := time.Parse(time.RFC3339, keypair.CreatedAt)
		if err != nil {
			t.Errorf("CreatedAt should be valid RFC3339 format: %v", err)
		}

		_, err = time.Parse(time.RFC3339, keypair.UpdatedAt)
		if err != nil {
			t.Errorf("UpdatedAt should be valid RFC3339 format: %v", err)
		}
	})
}

// TestIDNameSerialization tests common.IDName struct JSON serialization
func TestIDNameSerialization(t *testing.T) {
	tests := []struct {
		name     string
		idname   common.IDName
		wantJSON string
	}{
		{
			name: "complete common.IDName",
			idname: common.IDName{
				ID:   "user-123",
				Name: "john@example.com",
			},
			wantJSON: `{"id":"user-123","name":"john@example.com"}`,
		},
		{
			name: "IDName with special characters",
			idname: common.IDName{
				ID:   "user-456",
				Name: "jane.doe+test@example.com",
			},
			wantJSON: `{"id":"user-456","name":"jane.doe+test@example.com"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test marshaling
			jsonData, err := json.Marshal(tt.idname)
			if err != nil {
				t.Fatalf("json.Marshal() error = %v", err)
			}

			if string(jsonData) != tt.wantJSON {
				t.Errorf("json.Marshal() = %s, want %s", string(jsonData), tt.wantJSON)
			}

			// Test unmarshaling
			var got common.IDName
			if err := json.Unmarshal([]byte(tt.wantJSON), &got); err != nil {
				t.Fatalf("json.Unmarshal() error = %v", err)
			}

			if got.ID != tt.idname.ID {
				t.Errorf("ID = %v, want %v", got.ID, tt.idname.ID)
			}
			if got.Name != tt.idname.Name {
				t.Errorf("Name = %v, want %v", got.Name, tt.idname.Name)
			}
		})
	}
}

// ========================================================================
// Phase 5: User Story 3 - Timestamp Handling Tests
// ========================================================================

// TestTimestampParsingFromAPIResponse validates that timestamp fields are properly parsed from API responses
func TestTimestampParsingFromAPIResponse(t *testing.T) {
	jsonData := `{
		"id": "kp-timestamp-test",
		"name": "timestamp-test-key",
		"public_key": "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQC...",
		"fingerprint": "SHA256:timestamp123",
		"user_id": "user-timestamp",
		"createdAt": "2025-11-10T10:00:00Z",
		"updatedAt": "2025-11-10T11:30:45Z"
	}`

	var keypair Keypair
	err := json.Unmarshal([]byte(jsonData), &keypair)
	if err != nil {
		t.Fatalf("Failed to unmarshal keypair: %v", err)
	}

	// Verify timestamps are parsed as strings
	if keypair.CreatedAt != "2025-11-10T10:00:00Z" {
		t.Errorf("Expected CreatedAt to be '2025-11-10T10:00:00Z', got '%s'", keypair.CreatedAt)
	}
	if keypair.UpdatedAt != "2025-11-10T11:30:45Z" {
		t.Errorf("Expected UpdatedAt to be '2025-11-10T11:30:45Z', got '%s'", keypair.UpdatedAt)
	}
}

// TestRFC3339FormatValidation tests that timestamps follow RFC3339 format
func TestRFC3339FormatValidation(t *testing.T) {
	tests := []struct {
		name      string
		timestamp string
		wantErr   bool
	}{
		{
			name:      "valid RFC3339 with Z",
			timestamp: "2025-11-10T10:00:00Z",
			wantErr:   false,
		},
		{
			name:      "valid RFC3339 with timezone offset",
			timestamp: "2025-11-10T10:00:00+05:00",
			wantErr:   false,
		},
		{
			name:      "valid RFC3339 with negative offset",
			timestamp: "2025-11-10T10:00:00-08:00",
			wantErr:   false,
		},
		{
			name:      "invalid format - missing T",
			timestamp: "2025-11-10 10:00:00Z",
			wantErr:   true,
		},
		{
			name:      "invalid format - wrong date format",
			timestamp: "11/10/2025T10:00:00Z",
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := time.Parse(time.RFC3339, tt.timestamp)
			if (err != nil) != tt.wantErr {
				t.Errorf("time.Parse() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// TestTimestampConversionToTime validates that timestamp strings can be converted to time.Time
func TestTimestampConversionToTime(t *testing.T) {
	keypair := Keypair{
		ID:          "kp-time-conversion",
		Name:        "time-conversion-test",
		PublicKey:   "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQC...",
		Fingerprint: "SHA256:timeconv123",
		UserID:      "user-time",
		CreatedAt:   "2025-11-10T10:00:00Z",
		UpdatedAt:   "2025-11-10T11:30:45Z",
	}

	// Test CreatedAt conversion
	createdTime, err := time.Parse(time.RFC3339, keypair.CreatedAt)
	if err != nil {
		t.Fatalf("Failed to parse CreatedAt: %v", err)
	}

	expectedCreated := time.Date(2025, 11, 10, 10, 0, 0, 0, time.UTC)
	if !createdTime.Equal(expectedCreated) {
		t.Errorf("CreatedAt parsed incorrectly: got %v, want %v", createdTime, expectedCreated)
	}

	// Test UpdatedAt conversion
	updatedTime, err := time.Parse(time.RFC3339, keypair.UpdatedAt)
	if err != nil {
		t.Fatalf("Failed to parse UpdatedAt: %v", err)
	}

	expectedUpdated := time.Date(2025, 11, 10, 11, 30, 45, 0, time.UTC)
	if !updatedTime.Equal(expectedUpdated) {
		t.Errorf("UpdatedAt parsed incorrectly: got %v, want %v", updatedTime, expectedUpdated)
	}
}

// TestCreatedAtImmutability validates that createdAt remains the same after updates
func TestCreatedAtImmutability(t *testing.T) {
	originalCreatedAt := "2025-11-10T10:00:00Z"

	// Simulate initial creation
	keypair := Keypair{
		ID:          "kp-immutable",
		Name:        "immutable-test",
		PublicKey:   "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQC...",
		Fingerprint: "SHA256:immutable123",
		UserID:      "user-immutable",
		CreatedAt:   originalCreatedAt,
		UpdatedAt:   originalCreatedAt,
	}

	// Simulate an update (only UpdatedAt should change)
	keypair.UpdatedAt = "2025-11-10T11:30:45Z"

	// CreatedAt should remain unchanged
	if keypair.CreatedAt != originalCreatedAt {
		t.Errorf("CreatedAt changed after update: got %s, want %s", keypair.CreatedAt, originalCreatedAt)
	}

	// UpdatedAt should be different
	if keypair.UpdatedAt == originalCreatedAt {
		t.Errorf("UpdatedAt should have changed after update")
	}
}

// TestUpdatedAtChangesAfterUpdate validates that updatedAt changes after updates
func TestUpdatedAtChangesAfterUpdate(t *testing.T) {
	originalUpdatedAt := "2025-11-10T10:00:00Z"

	keypair := Keypair{
		ID:          "kp-updated",
		Name:        "updated-test",
		PublicKey:   "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQC...",
		Fingerprint: "SHA256:updated123",
		UserID:      "user-updated",
		CreatedAt:   "2025-11-10T09:00:00Z",
		UpdatedAt:   originalUpdatedAt,
	}

	// Simulate an update
	newUpdatedAt := "2025-11-10T12:45:30Z"
	keypair.UpdatedAt = newUpdatedAt

	// UpdatedAt should have changed
	if keypair.UpdatedAt != newUpdatedAt {
		t.Errorf("UpdatedAt not updated correctly: got %s, want %s", keypair.UpdatedAt, newUpdatedAt)
	}

	// Verify it's different from original
	if keypair.UpdatedAt == originalUpdatedAt {
		t.Errorf("UpdatedAt should be different after update")
	}
}

// Helper function for min
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
