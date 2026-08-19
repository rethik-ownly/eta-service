package types

type GetETAResponse struct {
	EtaID  				 string `json:"etaId" bson:"etaId"`
	EtaInSeconds 		 uint `json:"eta_in_seconds" bson:"eta_in_seconds"`
	CreatedAt            float64  `json:"createdAt" bson:"createdAt"`
	UpdatedAt            float64  `json:"updatedAt" bson:"updatedAt"`
}

type InsertETARequest struct {
	RestaurantID string `json:"restaurant_id" bson:"restaurant_id"`
	DayOfWeek string `json:"day_of_week" bson:"day_of_week"`
	TimeSlot string `json:"time_slot" bson:"time_slot"`
	Lat float64 `json:"lat" bson:"lat"`
	Lon float64 `json:"lon" bson:"lon"`
	Eta uint `json:"eta" bson:"eta"`
	CreatedAt            float64  `json:"createdAt" bson:"createdAt"`
	UpdatedAt            float64  `json:"updatedAt" bson:"updatedAt"`
}