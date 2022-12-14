package main

import (
	"log"
	"net/http"

	"moviehub.com/metadata/internal/controller/metadata"
	httphandler "moviehub.com/metadata/internal/handler/http"
	"moviehub.com/metadata/internal/repository/memory"
)

func main() {
	log.Println("Starting the movie metadata service")
	repo := memory.New()
	ctrl := metadata.New(repo)
	h := httphandler.New(ctrl)
	http.HandleFunc("/metadata", http.HandlerFunc(h.GetMetadata))
	if err := http.ListenAndServe(":8081", nil); err != nil {
		panic(err)
	}
}
