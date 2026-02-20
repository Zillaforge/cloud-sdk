package loadbalancers

import (
	"encoding/json"
	"testing"
)

func TestListenerJSONRoundtrip(t *testing.T) {
	original := `{
		"id": "listener-123",
		"name": "my-listener",
		"lb_id": "lb-456",
		"protocol": "HTTP",
		"protocol_port": 80,
		"status": "ACTIVE",
		"allowed_cidrs": ["10.0.0.0/24"],
		"insert_headers": {"X-Forwarded-For": "true"},
		"timeout_client_data": 50000,
		"timeout_member_connect": 5000,
		"timeout_member_data": 50000,
		"timeout_tcp_inspect": 0
	}`

	var listener Listener
	if err := json.Unmarshal([]byte(original), &listener); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}

	if listener.ID != "listener-123" {
		t.Errorf("expected ID listener-123, got %s", listener.ID)
	}
	if listener.Protocol != ListenerProtocolHTTP {
		t.Errorf("expected Protocol HTTP, got %s", listener.Protocol)
	}
	if listener.ProtocolPort != 80 {
		t.Errorf("expected ProtocolPort 80, got %d", listener.ProtocolPort)
	}
	if len(listener.AllowedCIDRs) != 1 || listener.AllowedCIDRs[0] != "10.0.0.0/24" {
		t.Errorf("expected AllowedCIDRs [10.0.0.0/24], got %v", listener.AllowedCIDRs)
	}
	if listener.InsertHeaders["X-Forwarded-For"] != "true" {
		t.Errorf("expected InsertHeaders X-Forwarded-For=true, got %v", listener.InsertHeaders)
	}

	// Re-marshal
	data, err := json.Marshal(&listener)
	if err != nil {
		t.Fatalf("failed to marshal: %v", err)
	}

	var roundtrip Listener
	if err := json.Unmarshal(data, &roundtrip); err != nil {
		t.Fatalf("failed to unmarshal roundtrip: %v", err)
	}

	if roundtrip.ID != listener.ID || roundtrip.Protocol != listener.Protocol {
		t.Error("roundtrip mismatch")
	}
}

func TestListenerCreateRequestValidate(t *testing.T) {
	tests := []struct {
		name    string
		req     ListenerCreateRequest
		wantErr string
	}{
		{
			name:    "valid request",
			req:     ListenerCreateRequest{Protocol: ListenerProtocolHTTP, ProtocolPort: 80},
			wantErr: "",
		},
		{
			name:    "missing protocol",
			req:     ListenerCreateRequest{ProtocolPort: 80},
			wantErr: "protocol is required",
		},
		{
			name:    "invalid protocol",
			req:     ListenerCreateRequest{Protocol: ListenerProtocol("INVALID"), ProtocolPort: 80},
			wantErr: "protocol must be one of HTTP, HTTPS, TERMINATED_HTTPS, TCP, UDP",
		},
		{
			name:    "missing protocol_port",
			req:     ListenerCreateRequest{Protocol: ListenerProtocolHTTP},
			wantErr: "protocol_port is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.req.Validate()
			if tt.wantErr == "" {
				if err != nil {
					t.Errorf("expected no error, got: %v", err)
				}
			} else {
				if err == nil {
					t.Errorf("expected error containing %q, got nil", tt.wantErr)
				} else if err.Error() != tt.wantErr {
					t.Errorf("expected error %q, got %q", tt.wantErr, err.Error())
				}
			}
		})
	}
}

func TestListenerListResponseJSON(t *testing.T) {
	jsonData := `{"listeners": [{"id": "l-1", "protocol": "HTTP", "protocol_port": 80}]}`

	var response ListenerListResponse
	if err := json.Unmarshal([]byte(jsonData), &response); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}

	if len(response.Listeners) != 1 {
		t.Fatalf("expected 1 listener, got %d", len(response.Listeners))
	}
	if response.Listeners[0].Protocol != ListenerProtocolHTTP {
		t.Errorf("expected Protocol HTTP, got %s", response.Listeners[0].Protocol)
	}
}
