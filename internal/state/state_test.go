package state

import "testing"

func TestNewStateInitializesMaps(t *testing.T) {
	s := New()
	if s == nil {
		t.Fatalf("expected state instance")
	}
	if s.Workcenters == nil || s.Shifts == nil || s.Operators == nil || s.Statuses == nil {
		t.Fatalf("expected all state maps to be initialized")
	}
}

func TestGetStateReturnsSelf(t *testing.T) {
	s := New()
	if s.GetState() != s {
		t.Fatalf("GetState should return same pointer")
	}
}


