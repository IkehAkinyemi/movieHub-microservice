package testutil

import (
	"moviehub.com/gen"
	"moviehub.com/rating/internal/controller/rating"
	grpchandler "moviehub.com/rating/internal/handler/grpc"
	"moviehub.com/rating/internal/repository/memory"
)

// NewTestRatingGRPCServer creates a new rating gRPC server to be used in tests.
func NewTestRatingGRPCServer() gen.RatingServiceServer {
	r := memory.New()
	ctrl := rating.New(r, nil)
	return grpchandler.New(ctrl)
}
