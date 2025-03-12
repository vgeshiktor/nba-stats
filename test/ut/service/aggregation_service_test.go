// test/ut/service/aggregation_service_test.go
package service_test

import (
	"testing"

	"github.com/vgeshiktor/nba-stats/internal/service"
	"github.com/vgeshiktor/nba-stats/test/ut/mocks"
)

func TestGetPlayerAggregate_Success(t *testing.T) {
	// Use FakePlayerStatsRepo from our mocks package.
	fakeStatsRepo := &mocks.FakePlayerStatsRepo{}
	aggService := service.NewAggregationService(fakeStatsRepo)

	agg, err := aggService.GetPlayerAggregate("valid")
	if err != nil {
		t.Errorf("expected no error, got: %v", err)
	}
	if agg == nil {
		t.Errorf("expected non-nil aggregate, got nil")
	}
	if agg.PlayerID != "valid" {
		t.Errorf("expected playerID 'valid', got %s", agg.PlayerID)
	}
}

func TestGetPlayerAggregate_Invalid(t *testing.T) {
	fakeStatsRepo := &mocks.FakePlayerStatsRepo{}
	aggService := service.NewAggregationService(fakeStatsRepo)

	// Use an ID that is not recognized by the fake repository.
	_, err := aggService.GetPlayerAggregate("invalid")
	if err == nil {
		t.Errorf("expected error for invalid playerID, got nil")
	}
}

func TestGetTeamAggregate_Success(t *testing.T) {
	fakeStatsRepo := &mocks.FakePlayerStatsRepo{}
	aggService := service.NewAggregationService(fakeStatsRepo)

	agg, err := aggService.GetTeamAggregate("team1")
	if err != nil {
		t.Errorf("expected no error, got: %v", err)
	}
	if agg == nil {
		t.Errorf("expected non-nil aggregate, got nil")
	}
	if agg.TeamID != "team1" {
		t.Errorf("expected teamID 'team1', got %s", agg.TeamID)
	}
}

func TestGetTeamAggregate_Invalid(t *testing.T) {
	fakeStatsRepo := &mocks.FakePlayerStatsRepo{}
	aggService := service.NewAggregationService(fakeStatsRepo)

	// Use an ID that is not recognized by the fake repository.
	_, err := aggService.GetTeamAggregate("invalid")
	if err == nil {
		t.Errorf("expected error for invalid teamID, got nil")
	}
}
