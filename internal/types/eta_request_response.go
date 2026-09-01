package types

import (
	"github.com/nutanalabs/eta-service/internal/constants"
)

type Location struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}

// Fetch Eta Reques-Response
type FetchEtaRequestOptions struct {
	QosLevel constants.QosLevel `json:"qosLevel"`
}

type FetchEtaRequestEntity struct {
	RestaurantID       string   `json:"restaurantId"`
	RestaurantLocation Location `json:"restaurantLocation"`
}

type FetchEtaRequest struct {
	Surface      constants.Surface       `json:"surface"`
	DeliveryType constants.DeliveryType  `json:"deliveryType"`
	UserID       string                  `json:"userId" binding:"required"`
	UserLocation Location                `json:"userLocation" binding:"required"`
	Options      FetchEtaRequestOptions  `json:"options"`
	Entities     []FetchEtaRequestEntity `json:"entities"`
	OrderId      string                  `json:"-"`
	RequestId    string                  `json:"-"`
}

type FetchEtaResponse struct {
	RestaurantID string `json:"restaurantId"`
	EtaInSeconds uint   `json:"etaInSeconds"`
	DisplayMin   uint   `json:"displayMin,omitempty"`
	DisplayMax   uint   `json:"displayMax,omitempty"`
}

// Insert Eta Request
//
// Insert creates a brand-new document covering the given set of days
// (DayType/Day). It fails with a conflict if any of those days already
// belong to another document for the same id (see service layer).
type InsertRestaurantEstimateRequest struct {
	RestaurantId  string          `json:"restaurantId"`
	DayType       []constants.Day `json:"dayType"`
	CityId        string          `json:"cityId"`
	ZoneId        string          `json:"zoneId"`
	SublocalityId string          `json:"sublocalityId"`

	Breakfast RestaurantMealEstimate `json:"breakfast,omitempty"`
	Lunch     RestaurantMealEstimate `json:"lunch,omitempty"`
	Snacks    RestaurantMealEstimate `json:"snacks,omitempty"`
	Dinner    RestaurantMealEstimate `json:"dinner,omitempty"`
	Latenight RestaurantMealEstimate `json:"latenight,omitempty"`

	UpdatedAt float64 `json:"-"`
}

type InsertSublocalityEstimateRequest struct {
	SublocalityId string          `json:"sublocalityId"`
	Day           []constants.Day `json:"day"`
	ZoneId        string          `json:"zoneId,omitempty"`
	CityId        string          `json:"cityId"`

	Breakfast SublocalityMealEstimate `json:"breakfast"`
	Lunch     SublocalityMealEstimate `json:"lunch"`
	Snacks    SublocalityMealEstimate `json:"snacks"`
	Dinner    SublocalityMealEstimate `json:"dinner"`
	Latenight SublocalityMealEstimate `json:"latenight"`

	UpdatedAt float64 `json:"-"`
}

// RestaurantId is not bound from the request body; it is taken from the
// ":restaurantId" path param.
//
// Day locates the existing document (it must be present in the document's
// dayType array); Update never creates a document. MealType is only
// required when patching one of the nested Rat/Kpt/Pickup/DelayDispatch sections. Days,
// if provided, replaces the whole dayType array (e.g. to add another day to
// this same document).
type UpdateRestaurantEstimateRequest struct {
	RestaurantId string `json:"-"`

	Day      constants.Day      `json:"day" binding:"required"`
	MealType constants.MealType `json:"mealType,omitempty"`

	Days          *[]constants.Day `json:"days,omitempty"`
	CityId        *string          `json:"cityId,omitempty"`
	ZoneId        *string          `json:"zoneId,omitempty"`
	SublocalityId *string          `json:"sublocalityId,omitempty"`

	RatSeconds               *float64 `json:"ratSeconds,omitempty"`
	RatSampleCount           *int     `json:"ratSampleCount,omitempty"`
	KptSeconds               *float64 `json:"kptSeconds,omitempty"`
	KptSampleCount           *int     `json:"kptSampleCount,omitempty"`
	PickupSeconds            *float64 `json:"pickupSeconds,omitempty"`
	PickupSampleCount        *int     `json:"pickupSampleCount,omitempty"`
	DelayDispatchSeconds     *float64 `json:"delayDispatchSeconds,omitempty"`
	DelayDispatchSampleCount *int     `json:"delayDispatchSampleCount,omitempty"`

	UpdatedAt float64 `json:"-"`
}

// SublocalityId is not bound from the request body; it is taken from the
// ":sublocalityId" path param.
//
// Day locates the existing document (it must be present in the document's
// day array); Update never creates a document. MealType is only required
// when patching one of the nested Cat/Fm sections. Days, if provided,
// replaces the whole day array (e.g. to add another day to this same
// document).
type UpdateSublocalityEstimateRequest struct {
	SublocalityId string `json:"-"`

	Day      constants.Day      `json:"day" binding:"required"`
	MealType constants.MealType `json:"mealType,omitempty"`

	Days   *[]constants.Day `json:"days,omitempty"`
	ZoneId *string          `json:"zoneId,omitempty"`
	CityId *string          `json:"cityId,omitempty"`

	CatSeconds     *float64 `json:"catSeconds,omitempty"`
	CatSampleCount *int     `json:"catSampleCount,omitempty"`
	FmSeconds      *float64 `json:"fmSeconds,omitempty"`
	FmSampleCount  *int     `json:"fmSampleCount,omitempty"`

	UpdatedAt float64 `json:"-"`
}
