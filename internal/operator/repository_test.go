package operator

import (
	"actions-service/internal/models"
	"actions-service/internal/state"
	"context"
	"testing"

	"github.com/alicebob/miniredis/v2"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

func TestRepositorySetAndFind(t *testing.T) {
	ctx := context.Background()
	mr, err := miniredis.Run()
	if err != nil {
		t.Fatalf("failed to start miniredis: %v", err)
	}
	defer mr.Close()

	client := redis.NewClient(&redis.Options{Addr: mr.Addr()})
	s := state.New()
	repo := NewOperatorRepository(s, client)

	id := uuid.New().String()
	op := models.OperatorDTO{OperatorID: uuid.MustParse(id), OperatorName: "Alice"}

	if err := repo.Set(ctx, id, op); err != nil {
		t.Fatalf("Set failed: %v", err)
	}

	found, source, err := repo.FindByID(ctx, id)
	if err != nil {
		t.Fatalf("FindByID failed: %v", err)
	}
	if source != models.SourceMemory {
		t.Fatalf("expected memory source, got %v", source)
	}
	if found.OperatorName != "Alice" {
		t.Fatalf("unexpected operator: %+v", found)
	}

	// Clear memory to force redis hit
	delete(s.Operators, id)
	found, source, err = repo.FindByID(ctx, id)
	if err != nil {
		t.Fatalf("FindByID from redis failed: %v", err)
	}
	if source != models.SourceRedis {
		t.Fatalf("expected redis source, got %v", source)
	}
	if found.OperatorName != "Alice" {
		t.Fatalf("unexpected operator from redis: %+v", found)
	}

	if _, _, err := repo.FindByID(ctx, uuid.NewString()); err == nil {
		t.Fatalf("expected error for unknown operator")
	}
}


