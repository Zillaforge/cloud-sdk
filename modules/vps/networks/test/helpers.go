package networks_test

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

// loadFixtureBytes reads a fixture file from the testdata directory.
func loadFixtureBytes(t *testing.T, name string) []byte {
	t.Helper()

	path := filepath.Join("testdata", name)
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("failed to read fixture %s: %v", name, err)
	}

	return data
}

func requireStructValue(t *testing.T, obj interface{}) reflect.Value {
	t.Helper()

	v := reflect.ValueOf(obj)
	if v.Kind() == reflect.Ptr {
		if v.IsNil() {
			t.Fatal("received nil pointer when struct value expected")
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
		t.Fatalf("expected non-nil field %s", fieldName)
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
			t.Fatalf("expected string element for field %s, got %s", fieldName, elem.Kind())
		}
		actual[i] = elem.String()
	}

	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("field %s mismatch: expected %v, got %v", fieldName, expected, actual)
	}
}
