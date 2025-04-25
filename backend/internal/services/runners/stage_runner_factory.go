package runners

import (
	"crm-uplift-ii24-backend/internal/domain/entity"
	"crm-uplift-ii24-backend/internal/domain/value"
	"crm-uplift-ii24-backend/internal/services"
)

type stageRunnerFactory struct {
	StageRunnerService *services.StageRunnerService
	stageService       *services.StageService
}

func NewStageRunnerFactory(StageRunnerService *services.StageRunnerService, stageService *services.StageService) entity.StageRunnerFactory {
	return &stageRunnerFactory{StageRunnerService: StageRunnerService, stageService: stageService}
}

func (srf *stageRunnerFactory) CreateRunner(stageType value.StageType) entity.StageRunner {
	switch stageType {
	case value.ParallelStage:
		return newParallelStageRunner(srf.StageRunnerService, srf.stageService)
	case value.ObserverStage:
		return newObserverStageRunner(srf.StageRunnerService, srf.stageService)
	case value.SequentialStage:
		return newSequentialStageRunner(srf.StageRunnerService)
	default:
		return newSequentialStageRunner(srf.StageRunnerService)
	}
}
