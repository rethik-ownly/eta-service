package constants

// Surface - Point of ingression
type Surface string

const (
	Search     Surface = "SEARCH"
	Feed       Surface = "FEED"
	Restaurant Surface = "RESTAURANT"
	Cart       Surface = "CART"
	Promise    Surface = "PROMISE"
)

func (s Surface) IsValid() bool {
	switch s {
	case Search, Feed, Restaurant, Cart, Promise:
		return true
	default:
		return false
	}
}

// Day_Type
type Day string

const (
	Monday    Day = "MONDAY"
	Tuesday   Day = "TUESDAY"
	Wednesday Day = "WEDNESDAY"
	Thursday  Day = "THURSDAY"
	Friday    Day = "FRIDAY"
	Saturday  Day = "SATURDAY"
	Sunday    Day = "SUNDAY"
	Weekday	  Day = "WEEKDAY"
	Weekend   Day = "WEEKEND"
)

func (d Day) IsValid() bool {
	switch d {
	case Monday, Tuesday, Wednesday, Thursday, Friday, Saturday, Sunday, Weekday, Weekend:
		return true
	default:
		return false
	}
}

// Meal_Type
type MealType string

const (
	Breakfast MealType = "BREAKFAST"
	Lunch     MealType = "LUNCH"
	Dinner    MealType = "DINNER"
	LateNight MealType = "LATE_NIGHT"
	Snack     MealType = "SNACK"
)

func (mt MealType) IsValid() bool {
	switch mt {
	case Breakfast, LateNight, Lunch, Dinner, Snack:
		return true
	default:
		return false
	}
}

// BsonKey returns the nested bson field name for this meal type within the
// eta_restaurant_estimates / eta_sublocality_estimates documents (e.g.
// "breakfast", "snacks", "latenight").
func (mt MealType) BsonKey() string {
	switch mt {
	case Breakfast:
		return "breakfast"
	case Lunch:
		return "lunch"
	case Snack:
		return "snacks"
	case Dinner:
		return "dinner"
	case LateNight:
		return "latenight"
	default:
		return ""
	}
}

// eta_delivery_type
type DeliveryType string

const (
	Standard DeliveryType = "STANDARD"
	Express  DeliveryType = "EXPRESS"
)

func (dt DeliveryType) IsValid() bool {
	switch dt {
	case Standard, Express:
		return true
	default:
		return false
	}
}

// Qos_Level
type QosLevel string

const (
	QosOne   QosLevel = "qos1"
	QosTwo   QosLevel = "qos2"
	QosThree QosLevel = "qos3"
)

func (q QosLevel) IsValid() bool {
	switch q {
	case QosOne, QosTwo, QosThree:
		return true
	default:
		return false
	}
}

// EventType tags a FetchEta analytics event.
type EventType string

const (
	EventTypeNewEta EventType = "NEW_ETA"
)
