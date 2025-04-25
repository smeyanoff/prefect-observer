package repository

import (
	"context"
	"crm-uplift-ii24-backend/internal/domain/entity"
)

type StageRepository interface {
	SaveStage(ctx context.Context, stage *entity.Stage) error
	GetStageByID(ctx context.Context, stageID uint) (*entity.Stage, error)
	GetSendpostStages(ctx context.Context, sendpostID uint) ([]*entity.Stage, error)
	GetSubStages(ctx context.Context, parentStageID uint) ([]*entity.Stage, error)
	DeleteStage(ctx context.Context, stageID uint) error
	GetPreviousStage(ctx context.Context, stageID uint) (*entity.Stage, error)
}
