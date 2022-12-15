package http

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"moviehub.com/movie/internal/gateway"
	"moviehub.com/rating/pkg/model"
)

// Gateway defines an HTTP gateway for rating service.
type Gateway struct {
	addr string
}

// New instantiates a new HTTP gateway for a rating
// service
func New(addr string) *Gateway {
	return &Gateway{addr}
}

// GetAggregateRating returns the aggregated rating for a
// record or ErrNotFound if there's no rating for it.
func (g *Gateway) GetAggregateRating(ctx context.Context, recordID model.RecordID, recordType model.RecordType) (float64, error) {
	req, err := http.NewRequest(http.MethodGet, g.addr+"/rating", nil)
	if err != nil {
		return 0, err
	}
	req = req.WithContext(ctx)
	values := req.URL.Query()
	values.Add("id", string(recordID))
	values.Add("type", fmt.Sprintf("%v", recordType))
	req.URL.RawQuery = values.Encode()

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0, nil
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return 0, gateway.ErrNotFound
	} else if resp.StatusCode/100 != 2 {
		return 0, fmt.Errorf("non-2xx response: %v", resp)
	}

	var v float64
	if err := json.NewDecoder(resp.Body).Decode(&v); err != nil {
		return 0, err
	}

	return v, nil
}

// PutRating writes a rating.
func (g *Gateway) PutRating(ctx context.Context, recordID model.RecordID, recordType model.RecordType, rating *model.Rating) error {
	req, err := http.NewRequest(http.MethodPut, g.addr+"/rating", nil)
	if err != nil {
		return err
	}
	req = req.WithContext(ctx)
	values := req.URL.Query()
	values.Add("id", string(recordID))
	values.Add("type", fmt.Sprintf("%v", recordType))
	values.Add("userId", string(rating.UserID))
	values.Add("value", fmt.Sprintf("%v", rating.Value))
	req.URL.RawQuery = values.Encode()

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode/100 != 2 {
		return fmt.Errorf("non-2xx response: %v", resp)
	}

	return nil
}
