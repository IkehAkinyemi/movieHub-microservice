package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	"moviehub.com/pkg/discovery"
	"moviehub.com/pkg/discovery/consul"
	"moviehub.com/rating/internal/controller/rating"
	httphandler "moviehub.com/rating/internal/handler/http"
	"moviehub.com/rating/internal/repository/memory"
)

const serviceName = "rating"

func main() {
	var port int
	flag.IntVar(&port, "port", 8082, "API Handler port")
	flag.Parse()

	log.Printf("Starting the rating service on port %d\n", port)
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

	repo := memory.New()
	svc := rating.New(repo)
	h := httphandler.New(svc)
	http.Handle("/rating", http.HandlerFunc(h.Handle))
	if err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil); err != nil {
		panic(err)
	}
}
