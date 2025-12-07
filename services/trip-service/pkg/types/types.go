package pkg

import (
	pb "ride-sharing/shared/proto/trip"
)

type OsrmRespBody struct {
	Routes []struct {
		Weight float64 `json:"weight"`

		Distance float64 `json:"distance"`
		Geometry struct {
			Coordinates [][]float64 `json:"coordinates"`
		} `json:"geometry"`
	} `json:"routes"`
}

func (orsm *OsrmRespBody) ToProto() *pb.Route {

	route := orsm.Routes[0]
	geometry := route.Geometry.Coordinates

	coordinates := make([]*pb.Coordinate, len(route.Geometry.Coordinates))

	for i, coord := range geometry {
		coordinates[i] = &pb.Coordinate{
			Longitude: coord[0],
			Latitude:  coord[1],
		}
	}
	return &pb.Route{
		Geometry: []*pb.Geometry{
			{
				Coordinates: coordinates,
			},
		},
		Distance: route.Distance,
		Duration: route.Weight,
	}

}
