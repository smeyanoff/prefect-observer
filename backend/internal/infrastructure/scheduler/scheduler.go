package scheduler

import (
	"context"
	"crm-uplift-ii24-backend/internal/domain/entity"
	"sync"
	"time"
)

type sendpostScheduler struct {
	queue           chan entity.SendpostSchedule
	activeSchedules map[uint]context.CancelFunc
	mu              sync.Mutex
}

func NewSendpostScheduler() entity.SendpostScheduler {
	return &sendpostScheduler{
		queue:           make(chan entity.SendpostSchedule, 100),
		activeSchedules: make(map[uint]context.CancelFunc),
	}
}

func (q *sendpostScheduler) AddSchedule(ctx context.Context, schedule entity.SendpostSchedule) {
	q.mu.Lock()
	if cancel, exists := q.activeSchedules[schedule.SendpostID]; exists {
		cancel() // Отмена старой задачи, если расписание обновилось
	}
	_, cancel := context.WithCancel(ctx)
	q.activeSchedules[schedule.SendpostID] = cancel
	q.mu.Unlock()
	delay := time.Until(schedule.PlannedAt)
	if delay > 0 {
		go func() {
			time.Sleep(delay)
			q.mu.Lock()
			delete(q.activeSchedules, schedule.SendpostID)
			q.mu.Unlock()
			q.queue <- schedule
		}()
	} else {
		q.queue <- schedule
	}
}

func (q *sendpostScheduler) RemoveSchedule(sendpostID uint) {
	q.mu.Lock()
	if cancel, exists := q.activeSchedules[sendpostID]; exists {
		cancel()
		delete(q.activeSchedules, sendpostID)
	}
	q.mu.Unlock()
}

func (q *sendpostScheduler) GetQueue() <-chan entity.SendpostSchedule {
	return q.queue
}
