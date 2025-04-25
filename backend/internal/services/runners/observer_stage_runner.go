package runners

import (
	"context"
	"crm-uplift-ii24-backend/internal/domain/entity"
	"crm-uplift-ii24-backend/internal/domain/value"
	"crm-uplift-ii24-backend/internal/services"
	"crm-uplift-ii24-backend/pkg/logging"
	"time"

	"go.uber.org/zap"
)

type observerStageRunner struct {
	stageRunnerService *services.StageRunnerService
	stageService       *services.StageService
}

func newObserverStageRunner(stageRunnerService *services.StageRunnerService, stageService *services.StageService) entity.StageRunner {
	return &observerStageRunner{stageRunnerService: stageRunnerService, stageService: stageService}
}

func (ost *observerStageRunner) Start(ctx context.Context, stage *entity.Stage) error {
	logging.Debug("[Stage Runner Observer] Start", zap.Uint("stage_id", stage.ID))

	return ost.stageService.UpdateStageState(ctx, stage, value.Running)
}

func (ost *observerStageRunner) CheckState(ctx context.Context, stage *entity.Stage) error {
	logging.Debug("[Stage Runner Observer] CheckState", zap.Uint("stage_id", stage.ID))

	now := time.Now()
	logging.Debug("[Stage Runner Observer] CheckState", zap.Time("current_time", now))

	return ost.stageRunnerService.CheckStageCompledSuccesfullyInPeriod(
		ctx,
		stage,
		time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location()),
		time.Now(),
	)
}
