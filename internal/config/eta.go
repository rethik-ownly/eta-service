// Feature config
package config

type EtaConfig struct {
	// Max number of entities allowed in fetch Eta request body
	MaxEntities int `mapstructure:"maxEntities"`
}

// EtaDefaultEstimates holds hardcoded food-time fallbacks (seconds).
type EtaDefaultEstimates struct {
	RestaurantAcceptanceTime int `mapstructure:"restaurantAcceptanceTime"`
	KitchenPreparationTime   int `mapstructure:"kitchenPreparationTime"`
	PickupTime               int `mapstructure:"pickupTime"`
	CaptainAssignmentTime    int `mapstructure:"captainAssignmentTime"`
	FirstMileTime            int `mapstructure:"firstMileTime"`
	DelayDispatchTime        int `mapstructure:"delayDispatchTime"`
}