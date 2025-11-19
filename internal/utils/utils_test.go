package utils

import (
	"testing"
	"github.com/google/uuid"
)

func TestGetUUIDFromJsonString(t *testing.T) {
	validID := uuid.New()
	input := map[string]interface{}{"id": validID.String()}
	result, err := GetUUIDFromJsonString(input, "id")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if result != validID {
		t.Fatalf("expected %s, got %s", validID, result)
	}

	_, err = GetUUIDFromJsonString(input, "missing")
	if err == nil {
		t.Fatalf("expected error for missing key")
	}

	input["id"] = 123
	_, err = GetUUIDFromJsonString(input, "id")
	if err == nil {
		t.Fatalf("expected error for non-string value")
	}

	input["id"] = "invalid"
	_, err = GetUUIDFromJsonString(input, "id")
	if err == nil {
		t.Fatalf("expected error for invalid uuid")
	}
}

func TestStringAsAPointer(t *testing.T) {
	value := "test"
	ptr := StringAsAPointer(value)
	if ptr == nil {
		t.Fatalf("expected non-nil pointer")
	}
	if *ptr != value {
		t.Fatalf("expected %s, got %s", value, *ptr)
	}
}


