package server

import "github.com/google/wire"

var Wireset = wire.NewSet(
	NewServer,
)