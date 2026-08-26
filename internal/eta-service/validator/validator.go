package validator

import (
	"github.com/nutanalabs/eta-service/internal/config"
	"github.com/nutanalabs/eta-service/internal/types"
)

type Validator interface {
	ValidateFetchEtaRequest(request *types.FetchEtaRequest) error
}

type validatorImpl struct {
	config *config.Config
}

func NewValidator(config *config.Config) Validator {
	return &validatorImpl{
		config: config,
	}
} 

func (v *validatorImpl) ValidateFetchEtaRequest(request *types.FetchEtaRequest) error {
	if !request.Surface.IsValid(){
		return types.NewBadRequestError("invalid surface")
	}
	if !request.DeliveryType.IsValid(){
		return types.NewBadRequestError("invalid delivery type")
	}
	if request.Options.QosLevel != "" && !request.Options.QosLevel.IsValid() {
		return types.NewBadRequestError("invalid qosLevel")
	}
	if request.UserID == "" {
		return types.NewBadRequestError("invalid userId")
	}
	if request.Entities == nil || len(request.Entities) > v.config.Eta.MaxEntities {
		return types.NewBadRequestError("invalid entities")
	}
	return nil
}