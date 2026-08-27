package utils

import (
	"github.com/google/wire"
	commonUtils "github.com/nutanalabs/eta-service/internal/utils/common"
	httpUtils "github.com/nutanalabs/eta-service/internal/utils/http"  
)

var Wireset = wire.NewSet(
	// NewUtils,
	commonUtils.NewCommonUtils,
	httpUtils.NewHttpUtils,
)