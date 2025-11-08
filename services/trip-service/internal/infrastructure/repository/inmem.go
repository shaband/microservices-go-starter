package repository

import (
	"context"
	"ride-sharing/services/trip-service/internal/domain"
	// "ride-sharing/shared/types"
)

type inmemRepository struct {
	trips     map[string]*domain.TripModel
	rideFares map[string]*domain.RideFareModel
	// routes    map[string]*types.Route
}

func NewInmemRepository() *inmemRepository {
	return &inmemRepository{
		trips:     make(map[string]*domain.TripModel),
		rideFares: make(map[string]*domain.RideFareModel),
	}
}

func (r *inmemRepository) CreateTrip(ctx context.Context, trip *domain.TripModel) (*domain.TripModel, error) {
	r.trips[trip.ID.Hex()] = trip
	return trip, nil
}

// func (r *inmemRepository) CreateRoute(ctx context.Context, route *types.Route) (*types.Route, error) {
// 	// r.routes[''] = trip
// 	log.Println("inmem repo create route is not impelemented Yet")
// 	return route, nil
// }
