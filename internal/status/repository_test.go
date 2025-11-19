package status

import (
	"actions-service/internal/models"
	"actions-service/internal/state"
	"context"
	"testing"

	"github.com/alicebob/miniredis/v2"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

func TestStatusRepositorySetAndFind(t *testing.T) {
	ctx := context.Background()
	mr, err := miniredis.Run()
	if err != nil {
		t.Fatalf("failed to start miniredis: %v", err)
	}
	defer mr.Close()

	client := redis.NewClient(&redis.Options{Addr: mr.Addr()})
	s := state.New()
	repo := NewStatusRepository(s, client)

	statusID := uuid.New()
	workcenterID := uuid.New()
	dto := models.StatusDTO{
		StatusId:     statusID,
		WorkcenterId: workcenterID,
		Description:  "Running",
	}
	key := workcenterID.String() + ":" + statusID.String()

	if err := repo.Set(ctx, key, dto); err != nil {
		t.Fatalf("Set failed: %v", err)
	}

	found, source, err := repo.FindByID(ctx, key)
	if err != nil {
		t.Fatalf("FindByID failed: %v", err)
	}
	if source != models.SourceMemory {
		t.Fatalf("expected memory source, got %v", source)
	}
	if found.Description != "Running" {
		t.Fatalf("unexpected status description: %s", found.Description)
	}

	delete(s.Statuses, key)
	found, source, err = repo.FindByID(ctx, key)
	if err != nil {
		t.Fatalf("FindByID redis failed: %v", err)
	}
	if source != models.SourceRedis {
		t.Fatalf("expected redis source, got %v", source)
	}

	if _, _, err := repo.FindByID(ctx, "missing"); err == nil {
		t.Fatalf("expected error for unknown status")
	}
}


