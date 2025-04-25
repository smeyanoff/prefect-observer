package services

import (
	"crm-uplift-ii24-backend/internal/domain/entity"
	"crm-uplift-ii24-backend/internal/domain/value"
	runstatus "crm-uplift-ii24-backend/internal/infrastructure/notifications/runStatus"
	"crm-uplift-ii24-backend/pkg/logging"
	"errors"
	"time"

	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

type SenpostRunNotificationService struct {
	sendopostsToNotify map[uint]entity.SendopostRunNotificator
}

func NewSenpostRunNotificationService(notificator entity.SendopostRunNotificator) *SenpostRunNotificationService {
	return &SenpostRunNotificationService{sendopostsToNotify: make(map[uint]entity.SendopostRunNotificator)}
}

func (s *SenpostRunNotificationService) AddRunSendpostToNotify(sendpostID uint) {
	s.sendopostsToNotify[sendpostID] = runstatus.NewNotificatorWS()
}

func (s *SenpostRunNotificationService) RemoveRunSendpostToNotify(sendpostID uint) {
	s.sendopostsToNotify[sendpostID].StopNotificate()
	delete(s.sendopostsToNotify, sendpostID)
}

func (s *SenpostRunNotificationService) NotifyRunSendpost(sendpostID uint, status value.StateType) error {
	if notificator, ok := s.sendopostsToNotify[sendpostID]; ok {
		notificator.NotifyListeners(status.String())
		return nil
	}
	return errors.New("[SenpostRunNotificationService] NotifyRunSendpost: sendpost not found")
}

func (s *SenpostRunNotificationService) AddListener(sendpostID uint, listener entity.Listener) error {
	for i := 0; i < 5; i++ {
		if notificator, ok := s.sendopostsToNotify[sendpostID]; ok {
			notificator.AddListener(listener)
			return nil
		}
		logging.Info("[SenpostRunNotificationService] AddListener: Trying to add listener to sendpost", zap.Int("attempt", i+1))
		time.Sleep(time.Second / 3)
	}

	listener.WriteMessage(websocket.TextMessage, []byte("sendpost not found"))

	return errors.New("[SenpostRunNotificationService] AddListener: sendpost not found")
}
