package runners

import (
	"context"
	"crm-uplift-ii24-backend/internal/domain/entity"
	"crm-uplift-ii24-backend/internal/services"
	"crm-uplift-ii24-backend/pkg/logging"

	"go.uber.org/zap"
)

type sequentialStageRunner struct {
	StageRunnerService *services.StageRunnerService
}

func newSequentialStageRunner(StageRunnerService *services.StageRunnerService) entity.StageRunner {
	return &sequentialStageRunner{StageRunnerService: StageRunnerService}
}

func (ssr *sequentialStageRunner) Start(ctx context.Context, stage *entity.Stage) error {
	logging.Debug("[Stage Runner Sequential] Start", zap.Uint("stage_id", stage.ID))
	return ssr.StageRunnerService.Start(ctx, stage)
}

func (ssr *sequentialStageRunner) CheckState(ctx context.Context, stage *entity.Stage) error {
	logging.Debug("[Stage Runner Sequential] CheckState", zap.Uint("stage_id", stage.ID))
	return ssr.StageRunnerService.CheckState(ctx, stage)
}
