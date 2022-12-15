package main

import (
	"log"
	"net/http"

	"moviehub.com/movie/internal/controller/movie"
	metadatagateway "moviehub.com/movie/internal/gateway/metadata/http"
	ratinggateway "moviehub.com/movie/internal/gateway/rating/http"
	httphandler "moviehub.com/movie/internal/handler/http"
)

func main() {
	log.Println("Starting this movie service")

	metadataGateway := metadatagateway.New("localhost:8081")
	ratingGateway := ratinggateway.New("localhost:8082")
	ctrl := movie.New(ratingGateway, metadataGateway)

	h := httphandler.New(ctrl)
	http.HandleFunc("/movie", http.HandlerFunc(h.GetMovieDetails))
	if err := http.ListenAndServe(":8083", nil); err != nil {
		panic(err)
	}
}
