package types

import "github.com/nutanalabs/eta-service/internal/constants"

type FetchEtaAnalyticsEvent struct {
	EventId      string                 `json:"eventId"`
	RequestId    string                 `json:"requestId"`
	CreatedAt    int64                  `json:"createdAt"`
	EventType    constants.EventType    `json:"eventType"`
	Surface      constants.Surface      `json:"surface"`
	DeliveryType constants.DeliveryType `json:"deliveryType"`
	Status       string                 `json:"status"`
	ErrorMessage string                 `json:"errorMessage,omitempty"`
	Request      FetchEtaRequest        `json:"request"`
	Response     []FetchEtaResponse     `json:"response,omitempty"`
}
