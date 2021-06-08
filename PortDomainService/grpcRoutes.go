package main

import (
	pb "client-domain-elasticsearch/PortsCommunication/PortsCommunication"
	"context"
)

type portsServer struct {
	pb.UnimplementedPortsServer
}

func (s *portsServer) GetPort(ctx context.Context, portID *pb.PortID) (*pb.Port, error) {
	return &pb.Port{
		Id:          portID.String(),
		Name:        "",
		City:        "",
		Country:     "",
		Alias:       []string{},
		Regions:     []string{},
		Coordinates: []float64{},
		Province:    "",
		Timezone:    "",
		Unlocs:      "",
		Code:        "",
	}, nil
}

func (s *portsServer) StreamPorts(stream pb.Ports_StreamPortsServer) error {
	return nil
}
