package runstatus

import (
	"crm-uplift-ii24-backend/internal/domain/entity"
	"crm-uplift-ii24-backend/pkg/logging"
	"sync"

	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

type NotificatorWS struct {
	pool map[entity.Listener]bool
	mu   sync.Mutex
}

func NewNotificatorWS() entity.SendopostRunNotificator {
	return &NotificatorWS{
		pool: make(map[entity.Listener]bool),
	}
}

func (n *NotificatorWS) AddListener(listener entity.Listener) {
	logging.Debug("[NotificatorWS] AddListener")
	n.mu.Lock()
	defer n.mu.Unlock()
	n.pool[listener] = true
}

func (n *NotificatorWS) RemoveListener(listener entity.Listener) {
	logging.Debug("[NotificatorWS] RemoveListener")
	n.mu.Lock()
	defer n.mu.Unlock()
	if err := listener.Close(); err != nil {
		logging.Error("[NotificatorWS] RemoveListener", zap.Error(err))
	}
	delete(n.pool, listener)
}

func (n *NotificatorWS) NotifyListeners(msg string) {
	logging.Debug("[NotificatorWS] NotifyListeners", zap.String("msg", msg))
	for ws := range n.pool {
		if err := ws.WriteMessage(websocket.TextMessage, []byte(msg)); err != nil {
			logging.Error("[NotificatorWS] NotifyListeners", zap.Error(err))
			n.RemoveListener(ws)
		}
	}
}

func (n *NotificatorWS) StopNotificate() {
	logging.Debug("[NotificatorWS] StopNotificate")
	for ws := range n.pool {
		n.RemoveListener(ws)
	}
}
