package services

import (
	"context"
	"testing"

	"gorm.io/gorm"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	"crm-uplift-ii24-backend/internal/domain/entity"
	"crm-uplift-ii24-backend/internal/domain/value"
	"crm-uplift-ii24-backend/internal/mocks"
	"crm-uplift-ii24-backend/pkg/logging"

	"go.uber.org/zap"
)

type SendpostServiceTestSuite struct {
	suite.Suite
	svc          *SendpostService
	sendpostRepo *mocks.SendpostRepository
	stageRepo    *mocks.StageRepository
}

func (s *SendpostServiceTestSuite) SetupTest() {
	s.sendpostRepo = new(mocks.SendpostRepository)
	s.stageRepo = new(mocks.StageRepository)
	stageService := NewStageService(s.stageRepo, s.sendpostRepo)
	s.svc = NewSendpostService(s.sendpostRepo, stageService)
	logging.Logger = zap.NewNop()

	sendpostID := uint(1)
	existingParams := value.JSONB{"foo": "bar"}
	sendpost := &entity.Sendpost{Model: gorm.Model{ID: sendpostID}, GlobalParameters: &existingParams}

	s.sendpostRepo.On("GetSendpostByID", mock.Anything, sendpostID).Return(sendpost, nil)
	s.sendpostRepo.On("GetSendpostParameters", mock.Anything, sendpostID).Return(&existingParams, nil)
}

func (s *SendpostServiceTestSuite) TestAddUpdateSendpostParameter() {

	s.sendpostRepo.On("SaveSendpost", mock.Anything, mock.MatchedBy(func(s *entity.Sendpost) bool {
		params := *s.GlobalParameters
		return params["foo"] == "bar" &&
			params["new_key"] == "new_value"
	})).Return(nil)

	err := s.svc.AddUpdateSendpostParameter(context.Background(), uint(1), "new_key", "new_value")
	assert.NoError(s.T(), err)
	s.sendpostRepo.AssertExpectations(s.T())

	s.sendpostRepo.On("SaveSendpost", mock.Anything, mock.MatchedBy(func(s *entity.Sendpost) bool {
		params := *s.GlobalParameters
		return params["foo"] == "foos"
	})).Return(nil)

	err = s.svc.AddUpdateSendpostParameter(context.Background(), uint(1), "foo", "foos")
	assert.NoError(s.T(), err)
	s.sendpostRepo.AssertExpectations(s.T())
}

func (s *SendpostServiceTestSuite) TestDeleteSendpostParameter() {

	s.sendpostRepo.On("SaveSendpost", mock.Anything, mock.MatchedBy(func(s *entity.Sendpost) bool {
		params := *s.GlobalParameters
		if _, ok := params["foo"]; ok {
			return false
		}
		return true
	})).Return(nil)

	err := s.svc.DeleteSendpostParameter(context.Background(), uint(1), "foo")
	assert.NoError(s.T(), err)
	s.sendpostRepo.AssertExpectations(s.T())
}

func (s *SendpostServiceTestSuite) TestGetSendpostParameters() {
	params, err := s.svc.GetSendpostParameters(context.Background(), uint(1))
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), value.JSONB{"foo": "bar"}, *params)
}

func TestSendpostServiceTestSuite(t *testing.T) {
	suite.Run(t, new(SendpostServiceTestSuite))
}
