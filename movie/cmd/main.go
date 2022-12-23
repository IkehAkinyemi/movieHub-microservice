package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"moviehub.com/gen"
	"moviehub.com/movie/internal/controller/movie"
	metadatagateway "moviehub.com/movie/internal/gateway/metadata/grpc"
	ratinggateway "moviehub.com/movie/internal/gateway/rating/grpc"
	grpchandler "moviehub.com/movie/internal/handler/grpc"
	"moviehub.com/pkg/discovery"
	"moviehub.com/pkg/discovery/consul"
)

const serviceName = "movie"

func main() {
	var port int

	flag.IntVar(&port, "port", 8083, "API Handler port")
	flag.Parse()

	log.Printf("Starting this movie service on port %d", port)

	registry, err := consul.NewRegistry("localhost:8500")
	if err != nil {
		panic(err)
	}

	ctx := context.Background()
	instanceID := discovery.GenerateInstanceID(serviceName)

	if err := registry.Register(ctx, instanceID, serviceName, fmt.Sprintf("localhost:%d", port)); err != nil {
		panic(err)
	}
	go func() {
		for {
			if err := registry.ReportHealthyState(instanceID, serviceName); err != nil {
				log.Println("Failed to report healthy state: " + err.Error())
			}

			time.Sleep(1 * time.Second)
		}
	}()
	defer registry.Deregister(ctx, instanceID, serviceName)

	metadataGateway := metadatagateway.New(registry)
	ratingGateway := ratinggateway.New(registry)
	ctrl := movie.New(ratingGateway, metadataGateway)

	h := grpchandler.New(ctrl)
	listener, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	srv := grpc.NewServer()
	reflection.Register(srv)
	gen.RegisterMovieServiceServer(srv, h)
	if err := srv.Serve(listener); err != nil {
		panic(err)
	}
}
