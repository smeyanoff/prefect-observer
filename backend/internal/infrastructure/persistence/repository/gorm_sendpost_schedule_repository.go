package repository

import (
	"context"
	"crm-uplift-ii24-backend/internal/domain/entity"
	"crm-uplift-ii24-backend/internal/domain/repository"
	"errors"

	"gorm.io/gorm"
)

type gormSendpostScheduleRepository struct {
	db *gorm.DB
}

func NewGormSendpostScheduleRepository(db *gorm.DB) repository.SendpostScheduleRepository {
	return &gormSendpostScheduleRepository{db: db}
}

// SaveSchedule saves a SendpostSchedule entity to the database using the provided context.
// It returns an error if the operation fails.
func (r *gormSendpostScheduleRepository) SaveSchedule(ctx context.Context, schedule *entity.SendpostSchedule) error {
	return r.db.WithContext(ctx).Save(schedule).Error
}

// DeleteSchedule removes a SendpostSchedule entity from the database using the provided scheduleID.
// It executes the operation within the given context and returns an error if the deletion fails.
//
// Parameters:
//
//	ctx - The context to execute the database operation within.
//	scheduleID - The unique identifier of the SendpostSchedule to be deleted.
//
// Returns:
//
//	An error if the deletion operation fails, otherwise nil.
func (r *gormSendpostScheduleRepository) DeleteSchedule(ctx context.Context, scheduleID uint) error {
	return r.db.WithContext(ctx).Delete(&entity.SendpostSchedule{}, scheduleID).Error
}

// GetScheduleByID retrieves a SendpostSchedule entity by its ID from the database.
// It uses the provided context for request scoping and cancellation.
// Returns the SendpostSchedule entity if found, or an error if not found or if a database error occurs.
func (r *gormSendpostScheduleRepository) GetScheduleByID(ctx context.Context, scheduleID uint) (*entity.SendpostSchedule, error) {
	var sendpost entity.SendpostSchedule

	if err := r.db.WithContext(ctx).
		First(&sendpost, scheduleID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("schedule not found")
		} else {
			return nil, err
		}
	}
	return &sendpost, nil
}

// GetScheduleBySendpostID retrieves a SendpostSchedule entity from the database
// using the provided sendpostID. It returns the SendpostSchedule if found,
// or an error if the schedule is not found or if any other database error occurs.
//
// Parameters:
//
//	ctx - The context for managing request-scoped values, cancellation, and deadlines.
//	sendpostID - The unique identifier for the SendpostSchedule to be retrieved.
//
// Returns:
//
//	*entity.SendpostSchedule - The retrieved SendpostSchedule entity.
//	error - An error if the schedule is not found or if a database error occurs.
func (r *gormSendpostScheduleRepository) GetScheduleBySendpostID(ctx context.Context, sendpostID uint) (*entity.SendpostSchedule, error) {
	var sendpost entity.SendpostSchedule

	if err := r.db.WithContext(ctx).
		Where("sendpost_id = ?", sendpostID).
		First(&sendpost).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("schedule not found")
		} else {
			return nil, err
		}
	}
	return &sendpost, nil
}
