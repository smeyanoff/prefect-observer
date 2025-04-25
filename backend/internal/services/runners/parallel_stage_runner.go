package runners

import (
	"context"
	"crm-uplift-ii24-backend/internal/domain/entity"
	"crm-uplift-ii24-backend/internal/domain/value"
	"crm-uplift-ii24-backend/internal/services"
	"crm-uplift-ii24-backend/pkg/logging"
	"sync"

	"go.uber.org/zap"
)

type parallelStageRunner struct {
	stageRunnerService *services.StageRunnerService
	stageService       *services.StageService
}

func newParallelStageRunner(stageRunnerService *services.StageRunnerService, stageService *services.StageService) entity.StageRunner {
	return &parallelStageRunner{stageRunnerService: stageRunnerService, stageService: stageService}
}

func (psr *parallelStageRunner) Start(ctx context.Context, stage *entity.Stage) error {
	logging.Debug("[Stage Runner Parallel] Start", zap.Uint("stage_id", stage.ID))

	if err := psr.stageService.UpdateStageState(ctx, stage, value.Running); err != nil {
		return err
	}

	subStages, err := psr.stageService.GetSubStages(ctx, stage.ID)
	if err != nil {
		return err
	}

	var wg sync.WaitGroup
	errorsChan := make(chan error, len(subStages))

	for _, subStage := range subStages {
		wg.Add(1)
		go func(subStage *entity.Stage) {
			defer wg.Done()
			logging.Debug("[Stage Runner Parallel] Start", zap.Uint("sub_stage_id", subStage.ID))
			var err error
			if subStage.IsParallel() {
				err = psr.Start(ctx, subStage)
			} else {
				err = psr.stageRunnerService.Start(ctx, subStage)
			}

			if err != nil {
				errorsChan <- err
			}
		}(subStage)
	}

	go func() {
		wg.Wait()
		close(errorsChan)
	}()

	for err := range errorsChan {
		if err != nil {
			return psr.stageRunnerService.HandleFailedStage(ctx, stage, err)
		}
	}

	return nil
}

func (psr *parallelStageRunner) CheckState(ctx context.Context, stage *entity.Stage) error {
	logging.Debug("[Stage Runner Parallel] CheckState", zap.Uint("stage_id", stage.ID))

	subStages, err := psr.stageService.GetSubStages(ctx, stage.ID)
	if err != nil {
		return err
	}

	var wg sync.WaitGroup
	errorsChan := make(chan error, len(subStages))

	for _, subStage := range subStages {
		wg.Add(1)
		go func() {
			defer wg.Done()
			logging.Debug("[Stage Runner Parallel] CheckState", zap.Uint("sub_stage_id", subStage.ID))
			var err error
			if subStage.IsParallel() {
				err = psr.CheckState(ctx, subStage)
			} else {
				err = psr.stageRunnerService.CheckState(ctx, subStage)
			}

			if err != nil {
				errorsChan <- err
			}
		}()
	}

	go func() {
		wg.Wait()
		close(errorsChan)
	}()

	for err := range errorsChan {
		if err != nil {
			return psr.stageRunnerService.HandleFailedStage(ctx, stage, err)
		}
	}

	return psr.stageService.UpdateStageState(ctx, stage, value.Completed)
}
