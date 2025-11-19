package shift

import (
	"actions-service/internal/models"
	"actions-service/internal/state"
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestRepositorySetAndFindCurrent(t *testing.T) {
	ctx := context.Background()
	s := state.New()
	repo := NewShiftRepository(s)

	shiftID := uuid.New()
	detailID := uuid.New()
	shift := models.ShiftDTO{
		ID:   shiftID,
		Name: "Morning",
		ShiftDetails: []models.ShiftDetailDTO{
			{
				ID: detailID,
				StartTime: models.CustomTime{Time: time.Date(0, 1, 1, 6, 0, 0, 0, time.UTC)},
				EndTime:   models.CustomTime{Time: time.Date(0, 1, 1, 14, 0, 0, 0, time.UTC)},
				IsProductiveTime: true,
			},
		},
	}

	repo.Set(ctx, shiftID.String(), shift)

	now := time.Date(2025, 1, 1, 7, 0, 0, 0, time.UTC)
	current, err := repo.FindCurrent(ctx, now, shiftID.String())
	if err != nil {
		t.Fatalf("expected to find current shift detail, got %v", err)
	}
	if current.ID != detailID {
		t.Fatalf("expected detail %s, got %s", detailID, current.ID)
	}

	noShiftTime := time.Date(2025, 1, 1, 23, 0, 0, 0, time.UTC)
	_, err = repo.FindCurrent(ctx, noShiftTime, shiftID.String())
	if err == nil {
		t.Fatalf("expected error for time outside shift window")
	}

	_, err = repo.FindCurrent(ctx, now, uuid.NewString())
	if err == nil {
		t.Fatalf("expected error for unknown shift")
	}
}


