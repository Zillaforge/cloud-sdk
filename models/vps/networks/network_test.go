package networks

import (
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestNetworkJSONRoundTripFullPayload(t *testing.T) {
	data := loadFullNetworkFixture(t)

	var network Network
	if err := json.Unmarshal(data, &network); err != nil {
		t.Fatalf("failed to unmarshal full network fixture: %v", err)
	}

	assertStringField(t, &network, "ID", "net-full-001")
	assertStringField(t, &network, "Name", "full-network")
	assertStringField(t, &network, "Description", "Network fixture with all fields populated")
	assertStringField(t, &network, "CIDR", "10.42.0.0/24")
	assertBoolField(t, &network, "Bonding", true)
	assertStringField(t, &network, "Gateway", "10.42.0.1")
	assertBoolField(t, &network, "GWState", false)
	assertBoolField(t, &network, "IsDefault", true)
	assertStringSliceField(t, &network, "Nameservers", []string{"1.1.1.1", "8.8.8.8"})
	assertStringField(t, &network, "Namespace", "tenant-alpha")

	project := requirePointerStructField(t, &network, "Project")
	assertStringField(t, project.Interface(), "ID", "proj-001")
	assertStringField(t, project.Interface(), "Name", "Tenant Alpha")
	assertStringField(t, &network, "ProjectID", "proj-001")

	router := requirePointerStructField(t, &network, "Router")
	assertStringField(t, router.Interface(), "ID", "router-123")
	assertStringField(t, router.Interface(), "Name", "alpha-router")
	assertStringField(t, router.Interface(), "Description", "Primary router")
	assertBoolField(t, router.Interface(), "Bonding", false)
	assertBoolField(t, router.Interface(), "IsDefault", false)
	assertBoolField(t, router.Interface(), "Shared", true)
	assertBoolField(t, router.Interface(), "State", true)
	assertStringField(t, router.Interface(), "Status", "ACTIVE")
	assertStringField(t, router.Interface(), "StatusReason", "OK")
	assertStringField(t, router.Interface(), "Namespace", "tenant-alpha")

	routerProject := requirePointerStructField(t, router.Interface(), "Project")
	assertStringField(t, routerProject.Interface(), "ID", "proj-001")
	assertStringField(t, routerProject.Interface(), "Name", "Tenant Alpha")

	routerUser := requirePointerStructField(t, router.Interface(), "User")
	assertStringField(t, routerUser.Interface(), "ID", "user-123")
	assertStringField(t, routerUser.Interface(), "Name", "Alice Ops")

	extNetwork := requirePointerStructField(t, router.Interface(), "ExtNetwork")
	assertStringField(t, extNetwork.Interface(), "ID", "extnet-001")
	assertStringField(t, extNetwork.Interface(), "Name", "public-ext")
	assertStringField(t, extNetwork.Interface(), "Description", "Public external network")
	assertStringField(t, extNetwork.Interface(), "CIDR", "203.0.113.0/24")
	assertStringField(t, extNetwork.Interface(), "Namespace", "global")
	assertStringField(t, extNetwork.Interface(), "SegmentID", "seg-42")
	assertStringField(t, extNetwork.Interface(), "Type", "flat")
	assertBoolField(t, extNetwork.Interface(), "IsDefault", false)

	assertStringField(t, router.Interface(), "ExtNetworkID", "extnet-001")
	assertStringSliceField(t, router.Interface(), "GWAddrs", []string{"192.0.2.1", "192.0.2.2"})
	assertStringField(t, router.Interface(), "CreatedAt", "2025-01-05T12:00:00Z")
	assertStringField(t, router.Interface(), "UpdatedAt", "2025-01-06T12:00:00Z")

	assertStringField(t, &network, "RouterID", "router-123")
	assertBoolField(t, &network, "Shared", true)
	assertStringField(t, &network, "Status", "ACTIVE")
	assertStringField(t, &network, "StatusReason", "OK")
	assertStringField(t, &network, "SubnetID", "subnet-777")

	user := requirePointerStructField(t, &network, "User")
	assertStringField(t, user.Interface(), "ID", "user-123")
	assertStringField(t, user.Interface(), "Name", "Alice Ops")
	assertStringField(t, &network, "UserID", "user-123")

	assertStringField(t, &network, "CreatedAt", "2025-01-05T00:00:00Z")
	assertStringField(t, &network, "UpdatedAt", "2025-01-06T00:00:00Z")

	marshaled, err := json.Marshal(&network)
	if err != nil {
		t.Fatalf("failed to marshal network: %v", err)
	}

	var actual map[string]interface{}
	if err := json.Unmarshal(marshaled, &actual); err != nil {
		t.Fatalf("failed to unmarshal marshaled network: %v", err)
	}

	if _, ok := actual["bonding"]; !ok {
		t.Fatal("expected bonding key to be present in marshaled payload")
	}
	if _, ok := actual["is_default"]; !ok {
		t.Fatal("expected is_default key to be present in marshaled payload")
	}
	if _, ok := actual["shared"]; !ok {
		t.Fatal("expected shared key to be present in marshaled payload")
	}

	var roundtrip Network
	if err := json.Unmarshal(marshaled, &roundtrip); err != nil {
		t.Fatalf("failed to unmarshal roundtrip network: %v", err)
	}

	if !reflect.DeepEqual(network, roundtrip) {
		t.Fatalf("roundtrip network mismatch\nexpected: %#v\nactual: %#v", network, roundtrip)
	}
}

func TestNetworkJSONOptionalFieldsOmitted(t *testing.T) {
	minimal := []byte(`{
		"id": "net-min-001",
		"name": "minimal-network",
		"cidr": "10.99.0.0/24",
		"createdAt": "2025-02-01T00:00:00Z"
	}`)

	var network Network
	if err := json.Unmarshal(minimal, &network); err != nil {
		t.Fatalf("failed to unmarshal minimal payload: %v", err)
	}

	if network.Bonding {
		t.Fatal("expected bonding to default to false")
	}
	if network.Gateway != "" {
		t.Fatalf("expected empty gateway, got %q", network.Gateway)
	}
	if network.GWState {
		t.Fatal("expected deprecated gw_state to default to false")
	}
	if network.IsDefault {
		t.Fatal("expected is_default to default to false")
	}
	if len(network.Nameservers) != 0 {
		t.Fatalf("expected no nameservers, got %v", network.Nameservers)
	}
	if network.Namespace != "" {
		t.Fatalf("expected empty namespace, got %q", network.Namespace)
	}
	if network.Project != nil {
		t.Fatal("expected nil project reference")
	}
	if network.Router != nil {
		t.Fatal("expected nil router reference")
	}
	if network.Shared {
		t.Fatal("expected shared to default to false")
	}
	if network.Status != "" {
		t.Fatalf("expected empty status, got %q", network.Status)
	}
	if network.SubnetID != "" {
		t.Fatalf("expected empty subnet_id, got %q", network.SubnetID)
	}
	if network.User != nil {
		t.Fatal("expected nil user reference")
	}
	if network.UserID != "" {
		t.Fatalf("expected empty user_id, got %q", network.UserID)
	}

	marshaled, err := json.Marshal(&network)
	if err != nil {
		t.Fatalf("failed to marshal minimal network: %v", err)
	}

	var payload map[string]interface{}
	if err := json.Unmarshal(marshaled, &payload); err != nil {
		t.Fatalf("failed to unmarshal minimal network payload: %v", err)
	}

	for _, key := range []string{"bonding", "gateway", "gw_state", "is_default", "nameservers", "namespace", "project", "router", "shared", "status", "status_reason", "subnet_id", "user", "user_id"} {
		if _, exists := payload[key]; exists {
			t.Fatalf("expected key %q to be omitted, but found in payload: %v", key, payload[key])
		}
	}
}

func TestNetworkCreateRequestJSONWithOptionalFields(t *testing.T) {
	req := &NetworkCreateRequest{
		Name:        "gateway-network",
		Description: "Network with optional fields",
		CIDR:        "10.50.0.0/24",
		Gateway:     "10.50.0.1",
		RouterID:    "router-optional-1",
	}

	data, err := json.Marshal(req)
	if err != nil {
		t.Fatalf("failed to marshal create request: %v", err)
	}

	var payload map[string]interface{}
	if err := json.Unmarshal(data, &payload); err != nil {
		t.Fatalf("failed to unmarshal create request payload: %v", err)
	}

	if payload["name"] != "gateway-network" {
		t.Fatalf("expected name 'gateway-network', got %v", payload["name"])
	}
	if payload["description"] != "Network with optional fields" {
		t.Fatalf("expected description 'Network with optional fields', got %v", payload["description"])
	}
	if payload["cidr"] != "10.50.0.0/24" {
		t.Fatalf("expected cidr '10.50.0.0/24', got %v", payload["cidr"])
	}
	if payload["gateway"] != "10.50.0.1" {
		t.Fatalf("expected gateway '10.50.0.1', got %v", payload["gateway"])
	}
	if payload["router_id"] != "router-optional-1" {
		t.Fatalf("expected router_id 'router-optional-1', got %v", payload["router_id"])
	}
}

func TestNetworkCreateRequestJSONOmitOptionalFields(t *testing.T) {
	req := &NetworkCreateRequest{
		Name: "minimal-network",
		CIDR: "10.60.0.0/24",
	}

	data, err := json.Marshal(req)
	if err != nil {
		t.Fatalf("failed to marshal minimal create request: %v", err)
	}

	var payload map[string]interface{}
	if err := json.Unmarshal(data, &payload); err != nil {
		t.Fatalf("failed to unmarshal minimal create request payload: %v", err)
	}

	if _, ok := payload["gateway"]; ok {
		t.Fatalf("unexpected gateway field present: %v", payload["gateway"])
	}
	if _, ok := payload["router_id"]; ok {
		t.Fatalf("unexpected router_id field present: %v", payload["router_id"])
	}
	if _, ok := payload["description"]; ok {
		t.Fatalf("unexpected description field present: %v", payload["description"])
	}
}

func loadFullNetworkFixture(t *testing.T) []byte {
	t.Helper()

	path := filepath.Join("..", "..", "..", "modules", "vps", "networks", "test", "testdata", "network_full.json")
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("failed to read full network fixture: %v", err)
	}
	return data
}

func requireStructValue(t *testing.T, obj interface{}) reflect.Value {
	t.Helper()

	v := reflect.ValueOf(obj)
	if v.Kind() == reflect.Ptr {
		if v.IsNil() {
			t.Fatal("nil pointer encountered while accessing struct value")
		}
		v = v.Elem()
	}
	if v.Kind() != reflect.Struct {
		t.Fatalf("expected struct value, got %s", v.Kind())
	}
	return v
}

func requireStructField(t *testing.T, obj interface{}, fieldName string) reflect.Value {
	t.Helper()

	v := requireStructValue(t, obj)
	field := v.FieldByName(fieldName)
	if !field.IsValid() {
		t.Fatalf("missing field %s", fieldName)
	}
	return field
}

func requirePointerStructField(t *testing.T, obj interface{}, fieldName string) reflect.Value {
	t.Helper()

	field := requireStructField(t, obj, fieldName)
	if field.Kind() != reflect.Ptr {
		t.Fatalf("expected pointer field %s, got %s", fieldName, field.Kind())
	}
	if field.IsNil() {
		t.Fatalf("expected field %s to be non-nil", fieldName)
	}

	elem := field.Elem()
	if elem.Kind() != reflect.Struct {
		t.Fatalf("expected struct pointer for field %s, got %s", fieldName, elem.Kind())
	}
	return elem
}

func assertStringField(t *testing.T, obj interface{}, fieldName, expected string) {
	t.Helper()

	field := requireStructField(t, obj, fieldName)
	if field.Kind() != reflect.String {
		t.Fatalf("expected string field %s, got %s", fieldName, field.Kind())
	}
	if field.String() != expected {
		t.Fatalf("field %s mismatch: expected %s, got %s", fieldName, expected, field.String())
	}
}

func assertBoolField(t *testing.T, obj interface{}, fieldName string, expected bool) {
	t.Helper()

	field := requireStructField(t, obj, fieldName)
	if field.Kind() != reflect.Bool {
		t.Fatalf("expected bool field %s, got %s", fieldName, field.Kind())
	}
	if field.Bool() != expected {
		t.Fatalf("field %s mismatch: expected %v, got %v", fieldName, expected, field.Bool())
	}
}

func assertStringSliceField(t *testing.T, obj interface{}, fieldName string, expected []string) {
	t.Helper()

	field := requireStructField(t, obj, fieldName)
	if field.Kind() != reflect.Slice {
		t.Fatalf("expected slice field %s, got %s", fieldName, field.Kind())
	}

	actual := make([]string, field.Len())
	for i := 0; i < field.Len(); i++ {
		elem := field.Index(i)
		if elem.Kind() != reflect.String {
			t.Fatalf("expected string slice element for field %s, got %s", fieldName, elem.Kind())
		}
		actual[i] = elem.String()
	}

	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("field %s mismatch: expected %v, got %v", fieldName, expected, actual)
	}
}
