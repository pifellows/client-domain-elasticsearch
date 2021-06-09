package main

import (
	pb "client-domain-elasticsearch/PortsCommunication/PortsCommunication"
	"errors"
	"flag"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

var esClient ElasticsearchCommunication

func main() {
	var elasticsearchHost string
	var indexname string
	flag.StringVar(&elasticsearchHost, "es", "http://localhost:9200", "hostname and port of the elasticsearch service of the form 'http://localhost:9200'")
	flag.StringVar(&indexname, "index", "ports", "name of the index to search and index to")
	flag.Parse()

	esClient = ElasticsearchCommunication{}
	err := esClient.Initalise(elasticsearchHost, indexname, 100)
	if err != nil {
		log.Fatalf("failed to connect to elasticsearch: %v", err)
	}

	// create PortService listener
	listener, err := net.Listen("tcp", fmt.Sprintf(":5001"))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	var opts []grpc.ServerOption

	grpcServer := grpc.NewServer(opts...)
	pb.RegisterPortsServer(grpcServer, &portsServer{})

	// start the server in a goroutine so that interrupt signals
	// can be used to gracefully shut the service down
	go func() {
		if err := grpcServer.Serve(listener); err != nil && errors.Is(err, grpc.ErrServerStopped) {
			log.Printf("listen: %s\n", err)
		}
	}()

	// Wait for interrupt signal to initiate shutdown
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		grpcServer.GracefulStop()
		wg.Done()
	}()
	log.Println("Complete...")
}
