package model

import "moviehub.com/metadata/pkg/model"

// MovieDetails wraps movie metadata and aggregated
// rating.
type MovieDetails struct {
	Rating *float64 `json:"rating,omitempty"`
	Metadata model.Metadata `json:"metadata"`
}