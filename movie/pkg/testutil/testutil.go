package testutil

import (
	"moviehub.com/gen"
	"moviehub.com/movie/internal/controller/movie"
	metadatagateway "moviehub.com/movie/internal/gateway/metadata/grpc"
	ratinggateway "moviehub.com/movie/internal/gateway/rating/grpc"
	grpchandler "moviehub.com/movie/internal/handler/grpc"
	"moviehub.com/pkg/discovery"
)

// NewTestMovieGRPCServer creates a new movie gRPC server to be used in tests.
func NewTestMovieGRPCServer(registry discovery.Registry) gen.MovieServiceServer {
	metadataGateway := metadatagateway.New(registry)
	ratingGateway := ratinggateway.New(registry)
	ctrl := movie.New(ratingGateway, metadataGateway)
	return grpchandler.New(ctrl)
}
