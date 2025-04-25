package entity

import (
	"context"
	"crm-uplift-ii24-backend/internal/domain/value"
)

type StageRunner interface {
	Start(ctx context.Context, stage *Stage) error
	CheckState(ctx context.Context, stage *Stage) error
}

type StageRunnerFactory interface {
	CreateRunner(stageType value.StageType) StageRunner
}
