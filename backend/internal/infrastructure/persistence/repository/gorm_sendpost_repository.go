package repository

import (
	"context"
	"crm-uplift-ii24-backend/internal/domain/entity"
	"crm-uplift-ii24-backend/internal/domain/repository"
	"crm-uplift-ii24-backend/internal/domain/value"
	"crm-uplift-ii24-backend/pkg/logging"
	"fmt"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type gormSendpostRepository struct {
	db *gorm.DB
}

// NewGormSendpostRepository creates a new instance of gormSendpostRepository
// using the provided gorm.DB connection. It returns an implementation of
// the SendpostRepository interface.
func NewGormSendpostRepository(db *gorm.DB) repository.SendpostRepository {
	return &gormSendpostRepository{db: db}
}

// SaveSendpost persists a Sendpost entity to the database using the provided context.
// It returns an error if the operation fails.
func (sr *gormSendpostRepository) SaveSendpost(ctx context.Context, sendpost *entity.Sendpost) error {
	logging.Debug("[Sendpost repo] SaveSendpost", zap.Any("sendpost", sendpost))
	return sr.db.WithContext(ctx).Save(sendpost).Error
}

// DeleteSendpost deletes a Sendpost entity from the database using its ID.
// It utilizes the provided context for request scoping and cancellation.
// Returns an error if the deletion fails.
func (sr *gormSendpostRepository) DeleteSendpost(ctx context.Context, sendpostID uint) error {
	logging.Debug("[Sendpost repo] DeleteSendpost", zap.Uint("sendpost_id", sendpostID))
	return sr.db.WithContext(ctx).Delete(&entity.Sendpost{}, sendpostID).Error
}

// GetSendpostByID retrieves a Sendpost entity by its ID from the database.
// It uses the provided context for database operations and preloads the "FirstStage" relation.
// Returns the Sendpost entity if found, or an error if not found or if a database error occurs.
func (sr *gormSendpostRepository) GetSendpostByID(ctx context.Context, sendpostID uint) (*entity.Sendpost, error) {
	var sendpost *entity.Sendpost
	if err := sr.db.WithContext(ctx).
		First(&sendpost, sendpostID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("sendpost not found: %s", err)
		} else {
			return nil, err
		}
	}
	logging.Debug("[Senpost repo] GetSendposts", zap.Uint("sendpost_id", sendpostID), zap.Any("sendpost", sendpost))
	return sendpost, nil
}

// GetSendposts retrieves all Sendpost entities from the database that are not templates.
// It uses the provided context for database operations and returns a slice of Sendpost
// pointers or an error if the operation fails. Debug logs are generated with the retrieved
// sendposts.
func (sr *gormSendpostRepository) GetSendposts(ctx context.Context) ([]*entity.Sendpost, error) {
	var sendposts []*entity.Sendpost

	if err := sr.db.WithContext(ctx).
		Find(&sendposts).Error; err != nil {
		return nil, err
	}
	logging.Debug("[Senpost repo] GetSendposts", zap.Any("sendposts", sendposts))
	return sendposts, nil
}

// GetFirstStage retrieves the first stage associated with a given sendpost ID.
// It uses the provided context for database operations and returns the first stage
// entity or an error if the operation fails.
func (sr *gormSendpostRepository) GetFirstStage(ctx context.Context, sendpostID uint) (*entity.Stage, error) {
	var sendpost entity.Sendpost
	if err := sr.db.WithContext(ctx).
		Preload("FirstStage").
		Where("id = ?", sendpostID).
		First(&sendpost).Error; err != nil {
		return nil, err
	}
	logging.Debug("[Senpost repo] GetFirstStage", zap.Any("first_stage", sendpost.FirstStage))
	return sendpost.FirstStage, nil
}

// GetSendpostParameters retrieves the global parameters of a Sendpost entity
// from the database using the provided sendpostID.
// It executes the query within the given context and returns the global parameters
// as a JSONB entity or an error if the operation fails.
//
// Parameters:
//
//	ctx - The context for managing request-scoped values, cancellation, and deadlines.
//	sendpostID - The unique identifier of the Sendpost entity to retrieve the parameters for.
//
// Returns:
//
//	*entity.JSONB - A pointer to the global parameters of the Sendpost entity.
//	error - An error if the operation fails, otherwise nil.
func (sr *gormSendpostRepository) GetSendpostParameters(ctx context.Context, sendpostID uint) (*value.JSONB, error) {
	var sendpost entity.Sendpost
	if err := sr.db.WithContext(ctx).
		Where("id = ?", sendpostID).
		First(&sendpost).Error; err != nil {
		return nil, err
	}
	return sendpost.GlobalParameters, nil
}

// UpdateSendpostParameters updates the global parameters of a Sendpost entity
// in the database using the provided sendpostID and parameters map.
// It executes the update operation within the given context and returns an error
// if the operation fails.
//
// Parameters:
//
//	ctx - The context for managing request-scoped values, cancellation, and deadlines.
//	sendpostID - The unique identifier of the Sendpost entity to update.
//	parameters - A pointer to a map containing the new global parameters to set.
//
// Returns:
//
//	An error if the update operation fails, otherwise nil.
func (sr *gormSendpostRepository) UpdateSendpostParameters(ctx context.Context, sendpostID uint, parameters *map[string]interface{}) error {
	return sr.db.WithContext(ctx).
		Model(&entity.Sendpost{}).
		Where("id = ?", sendpostID).
		Update("global_parameters", parameters).Error
}
