package main

import (
	pb "client-domain-elasticsearch/PortsCommunication/PortsCommunication"
	"context"
	"io"
)

type portsServer struct {
	pb.UnimplementedPortsServer
}

func (s *portsServer) GetPort(ctx context.Context, portID *pb.PortID) (*pb.Port, error) {
	return esClient.GetPort(ctx, portID.GetId())
}

func (s *portsServer) StreamPorts(stream pb.Ports_StreamPortsServer) error {
	// if the stream closes or we return an error, we want to cancel
	// any requests to Elasticsearch that are in flight
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	for {
		port, err := stream.Recv()
		if err == io.EOF {
			// if we have reached the end of the stream
			// we need to commit any remaining documents
			err = esClient.IndexBatch(ctx)
			if err != nil {
				return err
			}
			return stream.SendAndClose(&pb.Summary{
				Total:  0,
				Failed: 0,
			})
		}
		if err != nil {
			return err
		}
		err = esClient.AddDocument(ctx, port)
		if err != nil {
			return err
		}
	}
}
