package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"gopkg.in/yaml.v2"
	"moviehub.com/gen"
	"moviehub.com/pkg/discovery"
	srvdiscovery "moviehub.com/pkg/discovery/memory"
	"moviehub.com/rating/internal/controller/rating"
	grpchandler "moviehub.com/rating/internal/handler/grpc"
	"moviehub.com/rating/internal/repository/mysql"
)

const serviceName = "rating"

func main() {
	log.Printf("Starting the rating service")

	f, err := os.Open("base.yaml")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	var cfg config
	if err := yaml.NewDecoder(f).Decode(&cfg); err != nil {
		panic(err)
	}

	registry := srvdiscovery.NewRegistry()

	ctx := context.Background()
	instanceID := discovery.GenerateInstanceID(serviceName)

	if err := registry.Register(ctx, instanceID, serviceName, fmt.Sprintf("localhost:%d", cfg.API.Port)); err != nil {
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

	repo, err := mysql.New()
	if err != nil {
		panic(err)
	}

	ctrl := rating.New(repo, nil)
	h := grpchandler.New(ctrl)
	listener, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", cfg.API.Port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	srv := grpc.NewServer()
	reflection.Register(srv)
	gen.RegisterRatingServiceServer(srv, h)
	if err := srv.Serve(listener); err != nil {
		panic(err)
	}
}
