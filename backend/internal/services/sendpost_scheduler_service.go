package services

import (
	"context"
	"crm-uplift-ii24-backend/internal/domain/entity"
	"crm-uplift-ii24-backend/pkg/logging"
	"sync"

	"go.uber.org/zap"
)

type SendpostSchedulerService struct {
	scheduler             entity.SendpostScheduler
	ScheduleService       ScheduleService
	SendpostRunnerService SendpostRunnerService
	workerCount           int
	wg                    sync.WaitGroup
}

func NewSendpostSchedulerService(scheduler entity.SendpostScheduler, ScheduleService ScheduleService, SendpostRunnerService SendpostRunnerService, workerCount int) *SendpostSchedulerService {
	return &SendpostSchedulerService{
		scheduler:             scheduler,
		ScheduleService:       ScheduleService,
		SendpostRunnerService: SendpostRunnerService,
		workerCount:           workerCount,
	}
}

func (s *SendpostSchedulerService) Start(ctx context.Context) {
	for i := 0; i < s.workerCount; i++ {
		s.wg.Add(1)
		go s.worker(ctx)
	}
}

func (s *SendpostSchedulerService) worker(ctx context.Context) {
	defer s.wg.Done()
	for {
		select {
		case schedule := <-s.scheduler.GetQueue():
			logging.Info("[SchedulerService] Start Sendpost ID: %d", zap.Uint("sendpost_id", schedule.SendpostID))

			// Обновляем, что выполнение началось
			if err := s.ScheduleService.UpdateScheduleStartedAt(ctx, &schedule); err != nil {
				logging.Error("[SchedulerService] Error updating start schedule", zap.Error(err))
				continue
			}

			// Запускаем выполнение рассылки
			s.SendpostRunnerService.Start(ctx, schedule.SendpostID)

			// Обновляем, что выполнение завершилось
			if err := s.ScheduleService.UpdateScheduleCompletedAt(ctx, &schedule); err != nil {
				logging.Error("[SchedulerService] Error updating completed stage", zap.Error(err))
			}

		case <-ctx.Done():
			return
		}
	}
}

func (s *SendpostSchedulerService) Wait() {
	s.wg.Wait()
}
