package utils

import (
	"encoding/json"
	"fmt"
	"math"
	"time"

	"github.com/nutanalabs/eta-service/internal/constants"
	logger "github.com/nutanalabs/rapido-logger-go"
)

type CommonUtils interface {
	GetDayFromTime(t time.Time) constants.Day
	GetMealTypeFromTime(t time.Time) constants.MealType
    ToJSON(data interface{}) string

    GetHaversineDistance(lat1, lng1, lat2, lng2 float64) float64 
}

type commonUtilImpl struct {}

func NewCommonUtils() CommonUtils {
	return &commonUtilImpl{}
}

func (c *commonUtilImpl) GetDayFromTime(t time.Time) constants.Day {
	switch t.Weekday() {
    case time.Monday:
        return constants.Monday
    case time.Tuesday:
        return constants.Tuesday
    case time.Wednesday:
        return constants.Wednesday
    case time.Thursday:
        return constants.Thursday
    case time.Friday:
        return constants.Friday
    case time.Saturday:
        return constants.Saturday
    case time.Sunday:
        return constants.Sunday
	default:
        return constants.Monday
    }
}

func (c *commonUtilImpl) GetMealTypeFromTime(t time.Time) constants.MealType {
	hour := t.Hour()
    switch {
    case hour >= 6 && hour < 11:
        return constants.Breakfast
    case hour >= 11 && hour < 16:
        return constants.Lunch
    case hour >= 16 && hour < 19:
        return constants.Snack
    case hour >= 19 && hour < 23:
        return constants.Dinner
    default:
        return constants.LateNight
    }
}

func (c *commonUtilImpl) ToJSON(data interface{}) string {
    b, err := json.Marshal(data)
	if err != nil {
		logger.Error(logger.Format{
			Message: fmt.Sprintf("Error in ToJSON - %s", err),
		})
		return ""
	}
	return string(b)
}

// Lat, lng in degree and output is distance in km
func (c *commonUtilImpl) GetHaversineDistance(lat1, lng1, lat2, lng2 float64) float64 {
    const EarthRadius = 6371

	toRad := func(deg float64) float64 {
		return deg * math.Pi / 180
	}

	deltaLat := toRad(lat2 - lat1)
	deltaLng := toRad(lng2 - lng1)

	lat1Rad := toRad(lat1)
	lat2Rad := toRad(lat2)

	a := math.Sin(deltaLat/2)*math.Sin(deltaLat/2) +
		    math.Cos(lat1Rad)*math.Cos(lat2Rad)*math.Sin(deltaLng/2)*math.Sin(deltaLng/2)

	d := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	return EarthRadius * d
}