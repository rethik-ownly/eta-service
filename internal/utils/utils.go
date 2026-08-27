package utils

import (
	common "github.com/nutanalabs/eta-service/internal/utils/common"
	http "github.com/nutanalabs/eta-service/internal/utils/http"
)

type Utils interface {
	GetCommonUtil() common.CommonUtils
	GetHttpUtil() http.HTTPUtils
}

type UtilsImpl struct {
	commonUtils common.CommonUtils
	httpUtils http.HTTPUtils
}

func NewUtils() Utils {
	return &UtilsImpl{
		commonUtils: common.NewCommonUtils(),
		httpUtils:   http.NewHttpUtils(),
	}
}

func (u *UtilsImpl) GetCommonUtil() common.CommonUtils {
	return u.commonUtils
}

func (u *UtilsImpl) GetHttpUtil() http.HTTPUtils {
	return u.httpUtils
}