package repository

import (
	"context"
	"crm-uplift-ii24-backend/internal/domain/entity"
	"crm-uplift-ii24-backend/internal/domain/repository"
	"crm-uplift-ii24-backend/pkg/logging"
	"errors"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type gormSendpostStageRepository struct {
	db *gorm.DB
}

// NewGormSendpostStageRepository creates a new instance of gormSendpostStageRepository
// using the provided gorm.DB connection. It returns an implementation of the
// repository.StageRepository interface.
func NewGormSendpostStageRepository(db *gorm.DB) repository.StageRepository {
	return &gormSendpostStageRepository{db: db}
}

// SaveStage saves the given stage entity to the database using the provided context.
// It returns an error if the operation fails.
//
// Parameters:
//
//	ctx - The context to use for the database operation.
//	stage - A pointer to the Stage entity to be saved.
//
// Returns:
//
//	An error if the save operation fails, otherwise nil.
func (ssr *gormSendpostStageRepository) SaveStage(ctx context.Context, stage *entity.Stage) error {
	logging.Debug("[Stage repo] SaveStage", zap.Any("stage", stage))
	return ssr.db.WithContext(ctx).Save(stage).Error
}

// GetStageByID retrieves a Stage entity by its ID from the database.
// It uses the provided context for database operations and returns
// an error if the stage is not found or if any other database error occurs.
//
// Parameters:
//
//	ctx - The context for managing request-scoped values, cancellation signals, and deadlines.
//	stageID - The unique identifier of the stage to be retrieved.
//
// Returns:
//
//	*entity.Stage - The retrieved stage entity, or nil if not found.
//	error - An error if the stage is not found or if a database error occurs.
func (ssr *gormSendpostStageRepository) GetStageByID(ctx context.Context, stageID uint) (*entity.Stage, error) {
	var stage *entity.Stage

	if err := ssr.db.WithContext(ctx).
		First(&stage, stageID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("sendpost not found")
		} else {
			return nil, err
		}
	}
	logging.Debug("[Stage repo] GetStageByID", zap.Any("stage", stage))
	return stage, nil
}

// GetSendpostStages retrieves all stages associated with a given sendpost ID
// from the database, filtering out stages with a parent stage ID. The stages
// are ordered by the next stage ID in descending order. It returns a slice of
// Stage entities or an error if the operation fails.
//
// Parameters:
//
//	ctx - The context for managing request-scoped values, cancellation, and deadlines.
//	sendpostID - The ID of the sendpost for which stages are to be retrieved.
//
// Returns:
//
//	A slice of pointers to Stage entities or an error if the retrieval fails.
func (ssr *gormSendpostStageRepository) GetSendpostStages(ctx context.Context, sendpostID uint) ([]*entity.Stage, error) {
	var stages []*entity.Stage

	if err := ssr.db.WithContext(ctx).
		Where("sendpost_id = ?", sendpostID).
		Where("parent_stage_id IS NULL").
		Order("next_stage_id asc").
		Find(&stages).Error; err != nil {
		return nil, err
	}
	logging.Debug("[Stage repo] GetSendpostStages", zap.Any("stages", stages))
	return stages, nil
}

// GetSubStages retrieves sub-stages associated with a given parent stage ID from the database.
// It uses the provided context for database operations to ensure proper request scoping and cancellation.
// Returns a slice of Stage entities or an error if the operation fails.
func (ssr *gormSendpostStageRepository) GetSubStages(ctx context.Context, parentStageID uint) ([]*entity.Stage, error) {
	var subStages []*entity.Stage

	if err := ssr.db.WithContext(ctx).
		Where("parent_stage_id = ?", parentStageID).
		Find(&subStages).Error; err != nil {
		return nil, err
	}
	logging.Debug("[Stage repo] GetSubStages", zap.Any("sub-stages", subStages))
	return subStages, nil
}

// DeleteStage removes a stage from the database using the provided stageID.
// It utilizes the context for managing request-scoped values and deadlines.
// Returns an error if the deletion fails.
func (ssr *gormSendpostStageRepository) DeleteStage(ctx context.Context, stageID uint) error {
	logging.Debug("[Stage repo] DeleteStage", zap.Uint("stageID", stageID))
	return ssr.db.WithContext(ctx).Delete(&entity.Stage{}, stageID).Error
}

// GetPreviousStage retrieves the previous stage from the database based on the given stageID.
// It uses the context for database operations and returns the previous stage entity or an error if not found.
// Logs the retrieved stage for debugging purposes.
func (ssr *gormSendpostStageRepository) GetPreviousStage(ctx context.Context, stageID uint) (*entity.Stage, error) {
	var stage *entity.Stage

	if err := ssr.db.WithContext(ctx).
		Where("next_stage_id = ?", stageID).
		First(&stage).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	logging.Debug("[Stage repo] GetPreviousStage", zap.Any("stage", stage))
	return stage, nil
}
