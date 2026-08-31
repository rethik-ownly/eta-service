package health

import (
	"github.com/google/wire"
)

var Wireset = wire.NewSet(
	NewHandler,
)

