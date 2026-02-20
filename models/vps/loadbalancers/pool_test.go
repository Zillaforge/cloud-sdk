package loadbalancers

import (
	"encoding/json"
	"testing"
)

func TestPoolJSONRoundtrip(t *testing.T) {
	original := `{
		"id": "pool-123",
		"name": "my-pool",
		"protocol": "HTTP",
		"member_port": 80,
		"method": "ROUND_ROBIN",
		"status": "ACTIVE",
		"members": [
			{"id": "member-1", "name": "vm1", "address": "192.168.1.10", "protocol_port": 80, "status": "ACTIVE"}
		]
	}`

	var pool Pool
	if err := json.Unmarshal([]byte(original), &pool); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}

	if pool.ID != "pool-123" {
		t.Errorf("expected ID pool-123, got %s", pool.ID)
	}
	if pool.Method != PoolMethodRoundRobin {
		t.Errorf("expected Method ROUND_ROBIN, got %s", pool.Method)
	}
	if len(pool.Members) != 1 {
		t.Fatalf("expected 1 member, got %d", len(pool.Members))
	}
	if pool.Members[0].Address != "192.168.1.10" {
		t.Errorf("expected member address 192.168.1.10, got %s", pool.Members[0].Address)
	}

	// Re-marshal
	data, err := json.Marshal(&pool)
	if err != nil {
		t.Fatalf("failed to marshal: %v", err)
	}

	var roundtrip Pool
	if err := json.Unmarshal(data, &roundtrip); err != nil {
		t.Fatalf("failed to unmarshal roundtrip: %v", err)
	}

	if roundtrip.ID != pool.ID || roundtrip.Method != pool.Method {
		t.Error("roundtrip mismatch")
	}
}

func TestPoolCreateRequestValidate(t *testing.T) {
	tests := []struct {
		name    string
		req     PoolCreateRequest
		wantErr string
	}{
		{
			name:    "valid request",
			req:     PoolCreateRequest{Name: "pool", Protocol: PoolProtocolHTTP, MemberPort: 80, Method: PoolMethodRoundRobin},
			wantErr: "",
		},
		{
			name:    "missing name",
			req:     PoolCreateRequest{Protocol: PoolProtocolHTTP, MemberPort: 80, Method: PoolMethodRoundRobin},
			wantErr: "name is required",
		},
		{
			name:    "missing protocol",
			req:     PoolCreateRequest{Name: "pool", MemberPort: 80, Method: PoolMethodRoundRobin},
			wantErr: "protocol is required",
		},
		{
			name:    "invalid protocol",
			req:     PoolCreateRequest{Name: "pool", Protocol: PoolProtocol("GRPC"), MemberPort: 80, Method: PoolMethodRoundRobin},
			wantErr: "protocol must be one of HTTP, HTTPS, TCP, UDP",
		},
		{
			name:    "missing member_port",
			req:     PoolCreateRequest{Name: "pool", Protocol: PoolProtocolHTTP, Method: PoolMethodRoundRobin},
			wantErr: "member_port is required",
		},
		{
			name:    "missing method",
			req:     PoolCreateRequest{Name: "pool", Protocol: PoolProtocolHTTP, MemberPort: 80},
			wantErr: "method is required",
		},
		{
			name:    "invalid method",
			req:     PoolCreateRequest{Name: "pool", Protocol: PoolProtocolHTTP, MemberPort: 80, Method: PoolMethod("RANDOM")},
			wantErr: "method must be one of ROUND_ROBIN, LEAST_CONNECTIONS, SOURCE_IP",
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

func TestMemberCreateRequestValidate(t *testing.T) {
	tests := []struct {
		name    string
		req     MemberCreateRequest
		wantErr string
	}{
		{
			name:    "valid request",
			req:     MemberCreateRequest{Address: "192.168.1.10"},
			wantErr: "",
		},
		{
			name:    "missing address",
			req:     MemberCreateRequest{},
			wantErr: "address is required",
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

func TestPoolMethodValid(t *testing.T) {
	if !PoolMethodRoundRobin.Valid() {
		t.Error("ROUND_ROBIN should be valid")
	}
	if !PoolMethodLeastConnections.Valid() {
		t.Error("LEAST_CONNECTIONS should be valid")
	}
	if !PoolMethodSourceIP.Valid() {
		t.Error("SOURCE_IP should be valid")
	}
	if PoolMethod("RANDOM").Valid() {
		t.Error("RANDOM should not be valid")
	}
}

func TestMemberJSONRoundtrip(t *testing.T) {
	original := `{"id": "m-1", "name": "vm1", "address": "10.0.0.1", "protocol_port": 80, "status": "ACTIVE"}`

	var member Member
	if err := json.Unmarshal([]byte(original), &member); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}

	if member.ID != "m-1" {
		t.Errorf("expected ID m-1, got %s", member.ID)
	}
	if member.Address != "10.0.0.1" {
		t.Errorf("expected Address 10.0.0.1, got %s", member.Address)
	}
}

func TestPoolListResponseJSON(t *testing.T) {
	jsonData := `{"pools": [{"id": "p-1", "name": "pool1", "protocol": "HTTP", "member_port": 80, "method": "ROUND_ROBIN"}]}`

	var response PoolListResponse
	if err := json.Unmarshal([]byte(jsonData), &response); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}

	if len(response.Pools) != 1 {
		t.Fatalf("expected 1 pool, got %d", len(response.Pools))
	}
	if response.Pools[0].Name != "pool1" {
		t.Errorf("expected Name pool1, got %s", response.Pools[0].Name)
	}
}
