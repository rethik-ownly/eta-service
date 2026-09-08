package routingengine

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/nutanalabs/eta-service/internal/config"
	"github.com/nutanalabs/eta-service/internal/constants"
	"github.com/nutanalabs/rapido-http-go/v3/httpclient"
)

type RoutingEngineClient interface {
	GetDistanceMatrix(request *DistanceMatrixRequest) (*DistanceMatrixResponse, error)
	GetDistanceMatrixWithQoS(request *DistanceMatrixRequest, qos constants.QosLevel) (*DistanceMatrixResponse, error)
}

type routingEngineClientImpl struct {
	config *config.Config
	httpClient httpclient.Client
}

func NewRoutingEngineClient(config *config.Config, httpClient httpclient.Client) RoutingEngineClient {
	return &routingEngineClientImpl{
		config: config,
		httpClient: httpClient,
	}
}

func (rec *routingEngineClientImpl) GetDistanceMatrix(request *DistanceMatrixRequest) (*DistanceMatrixResponse, error) {
	return rec.GetDistanceMatrixWithQoS(request, rec.config.ExternalServices.RoutingEngine.DistanceMatrixAPI.QosLevel)
}

func (rec *routingEngineClientImpl) GetDistanceMatrixWithQoS(request *DistanceMatrixRequest, qos constants.QosLevel) (*DistanceMatrixResponse, error) {
	if request == nil {
		return nil, fmt.Errorf("request cannot be nil")
	}
	if len(request.Sources) == 0 {
		return nil, fmt.Errorf("sources cannot be empty")
	}
	if len(request.Destinations) == 0 {
		return nil, fmt.Errorf("destinations cannot be empty")
	}

	if request.Vehicle == "" {
		request.Vehicle = constants.VEHICLE_TWO_WHEELER
	}

	if request.RoutingPreference == "" {
		request.RoutingPreference = constants.ROUTING_PREFERENCE_TRAFFIC_AWARE
	}

	routingEngineURL := fmt.Sprintf("%s:%d%s?qos=%s",
		rec.config.ExternalServices.RoutingEngine.Host,
		rec.config.ExternalServices.RoutingEngine.Port,
		rec.config.ExternalServices.RoutingEngine.DistanceMatrixAPI.Path,
		string(qos))

	requestBody, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %w", err)
	}

	httpRequest := httpclient.
		NewRequestBuilder().
		SetHeader("content-type", "application/json").
		SetURL(routingEngineURL).
		SetBody(requestBody).
		SetMaxTimeout(time.Duration(rec.config.ExternalServices.RoutingEngine.DistanceMatrixAPI.TimeoutInMs)*time.Millisecond).
		Build()
	
	httpResponse, err := httpRequest.Post(context.Background(), rec.httpClient)

	if err != nil {
		return nil, fmt.Errorf("failed to call routing engine API: %w", err)
	}

	var distanceMatrixResponse DistanceMatrixResponse

	if err := json.Unmarshal(httpResponse.Body, &distanceMatrixResponse); err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	return &distanceMatrixResponse, nil
}