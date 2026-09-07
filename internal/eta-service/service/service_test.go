package service

import (
	"errors"
	"testing"
	"time"

	"github.com/nutanalabs/eta-service/internal/config"
	"github.com/nutanalabs/eta-service/internal/constants"
	routingengine "github.com/nutanalabs/eta-service/internal/serviceclients/routing-engine"
	"github.com/nutanalabs/eta-service/internal/types"
)

type mockRepository struct {
	restaurantEstimates  []types.EtaRestaurantEstimates
	restaurantErr        error
	sublocalityEstimates []types.EtaSublocalityEstimates
	sublocalityErr       error
	publishCalled        bool
}

func (m *mockRepository) FetchRestaurantEstimates(_ string, _ constants.Day, _ constants.MealType) (*types.EtaRestaurantEstimates, error) {
	return nil, nil
}

func (m *mockRepository) FetchSublocalityEstimates(_ string, _ constants.Day, _ constants.MealType) (*types.EtaSublocalityEstimates, error) {
	return nil, nil
}

func (m *mockRepository) FetchRestaurantEstimatesByIDs(_ []string, _ constants.Day, _ constants.MealType, _ int) ([]types.EtaRestaurantEstimates, error) {
	return m.restaurantEstimates, m.restaurantErr
}

func (m *mockRepository) FetchSublocalityEstimatesByIDs(_ []string, _ constants.Day, _ constants.MealType, _ int) ([]types.EtaSublocalityEstimates, error) {
	return m.sublocalityEstimates, m.sublocalityErr
}

func (m *mockRepository) RestaurantDayOverlapExists(_ string, _ []constants.Day) (bool, error) {
	return false, nil
}

func (m *mockRepository) SublocalityDayOverlapExists(_ string, _ []constants.Day) (bool, error) {
	return false, nil
}

func (m *mockRepository) InsertRestaurantEstimates(_ *types.InsertRestaurantEstimateRequest) error {
	return nil
}

func (m *mockRepository) InsertSublocalityEstimates(_ *types.InsertSublocalityEstimateRequest) error {
	return nil
}

func (m *mockRepository) UpdateRestaurantEstimates(_ string, _ *types.UpdateRestaurantEstimateRequest) error {
	return nil
}

func (m *mockRepository) UpdateSublocalityEstimates(_ string, _ *types.UpdateSublocalityEstimateRequest) error {
	return nil
}

func (m *mockRepository) PublishFetchEtaEvent(_ constants.EventType, _ *types.FetchEtaRequest, _ []types.FetchEtaResponse, _ error) {
	m.publishCalled = true
}

type mockCommonUtils struct {
	day      constants.Day
	mealType constants.MealType
	distance float64
}

func (m *mockCommonUtils) GetDayFromTime(_ time.Time) constants.Day {
	return m.day
}

func (m *mockCommonUtils) GetMealTypeFromTime(_ time.Time) constants.MealType {
	return m.mealType
}

func (m *mockCommonUtils) ToJSON(_ interface{}) string {
	return ""
}

func (m *mockCommonUtils) GetHaversineDistance(_, _, _, _ float64) float64 {
	return m.distance
}

type mockRoutingClient struct {
	response *routingengine.DistanceMatrixResponse
	err      error
}

func (m *mockRoutingClient) GetDistanceMatrix(_ *routingengine.DistanceMatrixRequest) (*routingengine.DistanceMatrixResponse, error) {
	return m.response, m.err
}

func (m *mockRoutingClient) GetDistanceMatrixWithQoS(_ *routingengine.DistanceMatrixRequest, _ constants.QosLevel) (*routingengine.DistanceMatrixResponse, error) {
	return m.response, m.err
}

func testConfig() *config.Config {
	return &config.Config{
		Mongo: config.MongoConfig{QueryBatchSize: 100},
		EtaDefaultEstimates: config.EtaDefaultEstimates{
			RestaurantAcceptanceTime: 120,
			KitchenPreparationTime:   900,
			PickupTime:               180,
			CaptainAssignmentTime:    300,
			FirstMileTime:            600,
			DelayDispatchTime:        60,
		},
	}
}

func newTestService(repo *mockRepository, utils *mockCommonUtils, routing *mockRoutingClient) *serviceImpl {
	return &serviceImpl{
		repository:    repo,
		commonUtils:   utils,
		routingClient: routing,
		config:        testConfig(),
	}
}

func populatedRestaurantMealEstimate() types.RestaurantMealEstimate {
	return types.RestaurantMealEstimate{
		Rat:           types.TimeSample{Seconds: 100, SampleCount: 1},
		Kpt:           types.TimeSample{Seconds: 800, SampleCount: 1},
		Pickup:        types.TimeSample{Seconds: 150, SampleCount: 1},
		DelayDispatch: types.TimeSample{Seconds: 50, SampleCount: 1},
	}
}

func populatedSublocalityMealEstimate() types.SublocalityMealEstimate {
	return types.SublocalityMealEstimate{
		Cat: types.TimeSample{Seconds: 200, SampleCount: 1},
		Fm:  types.TimeSample{Seconds: 400, SampleCount: 1},
	}
}

func TestExtractRestaurantLocationsAndIDs(t *testing.T) {
	entities := []types.FetchEtaRequestEntity{
		{RestaurantID: "r1", RestaurantLocation: types.Location{Lat: 12.1, Lng: 77.1}},
		{RestaurantID: "r2", RestaurantLocation: types.Location{Lat: 12.2, Lng: 77.2}},
	}

	sources, ids := extractRestaurantLocationsAndIDs(entities)

	if len(sources) != 2 || len(ids) != 2 {
		t.Fatalf("expected 2 sources and ids, got %d and %d", len(sources), len(ids))
	}
	if ids[0] != "r1" || ids[1] != "r2" {
		t.Fatalf("unexpected ids: %v", ids)
	}
	if sources[0].Lat != 12.1 || sources[1].Lng != 77.2 {
		t.Fatalf("unexpected sources: %v", sources)
	}
}

func TestRestaurantSublocalityID(t *testing.T) {
	byRestaurant := map[string]types.EtaRestaurantEstimates{
		"r1": {SublocalityId: "sub1"},
	}

	if got := restaurantSublocalityID("r1", byRestaurant, false); got != "sub1" {
		t.Fatalf("expected sub1, got %q", got)
	}
	if got := restaurantSublocalityID("missing", byRestaurant, false); got != "" {
		t.Fatalf("expected empty for missing restaurant, got %q", got)
	}
	if got := restaurantSublocalityID("r1", byRestaurant, true); got != "" {
		t.Fatalf("expected empty when restaurant fetch failed, got %q", got)
	}
}

func TestUniqueSublocalityIDs(t *testing.T) {
	estimates := []types.EtaRestaurantEstimates{
		{SublocalityId: "sub1"},
		{SublocalityId: "sub2"},
		{SublocalityId: "sub1"},
	}

	ids := uniqueSublocalityIDs(estimates)
	if len(ids) != 2 {
		t.Fatalf("expected 2 unique ids, got %d: %v", len(ids), ids)
	}
}

func TestResolveEtaSource(t *testing.T) {
	if got := resolveEtaSource(false); got != constants.EtaSourceHistoric {
		t.Fatalf("expected historic, got %q", got)
	}
	if got := resolveEtaSource(true); got != constants.EtaSourceFallback {
		t.Fatalf("expected fallback, got %q", got)
	}
}

func TestResolveRestaurantMealEstimate(t *testing.T) {
	svc := newTestService(&mockRepository{}, &mockCommonUtils{}, &mockRoutingClient{})
	meal := populatedRestaurantMealEstimate()
	byRestaurant := map[string]types.EtaRestaurantEstimates{
		"r1": {RestaurantId: "r1", Lunch: meal},
	}

	got, usedFallback := svc.resolveRestaurantMealEstimate("r1", constants.Lunch, byRestaurant, false)
	if usedFallback || got.Rat.Seconds != 100 {
		t.Fatalf("expected historic restaurant meal, got fallback=%v meal=%+v", usedFallback, got)
	}

	got, usedFallback = svc.resolveRestaurantMealEstimate("missing", constants.Lunch, byRestaurant, false)
	if !usedFallback || got.Rat.Seconds != 120 {
		t.Fatalf("expected default fallback for missing restaurant, got fallback=%v rat=%f", usedFallback, got.Rat.Seconds)
	}

	got, usedFallback = svc.resolveRestaurantMealEstimate("r1", constants.Lunch, byRestaurant, true)
	if !usedFallback || got.Kpt.Seconds != 900 {
		t.Fatalf("expected default fallback on mongo failure, got fallback=%v kpt=%f", usedFallback, got.Kpt.Seconds)
	}
}

func TestResolveSublocalityMealEstimate(t *testing.T) {
	svc := newTestService(&mockRepository{}, &mockCommonUtils{}, &mockRoutingClient{})
	meal := populatedSublocalityMealEstimate()
	bySublocality := map[string]types.EtaSublocalityEstimates{
		"sub1": {SublocalityId: "sub1", Lunch: meal},
	}

	got, usedFallback := svc.resolveSublocalityMealEstimate("sub1", constants.Lunch, bySublocality, false)
	if usedFallback || got.Cat.Seconds != 200 {
		t.Fatalf("expected historic sublocality meal, got fallback=%v meal=%+v", usedFallback, got)
	}

	got, usedFallback = svc.resolveSublocalityMealEstimate("", constants.Lunch, bySublocality, false)
	if !usedFallback || got.Fm.Seconds != 600 {
		t.Fatalf("expected default fallback for empty sublocality id, got fallback=%v fm=%f", usedFallback, got.Fm.Seconds)
	}

	got, usedFallback = svc.resolveSublocalityMealEstimate("sub1", constants.Lunch, bySublocality, true)
	if !usedFallback || got.Cat.Seconds != 300 {
		t.Fatalf("expected default fallback on mongo failure, got fallback=%v cat=%f", usedFallback, got.Cat.Seconds)
	}
}

func TestLoadRestaurantEstimates(t *testing.T) {
	repo := &mockRepository{
		restaurantEstimates: []types.EtaRestaurantEstimates{
			{RestaurantId: "r1", SublocalityId: "sub1"},
		},
	}
	svc := newTestService(repo, &mockCommonUtils{}, &mockRoutingClient{})

	byID, slice, failed := svc.loadRestaurantEstimates([]string{"r1"}, constants.Monday, constants.Lunch, 100)
	if failed {
		t.Fatal("expected success")
	}
	if len(byID) != 1 || len(slice) != 1 {
		t.Fatalf("unexpected results: map=%d slice=%d", len(byID), len(slice))
	}
	if byID["r1"].SublocalityId != "sub1" {
		t.Fatalf("unexpected estimate: %+v", byID["r1"])
	}

	repo.restaurantErr = errors.New("mongo down")
	byID, slice, failed = svc.loadRestaurantEstimates([]string{"r1"}, constants.Monday, constants.Lunch, 100)
	if !failed || byID != nil || slice != nil {
		t.Fatal("expected failure with nil map and slice")
	}
}

func TestLoadSublocalityEstimates(t *testing.T) {
	repo := &mockRepository{
		sublocalityEstimates: []types.EtaSublocalityEstimates{
			{SublocalityId: "sub1"},
		},
	}
	svc := newTestService(repo, &mockCommonUtils{}, &mockRoutingClient{})

	restaurantEstimates := []types.EtaRestaurantEstimates{{SublocalityId: "sub1"}}
	byID, failed := svc.loadSublocalityEstimates(restaurantEstimates, constants.Monday, constants.Lunch, 100, false)
	if failed || len(byID) != 1 {
		t.Fatalf("expected success with one sublocality, failed=%v map=%v", failed, byID)
	}

	_, failed = svc.loadSublocalityEstimates(restaurantEstimates, constants.Monday, constants.Lunch, 100, true)
	if !failed {
		t.Fatal("expected failure when restaurant fetch failed")
	}

	repo.sublocalityErr = errors.New("mongo down")
	_, failed = svc.loadSublocalityEstimates(restaurantEstimates, constants.Monday, constants.Lunch, 100, false)
	if !failed {
		t.Fatal("expected failure on sublocality mongo error")
	}
}

func TestCalculateLastMileDurationSeconds(t *testing.T) {
	svc := newTestService(&mockRepository{}, &mockCommonUtils{distance: 5.0}, &mockRoutingClient{})

	haversine := svc.calculateLastMileDurationSeconds(
		types.Location{Lat: 12.0, Lng: 77.0},
		types.Location{Lat: 12.1, Lng: 77.1},
		nil,
		true,
		0,
	)
	if haversine != 5.0/constants.HF_SPEED {
		t.Fatalf("expected haversine fallback, got %f", haversine)
	}

	matrix := &routingengine.DistanceMatrixResponse{
		Data: [][]routingengine.MatrixElement{
			{{Duration: routingengine.DurationInfo{Value: 420}}},
		},
	}
	routing := svc.calculateLastMileDurationSeconds(
		types.Location{Lat: 12.0, Lng: 77.0},
		types.Location{Lat: 12.1, Lng: 77.1},
		matrix,
		false,
		0,
	)
	if routing != 420 {
		t.Fatalf("expected routing duration 420, got %f", routing)
	}
}

func TestFetchEta_HistoricSource(t *testing.T) {
	meal := populatedRestaurantMealEstimate()
	subMeal := populatedSublocalityMealEstimate()

	repo := &mockRepository{
		restaurantEstimates: []types.EtaRestaurantEstimates{
			{
				RestaurantId:  "r1",
				SublocalityId: "sub1",
				Lunch:         meal,
			},
		},
		sublocalityEstimates: []types.EtaSublocalityEstimates{
			{SublocalityId: "sub1", Lunch: subMeal},
		},
	}
	routing := &mockRoutingClient{
		response: &routingengine.DistanceMatrixResponse{
			Data: [][]routingengine.MatrixElement{
				{{Duration: routingengine.DurationInfo{Value: 300}}},
			},
		},
	}
	svc := newTestService(repo, &mockCommonUtils{day: constants.Monday, mealType: constants.Lunch}, routing)

	resp, err := svc.FetchEta(&types.FetchEtaRequest{
		UserLocation: types.Location{Lat: 12.97, Lng: 77.59},
		Entities: []types.FetchEtaRequestEntity{
			{RestaurantID: "r1", RestaurantLocation: types.Location{Lat: 12.98, Lng: 77.60}},
		},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(resp) != 1 {
		t.Fatalf("expected 1 response, got %d", len(resp))
	}

	expectedETA := uint(1200)
	if resp[0].EtaInSeconds != expectedETA {
		t.Fatalf("expected eta %d, got %d", expectedETA, resp[0].EtaInSeconds)
	}
	if resp[0].Source != constants.EtaSourceHistoric {
		t.Fatalf("expected historic source, got %q", resp[0].Source)
	}
	if !repo.publishCalled {
		t.Fatal("expected analytics event to be published")
	}
}

func TestFetchEta_FallbackOnMongoAndRoutingFailure(t *testing.T) {
	repo := &mockRepository{restaurantErr: errors.New("mongo down")}
	routing := &mockRoutingClient{err: errors.New("routing down")}
	svc := newTestService(repo, &mockCommonUtils{day: constants.Monday, mealType: constants.Lunch, distance: 3.6}, routing)

	resp, err := svc.FetchEta(&types.FetchEtaRequest{
		UserLocation: types.Location{Lat: 12.97, Lng: 77.59},
		Entities: []types.FetchEtaRequestEntity{
			{RestaurantID: "r1", RestaurantLocation: types.Location{Lat: 12.98, Lng: 77.60}},
		},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	lastMile := 3.6 / constants.HF_SPEED
	expectedETA := uint(120 + 1140 + lastMile)
	if resp[0].EtaInSeconds != expectedETA {
		t.Fatalf("expected eta %d, got %d", expectedETA, resp[0].EtaInSeconds)
	}
	if resp[0].Source != constants.EtaSourceFallback {
		t.Fatalf("expected fallback source, got %q", resp[0].Source)
	}
}

func TestFetchEta_MultipleEntities(t *testing.T) {
	meal := populatedRestaurantMealEstimate()
	subMeal := populatedSublocalityMealEstimate()

	repo := &mockRepository{
		restaurantEstimates: []types.EtaRestaurantEstimates{
			{RestaurantId: "r1", SublocalityId: "sub1", Lunch: meal},
			{RestaurantId: "r2", SublocalityId: "sub1", Lunch: meal},
		},
		sublocalityEstimates: []types.EtaSublocalityEstimates{
			{SublocalityId: "sub1", Lunch: subMeal},
		},
	}
	routing := &mockRoutingClient{
		response: &routingengine.DistanceMatrixResponse{
			Data: [][]routingengine.MatrixElement{
				{{Duration: routingengine.DurationInfo{Value: 100}}},
				{{Duration: routingengine.DurationInfo{Value: 200}}},
			},
		},
	}
	svc := newTestService(repo, &mockCommonUtils{day: constants.Monday, mealType: constants.Lunch}, routing)

	resp, err := svc.FetchEta(&types.FetchEtaRequest{
		UserLocation: types.Location{Lat: 12.97, Lng: 77.59},
		Entities: []types.FetchEtaRequestEntity{
			{RestaurantID: "r1", RestaurantLocation: types.Location{Lat: 12.98, Lng: 77.60}},
			{RestaurantID: "r2", RestaurantLocation: types.Location{Lat: 12.99, Lng: 77.61}},
		},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(resp) != 2 {
		t.Fatalf("expected 2 responses, got %d", len(resp))
	}
	if resp[0].RestaurantID != "r1" || resp[1].RestaurantID != "r2" {
		t.Fatalf("unexpected restaurant ids: %+v", resp)
	}
	if resp[0].EtaInSeconds == resp[1].EtaInSeconds {
		t.Fatal("expected different ETAs due to different last-mile durations")
	}
}
