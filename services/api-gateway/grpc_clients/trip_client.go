package grpc_clients

import (
	"os"
	pb "ride-sharing/shared/proto/trip"

	"google.golang.org/grpc"
)

type TripServiceClient struct {
	Client pb.TripServiceClient
	conn   *grpc.ClientConn
}

func NewTripServiceClient() (*TripServiceClient, error) {

	tripServiceUrl := os.Getenv("TRIP_SERVICE_URL")
	if tripServiceUrl == "" {
		tripServiceUrl = "trip-service:9093"
	}
	conn, err := grpc.NewClient(tripServiceUrl)

	if err != nil {
		return nil, err
	}
	client := pb.NewTripServiceClient(conn)
	return &TripServiceClient{
		Client: client,
		conn:   conn,
	}, nil
}

func (t *TripServiceClient) Close() error {
	if t.conn == nil {
		return nil
	}
	return t.conn.Close()
}
