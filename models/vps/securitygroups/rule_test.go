package securitygroups

import (
	"encoding/json"
	"testing"
)

// TestProtocolMarshaling tests Protocol type marshaling to JSON.
func TestProtocolMarshaling(t *testing.T) {
	tests := []struct {
		name     string
		protocol Protocol
		want     string
	}{
		{"TCP", ProtocolTCP, `"tcp"`},
		{"UDP", ProtocolUDP, `"udp"`},
		{"ICMP", ProtocolICMP, `"icmp"`},
		{"Any", ProtocolAny, `"any"`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := json.Marshal(tt.protocol)
			if err != nil {
				t.Fatalf("Marshal() error = %v", err)
			}
			if string(data) != tt.want {
				t.Errorf("Marshal() = %s, want %s", data, tt.want)
			}
		})
	}
}

// TestProtocolUnmarshaling tests Protocol type unmarshaling from JSON.
func TestProtocolUnmarshaling(t *testing.T) {
	tests := []struct {
		name    string
		json    string
		want    Protocol
		wantErr bool
	}{
		{"TCP", `"tcp"`, ProtocolTCP, false},
		{"UDP", `"udp"`, ProtocolUDP, false},
		{"ICMP", `"icmp"`, ProtocolICMP, false},
		{"Any", `"any"`, ProtocolAny, false},
		{"Custom value", `"custom"`, Protocol("custom"), false}, // Go allows any string value
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got Protocol
			err := json.Unmarshal([]byte(tt.json), &got)
			if (err != nil) != tt.wantErr {
				t.Fatalf("Unmarshal() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && got != tt.want {
				t.Errorf("Unmarshal() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestDirectionMarshaling tests Direction type marshaling to JSON.
func TestDirectionMarshaling(t *testing.T) {
	tests := []struct {
		name      string
		direction Direction
		want      string
	}{
		{"Ingress", DirectionIngress, `"ingress"`},
		{"Egress", DirectionEgress, `"egress"`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := json.Marshal(tt.direction)
			if err != nil {
				t.Fatalf("Marshal() error = %v", err)
			}
			if string(data) != tt.want {
				t.Errorf("Marshal() = %s, want %s", data, tt.want)
			}
		})
	}
}

// TestDirectionUnmarshaling tests Direction type unmarshaling from JSON.
func TestDirectionUnmarshaling(t *testing.T) {
	tests := []struct {
		name    string
		json    string
		want    Direction
		wantErr bool
	}{
		{"Ingress", `"ingress"`, DirectionIngress, false},
		{"Egress", `"egress"`, DirectionEgress, false},
		{"Custom value", `"custom"`, Direction("custom"), false}, // Go allows any string value
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got Direction
			err := json.Unmarshal([]byte(tt.json), &got)
			if (err != nil) != tt.wantErr {
				t.Fatalf("Unmarshal() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && got != tt.want {
				t.Errorf("Unmarshal() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestSecurityGroupRuleMarshalingWithCustomTypes tests that SecurityGroupRule marshals correctly with custom types.
func TestSecurityGroupRuleMarshalingWithCustomTypes(t *testing.T) {
	rule := SecurityGroupRule{
		ID:         "rule-123",
		Direction:  DirectionIngress,
		Protocol:   ProtocolTCP,
		PortMin:    80,
		PortMax:    80,
		RemoteCIDR: "0.0.0.0/0",
	}

	data, err := json.Marshal(rule)
	if err != nil {
		t.Fatalf("Marshal() error = %v", err)
	}

	// Unmarshal back to verify
	var got SecurityGroupRule
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("Unmarshal() error = %v", err)
	}

	if got.Direction != DirectionIngress {
		t.Errorf("Direction = %v, want %v", got.Direction, DirectionIngress)
	}
	if got.Protocol != ProtocolTCP {
		t.Errorf("Protocol = %v, want %v", got.Protocol, ProtocolTCP)
	}
}

// TestSecurityGroupRuleCreateRequestMarshalingWithCustomTypes tests that request marshals correctly.
func TestSecurityGroupRuleCreateRequestMarshalingWithCustomTypes(t *testing.T) {
	portMin := 443
	portMax := 443
	req := SecurityGroupRuleCreateRequest{
		Direction:  DirectionEgress,
		Protocol:   ProtocolUDP,
		PortMin:    &portMin,
		PortMax:    &portMax,
		RemoteCIDR: "10.0.0.0/8",
	}

	data, err := json.Marshal(req)
	if err != nil {
		t.Fatalf("Marshal() error = %v", err)
	}

	// Unmarshal back to verify
	var got SecurityGroupRuleCreateRequest
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("Unmarshal() error = %v", err)
	}

	if got.Direction != DirectionEgress {
		t.Errorf("Direction = %v, want %v", got.Direction, DirectionEgress)
	}
	if got.Protocol != ProtocolUDP {
		t.Errorf("Protocol = %v, want %v", got.Protocol, ProtocolUDP)
	}
}

// TestSecurityGroupRuleWithAllFields tests SecurityGroupRule with all fields populated.
func TestSecurityGroupRuleWithAllFields(t *testing.T) {
	jsonData := `{
		"id": "rule-xyz789",
		"direction": "ingress",
		"protocol": "tcp",
		"port_min": 22,
		"port_max": 22,
		"remote_cidr": "192.168.1.0/24"
	}`

	var rule SecurityGroupRule
	err := json.Unmarshal([]byte(jsonData), &rule)
	if err != nil {
		t.Fatalf("Unmarshal() error = %v", err)
	}

	// Verify all fields
	if rule.ID != "rule-xyz789" {
		t.Errorf("ID = %s, want rule-xyz789", rule.ID)
	}
	if rule.Direction != DirectionIngress {
		t.Errorf("Direction = %v, want %v", rule.Direction, DirectionIngress)
	}
	if rule.Protocol != ProtocolTCP {
		t.Errorf("Protocol = %v, want %v", rule.Protocol, ProtocolTCP)
	}
	// Verify ports are plain int (not pointers)
	if rule.PortMin != 22 {
		t.Errorf("PortMin = %d, want 22", rule.PortMin)
	}
	if rule.PortMax != 22 {
		t.Errorf("PortMax = %d, want 22", rule.PortMax)
	}
	if rule.RemoteCIDR != "192.168.1.0/24" {
		t.Errorf("RemoteCIDR = %s, want 192.168.1.0/24", rule.RemoteCIDR)
	}
}

// TestSecurityGroupRuleWithICMPProtocol tests SecurityGroupRule with ICMP protocol (ports should be 0).
func TestSecurityGroupRuleWithICMPProtocol(t *testing.T) {
	jsonData := `{
		"id": "rule-icmp-001",
		"direction": "ingress",
		"protocol": "icmp",
		"port_min": 0,
		"port_max": 0,
		"remote_cidr": "0.0.0.0/0"
	}`

	var rule SecurityGroupRule
	err := json.Unmarshal([]byte(jsonData), &rule)
	if err != nil {
		t.Fatalf("Unmarshal() error = %v", err)
	}

	if rule.Protocol != ProtocolICMP {
		t.Errorf("Protocol = %v, want %v", rule.Protocol, ProtocolICMP)
	}
	// For ICMP, ports should be 0
	if rule.PortMin != 0 {
		t.Errorf("PortMin = %d, want 0 for ICMP", rule.PortMin)
	}
	if rule.PortMax != 0 {
		t.Errorf("PortMax = %d, want 0 for ICMP", rule.PortMax)
	}
}

// TestSecurityGroupRuleCreateRequestJSONMarshaling tests SecurityGroupRuleCreateRequest marshals correctly.
func TestSecurityGroupRuleCreateRequestJSONMarshaling(t *testing.T) {
	tests := []struct {
		name         string
		req          SecurityGroupRuleCreateRequest
		wantFields   []string
		unwantFields []string
	}{
		{
			name: "TCP rule with ports",
			req: SecurityGroupRuleCreateRequest{
				Direction:  DirectionIngress,
				Protocol:   ProtocolTCP,
				PortMin:    intPtr(80),
				PortMax:    intPtr(80),
				RemoteCIDR: "0.0.0.0/0",
			},
			wantFields:   []string{"direction", "protocol", "port_min", "port_max", "remote_cidr"},
			unwantFields: []string{},
		},
		{
			name: "UDP rule with port range",
			req: SecurityGroupRuleCreateRequest{
				Direction:  DirectionEgress,
				Protocol:   ProtocolUDP,
				PortMin:    intPtr(8000),
				PortMax:    intPtr(9000),
				RemoteCIDR: "192.168.0.0/16",
			},
			wantFields:   []string{"direction", "protocol", "port_min", "port_max", "remote_cidr"},
			unwantFields: []string{},
		},
		{
			name: "ICMP rule without ports",
			req: SecurityGroupRuleCreateRequest{
				Direction:  DirectionIngress,
				Protocol:   ProtocolICMP,
				RemoteCIDR: "0.0.0.0/0",
			},
			wantFields:   []string{"direction", "protocol", "remote_cidr"},
			unwantFields: []string{"port_min", "port_max"}, // omitempty should exclude these
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := json.Marshal(tt.req)
			if err != nil {
				t.Fatalf("Marshal() error = %v", err)
			}

			var raw map[string]interface{}
			if err := json.Unmarshal(data, &raw); err != nil {
				t.Fatalf("Unmarshal to map error = %v", err)
			}

			// Verify expected fields present
			for _, field := range tt.wantFields {
				if _, exists := raw[field]; !exists {
					t.Errorf("JSON should contain '%s' field", field)
				}
			}

			// Verify unwanted fields absent (due to omitempty)
			for _, field := range tt.unwantFields {
				if _, exists := raw[field]; exists {
					t.Errorf("JSON should NOT contain '%s' field (omitempty)", field)
				}
			}

			// Verify custom types marshal to strings
			if direction, ok := raw["direction"].(string); !ok {
				t.Error("direction should marshal as string")
			} else if direction != string(tt.req.Direction) {
				t.Errorf("direction = %s, want %s", direction, string(tt.req.Direction))
			}

			if protocol, ok := raw["protocol"].(string); !ok {
				t.Error("protocol should marshal as string")
			} else if protocol != string(tt.req.Protocol) {
				t.Errorf("protocol = %s, want %s", protocol, string(tt.req.Protocol))
			}
		})
	}
}
