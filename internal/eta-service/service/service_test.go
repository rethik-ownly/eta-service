package service

import (
	"testing"

	"github.com/nutanalabs/eta-service/internal/config"
	"github.com/nutanalabs/eta-service/internal/constants"
	"github.com/nutanalabs/eta-service/internal/eta-service/repository"
	routingengine "github.com/nutanalabs/eta-service/internal/serviceclients/routing-engine"
	"github.com/nutanalabs/eta-service/internal/types"
	common "github.com/nutanalabs/eta-service/internal/utils/common"
	logger "github.com/nutanalabs/rapido-logger-go"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

type ServiceTestSuite struct {
	suite.Suite
	ctrl              *gomock.Controller
	config            *config.Config
	repoMock          *repository.MockRepository
	commonUtilsMock   *common.MockCommonUtils
	routingClientMock *routingengine.MockRoutingEngineClient
	svc               Service
}

func (s *ServiceTestSuite) SetupTest() {
	s.ctrl = gomock.NewController(s.T())
	s.config = config.InitConfig("test")
	logger.Init(s.config.Log.Level)

	s.repoMock = repository.NewMockRepository(s.ctrl)
	s.commonUtilsMock = common.NewMockCommonUtils(s.ctrl)
	s.routingClientMock = routingengine.NewMockRoutingEngineClient(s.ctrl)

	s.svc = NewService(s.repoMock, s.commonUtilsMock, s.routingClientMock, s.config)
}

func (s *ServiceTestSuite) TestFetchEta() {
	// Arrange
	request := &types.FetchEtaRequest{
		Surface:      "Search",
		DeliveryType: "Standard",
		UserID:       "1234",
		UserLocation: types.Location{
			Lat: 12.9716,
			Lng: 77.5946,
		},
		Options: types.FetchEtaRequestOptions{
			QosLevel: "qos3",
		},
		Entities: []types.FetchEtaRequestEntity{
			{
				RestaurantID:       "rest-1",
				RestaurantLocation: types.Location{Lat: 12.9352, Lng: 77.6245},
			},
		},
	}

	s.commonUtilsMock.EXPECT().GetDayFromTime(gomock.Any()).Return(constants.Monday).Times(1)
	s.commonUtilsMock.EXPECT().GetMealTypeFromTime(gomock.Any()).Return(constants.Lunch).Times(1)

	restaurantComponents := []types.RestaurantComponents{
		{
			RestaurantId:  "rest-1",
			SublocalityId: "sub-1",
			Lunch: types.RestaurantMealComponents{
				Rat:           types.TimeSample{Seconds: 100, SampleCount: 10},
				Kpt:           types.TimeSample{Seconds: 500, SampleCount: 10},
				Pickup:        types.TimeSample{Seconds: 200, SampleCount: 10},
				DelayDispatch: types.TimeSample{Seconds: 0, SampleCount: 10},
			},
		},
	}

	s.repoMock.
		EXPECT().
		FetchRestaurantComponentsByIDs([]string{"rest-1"}, constants.Monday, constants.Lunch, s.config.Mongo.QueryBatchSize).
		Return(restaurantComponents, nil).
		Times(1)

	sublocalityComponents := []types.SublocalityComponents{
		{
			SublocalityId: "sub-1",
			Lunch: types.SublocalityMealComponents{
				Cat: types.TimeSample{Seconds: 150, SampleCount: 10},
				Fm:  types.TimeSample{Seconds: 250, SampleCount: 10},
			},
		},
	}

	s.repoMock.
		EXPECT().
		FetchSublocalityComponentsByIDs([]string{"sub-1"}, constants.Monday, constants.Lunch, s.config.Mongo.QueryBatchSize).
		Return(sublocalityComponents, nil).
		Times(1)

	distanceMatrixResponse := &routingengine.DistanceMatrixResponse{
		Data: [][]routingengine.MatrixElement{
			{
				{
					Duration: routingengine.DurationInfo{Value: 900, Unit: "seconds"},
				},
			},
		},
	}

	s.routingClientMock.
		EXPECT().
		GetDistanceMatrixWithQoS(gomock.Any(), constants.QosThree).
		Return(distanceMatrixResponse, nil).
		Times(1)

	s.repoMock.EXPECT().
		PublishFetchEtaEvent(
			constants.EventTypeNewEta,
			gomock.Any(),
			gomock.Any(),
			nil,
		).
		Times(1)

	// Act
	response, err := s.svc.FetchEta(request)

	// Assert
	s.Require().NoError(err)
	s.Equal("rest-1", response[0].RestaurantID)
	s.Equal(uint(1600), response[0].EtaInSeconds)
	s.Equal(constants.EtaSourceHistoric, response[0].Source)
}

func TestServiceTestSuite(t *testing.T) {
	suite.Run(t, new(ServiceTestSuite))
}
