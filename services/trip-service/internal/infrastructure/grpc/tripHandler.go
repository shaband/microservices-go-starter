package grpc

import (
	"context"
	"fmt"
	"ride-sharing/services/trip-service/internal/domain"
	pb "ride-sharing/shared/proto/trip"
	"ride-sharing/shared/types"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type gRPCHandler struct {
	pb.UnimplementedTripServiceServer

	service   domain.TripService
	publisher *events.TripEventPublisher
}


func NewtripHandler(s domain.TripService) *TripHandler {
	return &TripHandler{
		Service: s,
	}
}
func RegisterTripHandler(server *grpc.Server, svc domain.TripService) *TripHandler {
	handler := NewtripHandler(svc)
	pb.RegisterTripServiceServer(server, handler)
	return handler
}

func (s *TripHandler) PreviewTrip(ctx context.Context, req *pb.PreviewTripRequest) (*pb.PreviewTripResponse, error) {

	route, err := s.Service.GetRoute(ctx, types.Coordinate{
		Latitude:  req.StartLocation.Latitude,
		Longitude: req.StartLocation.Longitude,
	}, types.Coordinate{
		Latitude:  req.EndLocation.Latitude,
		Longitude: req.EndLocation.Longitude,
	})
	if err != nil {
		status.Errorf(codes.Internal, "failed to get route: %v", err)
		fmt.Println("Error fetching route:", err)
	}
	return &pb.PreviewTripResponse{
		Route:     route.ToProto(),
		RideFares: []*pb.RideFare{},
	}, nil
}
