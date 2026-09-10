package routingengine

import "github.com/google/wire"

var Wireset = wire.NewSet(
	NewRoutingEngineClient,
)