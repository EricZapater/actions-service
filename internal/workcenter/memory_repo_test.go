package workcenter

import (
	"actions-service/internal/models"
	"actions-service/internal/state"
	"context"
	"testing"

	"github.com/google/uuid"
)

func TestMemoryRepoSetAndFind(t *testing.T) {
	ctx := context.Background()
	s := state.New()
	repo := NewMemoryRepository(s)

	id := uuid.New().String()
	dto := models.WorkcenterDTO{WorkcenterID: uuid.MustParse(id), WorkcenterName: "WC"}

	if err := repo.Set(ctx, id, dto); err != nil {
		t.Fatalf("Set failed: %v", err)
	}

	found, err := repo.FindByID(ctx, id)
	if err != nil {
		t.Fatalf("FindByID failed: %v", err)
	}
	if found.WorkcenterName != "WC" {
		t.Fatalf("unexpected workcenter name: %s", found.WorkcenterName)
	}

	if _, err := repo.FindByID(ctx, uuid.NewString()); err == nil {
		t.Fatalf("expected error for missing workcenter")
	}
}


