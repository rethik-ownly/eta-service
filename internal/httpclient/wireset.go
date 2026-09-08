package httpclient

import (
	"github.com/google/wire"
)

var Wireset = wire.NewSet(NewHTTPClient)