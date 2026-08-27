package routingengine

import "github.com/nutanalabs/eta-service/internal/types"

type DistanceMatrixRequest struct {
	Sources                []types.Location              `json:"sources"`
	Destinations           []types.Location            `json:"destinations"`
	SourceDestinationPairs []SourceDestinationPair `json:"sourceDestinationPairs,omitempty"`
	Vehicle                string                  `json:"vehicle,omitempty"`
	RoutingPreference      string                  `json:"routingPreference,omitempty"`
}

type SourceDestinationPair struct {
	Source      int `json:"source"`
	Destination int `json:"destination"`
}

type DistanceMatrixResponse struct {
	Data [][]MatrixElement `json:"data"`
	Provider string `json:"provider,omitempty"`
}

type MatrixElement struct {
	Distance DistanceInfo `json:"distance"`
	Duration DurationInfo `json:"duration"`
}

type DistanceInfo struct {
	Value float64 `json:"value"`
	Unit string `json:"unit"`
}

type DurationInfo struct {
	Value float64 `json:"value"`
	Unit string `json:"unit"`
}