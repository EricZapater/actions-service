package shift

import (
	"actions-service/internal/models"
	"actions-service/internal/state"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
)

type fakeHTTPClient struct {
	responses map[string][]*http.Response
	counts    map[string]int
}

func newFakeHTTPClient() *fakeHTTPClient {
	return &fakeHTTPClient{
		responses: make(map[string][]*http.Response),
		counts:    make(map[string]int),
	}
}

func (f *fakeHTTPClient) addResponse(path string, resp *http.Response) {
	f.responses[path] = append(f.responses[path], resp)
}

func (f *fakeHTTPClient) DoGetRequest(_ context.Context, path string) (*http.Response, error) {
	list, ok := f.responses[path]
	if !ok {
		return nil, fmt.Errorf("no response configured for %s", path)
	}
	idx := f.counts[path]
	if idx >= len(list) {
		return nil, fmt.Errorf("no more responses for %s", path)
	}
	f.counts[path] = idx + 1
	return list[idx], nil
}

func (f *fakeHTTPClient) DoPostRequest(_ context.Context, _ string, _ interface{}) (*http.Response, error) {
	return nil, fmt.Errorf("not implemented")
}

func makeResponse(status int, body interface{}) *http.Response {
	var reader io.Reader
	switch v := body.(type) {
	case string:
		reader = strings.NewReader(v)
	default:
		b, _ := json.Marshal(v)
		reader = strings.NewReader(string(b))
	}
	return &http.Response{
		StatusCode: status,
		Body:       io.NopCloser(reader),
	}
}

func TestBuildDTOStoresShifts(t *testing.T) {
	ctx := context.Background()
	s := state.New()
	repo := NewShiftRepository(s)
	client := newFakeHTTPClient()

	shiftID := uuid.New()
	shiftResponse := []models.ShiftDTO{{ID: shiftID, Name: "Morning"}}
	detailID := uuid.New()
	start := time.Date(0, 1, 1, 6, 0, 0, 0, time.UTC)
	end := time.Date(0, 1, 1, 14, 0, 0, 0, time.UTC)
	detailResponse := []models.ShiftDetailDTO{{
		ID: detailID,
		StartTime: models.CustomTime{Time: start},
		EndTime:   models.CustomTime{Time: end},
	}}

	client.addResponse("/api/Shift", makeResponse(http.StatusOK, shiftResponse))
	pathDetail := fmt.Sprintf("/api/Shift/Detail/%v", shiftID)
	client.addResponse(pathDetail, makeResponse(http.StatusOK, detailResponse))

	service := NewShiftService(client, *repo)
	if err := service.BuildDTO(ctx); err != nil {
		t.Fatalf("BuildDTO failed: %v", err)
	}

	checkTime := time.Date(2025, 1, 1, 7, 0, 0, 0, time.UTC)
	stored, err := repo.FindCurrent(ctx, checkTime, shiftID.String())
	if err != nil {
		t.Fatalf("expected stored shift detail: %v", err)
	}
	if stored.ID != detailID {
		t.Fatalf("expected detail %s, got %s", detailID, stored.ID)
	}
}


