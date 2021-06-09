package main

import (
	"context"
	"errors"
	"flag"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var portsClient pb.PortsClient

var filePath string

func main() {
	var domainServiceHost string
	flag.StringVar(&domainServiceHost, "domainservice", "localhost:5001", "hostname and port of the PortDomainService of the form 'localhost:5001'")
	flag.StringVar(&filePath, "file", "", "location of the file to parse when the 'start' endpoint is called")
	flag.Parse()

	// set up GRPC client
	opts := []grpc.DialOption{grpc.WithInsecure(), grpc.WithBlock()}
	grpcConnection, err := grpc.Dial(domainServiceHost, opts...)
	if err != nil {
		log.Fatal("failed to dial domain service:", err)
	}

	portsClient = pb.NewPortsClient(grpcConnection)

	// set up HTTP server
	router := gin.Default()
	router.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Client API")
	})

	srv := &http.Server{
		Addr:    ":5000",
		Handler: router,
	}

	// start the server in a goroutine so that interrupt signals
	// can be used to gracefully shut the service down
	go func() {
		if err := srv.ListenAndServe(); err != nil && errors.Is(err, http.ErrServerClosed) {
			log.Printf("listen: %s\n", err)
		}
	}()

	// Wait for interrupt signal to initiate shutdown
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// Give the server 5 seconds to finish handling any current requests
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exiting")
}
