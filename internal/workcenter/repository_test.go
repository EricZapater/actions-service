package workcenter

import (
	"actions-service/internal/models"
	"actions-service/internal/state"
	"context"
	"testing"

	"github.com/alicebob/miniredis/v2"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

func TestRepositorySetFindList(t *testing.T) {
	ctx := context.Background()
	mr, err := miniredis.Run()
	if err != nil {
		t.Fatalf("failed to start miniredis: %v", err)
	}
	defer mr.Close()

	client := redis.NewClient(&redis.Options{Addr: mr.Addr()})
	s := state.New()
	repo := NewWorkcenterRepository(s, client)

	id := uuid.New()
	dto := models.WorkcenterDTO{WorkcenterID: id, WorkcenterName: "WC"}

	if err := repo.Set(ctx, id.String(), dto); err != nil {
		t.Fatalf("Set failed: %v", err)
	}

	found, source, err := repo.FindByID(ctx, id.String())
	if err != nil {
		t.Fatalf("FindByID failed: %v", err)
	}
	if source != models.SourceMemory {
		t.Fatalf("expected memory source, got %v", source)
	}
	if found.WorkcenterName != "WC" {
		t.Fatalf("unexpected workcenter name: %s", found.WorkcenterName)
	}

	delete(s.Workcenters, id.String())
	found, source, err = repo.FindByID(ctx, id.String())
	if err != nil {
		t.Fatalf("redis FindByID failed: %v", err)
	}
	if source != models.SourceRedis {
		t.Fatalf("expected redis source, got %v", source)
	}

	items, err := repo.List(ctx)
	if err != nil {
		t.Fatalf("List failed: %v", err)
	}
	if len(items) != 1 {
		t.Fatalf("expected one item, got %d", len(items))
	}
}


