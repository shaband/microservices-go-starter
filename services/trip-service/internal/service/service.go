package service

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"ride-sharing/services/trip-service/internal/domain"
	"ride-sharing/shared/types"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type service struct {
	repo domain.TripRepository
}

func NewService(repo domain.TripRepository) *service {
	return &service{
		repo: repo,
	}
}

func (s *service) CreateTrip(ctx context.Context, fare *domain.RideFareModel) (*domain.TripModel, error) {
	t := &domain.TripModel{
		ID:       primitive.NewObjectID(),
		UserID:   fare.UserID,
		Status:   "pending",
		RideFare: fare,
	}

	return s.repo.CreateTrip(ctx, t)
}
func (s *service) CreateRoute(ctx context.Context) (*types.Route, error) {
	t := &types.Route{
		Distance: 100,
		Duration: 13,
		Geometry: []*types.Geometry{},
	}

	return t, nil
}
func (s *service) GetRoute(ctx context.Context, pickup types.Coordinate, destination types.Coordinate) (*types.OsrmRespBody, error) {
	url := fmt.Sprintf("http://router.project-osrm.org/route/v1/driving/%f,%f;%f,%f?overview=full&geometries=geojson",
		pickup.Longitude,
		pickup.Latitude,
		destination.Longitude,
		destination.Latitude,
	)
	res, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to get route from OSRM: %w", err)
	}
	data, err := io.ReadAll(res.Body)
	defer res.Body.Close()
	if err != nil {
		return nil, fmt.Errorf("failed to read OSRM response body: %w", err)
	}
	var osrmResp types.OsrmRespBody
	json.Unmarshal(data, &osrmResp)
	return &osrmResp, nil
}
