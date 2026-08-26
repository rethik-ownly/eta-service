package utils

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/nutanalabs/eta-service/internal/constants"
	logger "github.com/nutanalabs/rapido-logger-go"
)

type CommonUtils interface {
	GetDayFromTime(t time.Time) constants.Day
	GetMealTypeFromTime(t time.Time) constants.MealType

    ToJSON(data interface{}) string
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