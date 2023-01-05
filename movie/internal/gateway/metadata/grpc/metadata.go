package grpc

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"moviehub.com/gen"
	"moviehub.com/internal/grpcutil"
	"moviehub.com/metadata/pkg/model"
	"moviehub.com/pkg/discovery"
)

// Gateway defines a movie metadata gRPC gateway.
type Gateway struct {
	registry discovery.Registry
}

// New creates a new gRPC gateway for the movie metadata
// service.
func New(registry discovery.Registry) *Gateway {
	return &Gateway{registry}
}

// Get returns movie metadata by a movie id.
func (g *Gateway) Get(ctx context.Context, id string) (*model.Metadata, error) {
	conn, err := grpcutil.ServiceConnection(ctx, "metadata", g.registry)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	client := gen.NewMetadataServiceClient(conn)
	var resp *gen.GetMetadataResponse
	const maxRetries = 5
	for i := 0; i < maxRetries; i++ {
		resp, err = client.GetMetadata(ctx, &gen.GetMetadataRequest{MovieId: id})
		if err != nil {
			if shouldRetry(err) {
				continue
			}
			return nil, err
		}
		return model.MetadataFromProto(resp.Metadata), nil
	}

	return model.MetadataFromProto(resp.Metadata), nil
}

func shouldRetry(err error) bool {
	e, ok := status.FromError(err)
	if !ok {
		return false
	}
	return e.Code() == codes.DeadlineExceeded || e.Code() == codes.ResourceExhausted || e.Code() == codes.Unavailable
}
