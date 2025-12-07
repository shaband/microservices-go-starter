package main

import (
	pb "ride-sharing/shared/proto/trip"
	"ride-sharing/shared/types"
)

type previewTripRequest struct {
	UserID      string           `json:"userID"`
	Pickup      types.Coordinate `json:"pickup"`
	Destination types.Coordinate `json:"destination"`
}

func (preview previewTripRequest) ToProto() *pb.PreviewTripRequest {
	return &pb.PreviewTripRequest{
		UserID: preview.UserID,
		StartLocation: &pb.Coordinate{
			Latitude:  preview.Pickup.Latitude,
			Longitude: preview.Pickup.Longitude,
		},
		EndLocation: &pb.Coordinate{
			Latitude:  preview.Destination.Latitude,
			Longitude: preview.Destination.Longitude,
		},
	}
}
