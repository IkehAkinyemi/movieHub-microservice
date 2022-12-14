package main

import (
	"log"
	"net/http"

	"moviehub.com/rating/internal/controller/rating"
	httphandler "moviehub.com/rating/internal/handler/http"
	"moviehub.com/rating/internal/repository/memory"
)

func main() {
	log.Printf("Starting the rating service")
	repo := memory.New()
	ctrl := rating.New(repo)
	h := httphandler.New(ctrl)
	http.HandleFunc("/rating", http.HandlerFunc(h.Handle))
	if err := http.ListenAndServe(":8082", nil); err != nil {
		panic(err)
	}
}
