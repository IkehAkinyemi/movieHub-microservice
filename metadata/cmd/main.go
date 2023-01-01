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
	"moviehub.com/metadata/internal/controller/metadata"
	grpchandler "moviehub.com/metadata/internal/handler/grpc"
	"moviehub.com/metadata/internal/repository/memory"
	"moviehub.com/pkg/discovery"
	srvdiscovery "moviehub.com/pkg/discovery/memory"
)

const serviceName = "metadata"

func main() {
	log.Println("Starting the movie metadata service")

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
				log.Println("Failed to report healthy state: ", err.Error())
			}
			time.Sleep(1 * time.Second)
		}
	}()
	defer registry.Deregister(ctx, instanceID, serviceName)

	repo := memory.New()
	ctrl := metadata.New(repo)
	h := grpchandler.New(ctrl)

	listener, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", cfg.API.Port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	srv := grpc.NewServer()
	reflection.Register(srv)
	gen.RegisterMetadataServiceServer(srv, h)

	if err := srv.Serve(listener); err != nil {
		panic(err)
	}
}
