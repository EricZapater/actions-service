package operator

import (
	"actions-service/internal/models"
	"actions-service/internal/state"
	"context"
	"testing"

	"github.com/google/uuid"
)

func TestMemoryRepoSetFindList(t *testing.T) {
	ctx := context.Background()
	s := state.New()
	repo := NewMemoryRepository(s)

	id := uuid.New().String()
	operator := models.OperatorDTO{
		OperatorID: uuid.MustParse(id),
		OperatorName: "Alice",
	}

	if err := repo.Set(ctx, id, operator); err != nil {
		t.Fatalf("Set failed: %v", err)
	}

	found, err := repo.FindByID(ctx, id)
	if err != nil {
		t.Fatalf("FindByID failed: %v", err)
	}
	if found.OperatorName != "Alice" {
		t.Fatalf("unexpected operator name: %s", found.OperatorName)
	}

	items, err := repo.List(ctx)
	if err != nil {
		t.Fatalf("List failed: %v", err)
	}
	if len(items) != 1 || items[0].OperatorID != operator.OperatorID {
		t.Fatalf("unexpected list result: %+v", items)
	}

	if _, err := repo.FindByID(ctx, uuid.NewString()); err == nil {
		t.Fatalf("expected error for missing operator")
	}
}

func TestMemoryRepoSetWorkcenterDTO(t *testing.T) {
	ctx := context.Background()
	s := state.New()
	repo := NewMemoryRepository(s)

	id := uuid.New().String()
	wc := models.WorkcenterDTO{WorkcenterID: uuid.MustParse(id), WorkcenterName: "WC"}

	if err := repo.SetWorkcenterDTO(ctx, id, wc); err != nil {
		t.Fatalf("SetWorkcenterDTO failed: %v", err)
	}
	stored, ok := s.Workcenters[id]
	if !ok || stored.WorkcenterName != "WC" {
		t.Fatalf("workcenter not stored correctly")
	}
}



