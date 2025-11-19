package status

import (
	"actions-service/internal/models"
	"actions-service/internal/state"
	"actions-service/internal/ws"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/alicebob/miniredis/v2"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type stubWorkcenterPort struct {
	workcenters map[string]*models.WorkcenterDTO
}

func (p *stubWorkcenterPort) GetWorkcenterDTO(_ context.Context, id string) (*models.WorkcenterDTO, error) {
	if wc, ok := p.workcenters[id]; ok {
		clone := *wc
		return &clone, nil
	}
	return nil, fmt.Errorf("workcenter %s not found", id)
}

type fakeStatusHTTPClient struct {
	getResponses  map[string][]*http.Response
	getCalls      map[string]int
	postResponse  *http.Response
	postErr       error
}

func newFakeStatusHTTPClient() *fakeStatusHTTPClient {
	return &fakeStatusHTTPClient{
		getResponses: make(map[string][]*http.Response),
		getCalls:     make(map[string]int),
	}
}

func (f *fakeStatusHTTPClient) addGetResponse(path string, resp *http.Response) {
	f.getResponses[path] = append(f.getResponses[path], resp)
}

func (f *fakeStatusHTTPClient) DoGetRequest(_ context.Context, path string) (*http.Response, error) {
	list, ok := f.getResponses[path]
	if !ok {
		return nil, fmt.Errorf("no response configured for %s", path)
	}
	idx := f.getCalls[path]
	if idx >= len(list) {
		return nil, fmt.Errorf("no more responses for %s", path)
	}
	f.getCalls[path] = idx + 1
	return list[idx], nil
}

func (f *fakeStatusHTTPClient) DoPostRequest(_ context.Context, _ string, _ interface{}) (*http.Response, error) {
	return f.postResponse, f.postErr
}

func httpResponse(status int, body interface{}) *http.Response {
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

func TestBuildDTOStoresStatuses(t *testing.T) {
	ctx := context.Background()
	mr, err := miniredis.Run()
	if err != nil {
		t.Fatalf("failed to start miniredis: %v", err)
	}
	defer mr.Close()
	redisClient := redis.NewClient(&redis.Options{Addr: mr.Addr()})

	s := state.New()
	repo := NewStatusRepository(s, redisClient)
	client := newFakeStatusHTTPClient()

	statusID := uuid.New()
	workcenterID := uuid.New()
	client.addGetResponse("/api/MachineStatus", httpResponse(http.StatusOK, []models.StatusResponse{{
		StatusId:     statusID,
		Description:  "Running",
		Color:        "green",
		Stopped:      false,
		OperatorsAllowed: true,
		Closed:       false,
	}}))
	client.addGetResponse("/api/WorkcenterCost", httpResponse(http.StatusOK, []models.StatusCostResponse{{
		StatusId:     statusID,
		WorkcenterId: workcenterID,
		Cost:         10,
	}}))

	service := NewStatusService(client, *repo, nil, nil)
	if err := service.BuildDTO(ctx); err != nil {
		t.Fatalf("BuildDTO failed: %v", err)
	}

	key := fmt.Sprintf("%s:%s", workcenterID, statusID)
	stored, _, err := repo.FindByID(ctx, key)
	if err != nil {
		t.Fatalf("expected stored status: %v", err)
	}
	if stored.Description != "Running" || stored.Cost != 10 {
		t.Fatalf("unexpected status stored: %+v", stored)
	}
}

func TestStatusInUpdatesWorkcenter(t *testing.T) {
	ctx := context.Background()
	mr, err := miniredis.Run()
	if err != nil {
		t.Fatalf("failed to start miniredis: %v", err)
	}
	defer mr.Close()
	redisClient := redis.NewClient(&redis.Options{Addr: mr.Addr()})

	s := state.New()
	repo := NewStatusRepository(s, redisClient)
	statusID := uuid.New()
	workcenterID := uuid.New()
	key := fmt.Sprintf("%s:%s", workcenterID, statusID)
	statusDTO := models.StatusDTO{
		StatusId:     statusID,
		WorkcenterId: workcenterID,
		Description:  "Running",
		OperatorsAllowed: true,
	}
	if err := repo.Set(ctx, key, statusDTO); err != nil {
		t.Fatalf("failed to seed status: %v", err)
	}

	workcenterDTO := &models.WorkcenterDTO{WorkcenterID: workcenterID}
	port := &stubWorkcenterPort{workcenters: map[string]*models.WorkcenterDTO{workcenterID.String(): workcenterDTO}}

	client := newFakeStatusHTTPClient()
	client.postResponse = httpResponse(http.StatusOK, "ok")

	hub := ws.NewHub()

	service := NewStatusService(client, *repo, port, hub)
	if err := service.StatusIn(ctx, workcenterID.String(), statusID.String()); err != nil {
		t.Fatalf("StatusIn failed: %v", err)
	}

	stored := s.Workcenters[workcenterID.String()]
	if stored.StatusName != "Running" {
		t.Fatalf("expected workcenter status to be updated, got %+v", stored)
	}
	if stored.StatusStartTime.IsZero() {
		t.Fatalf("expected status start time to be set")
	}
}


