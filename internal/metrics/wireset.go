package metrics

import "github.com/google/wire"

var Wireset = wire.NewSet(
	NewMetrics,
)