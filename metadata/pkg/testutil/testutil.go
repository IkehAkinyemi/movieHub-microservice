package testutil

import (
	"moviehub.com/gen"
	"moviehub.com/metadata/internal/controller/metadata"
	grpchandler "moviehub.com/metadata/internal/handler/grpc"
	"moviehub.com/metadata/internal/repository/memory"
)

// NewTestMetadataGRPCServer creates a new metadata gRPC server to be used in tests.
func NewTestMetadataGRPCServer() gen.MetadataServiceServer {
	r := memory.New()
	ctrl := metadata.New(r)
	return grpchandler.New(ctrl)
}
