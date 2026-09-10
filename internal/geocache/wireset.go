package geocache

import "github.com/google/wire"

var Wireset = wire.NewSet(
	NewGeocache,
	NewGeoRepository,
)
