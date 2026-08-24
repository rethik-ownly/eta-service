package constants

// Surface - Point of ingression
type Surface string

const (
	Search Surface = "SEARCH"
	Feed Surface = "FEED"
	Restaurant Surface = "RESTAURANT"
	Cart Surface = "CART"
	Promise Surface = "PROMISE"
)

// Day_Type
type Day string

const (
	Monday Day = "MONDAY"
	Tuesday Day = "TUESDAY"
	Wednesday Day = "WEDNESDAY"
	Thursday Day = "THURSDAY"
	Friday Day = "FRIDAY"
	Saturday Day = "SATURDAY"
	Sunday Day = "SUNDAY"
)

// Meal_Type
type MealType string

const (
	Breakfast MealType = "BREAKFAST"
	Lunch MealType = "LUNCH"
	Dinner MealType = "DINNER"
	LateNight MealType = "LATE_NIGHT"
	Snack MealType = "SNACK"
)

// eta_delivery_type
type DeliveryType string

const (
	Standard DeliveryType = "STANDARD"
	Express DeliveryType = "EXPRESS"
)

// Qos_Level
type QosLevel string

const (
	QosOne QosLevel = "qos1"
	QosTwo QosLevel = "qos2"
	QosThree QosLevel = "qos3"
)