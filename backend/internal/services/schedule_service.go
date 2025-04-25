package services

import (
	"context"
	"crm-uplift-ii24-backend/internal/domain/entity"
	"crm-uplift-ii24-backend/internal/domain/repository"
	"fmt"
	"time"
)

type ScheduleService struct {
	repo repository.SendpostScheduleRepository
}

func NewScheduleService(repo repository.SendpostScheduleRepository) *ScheduleService {
	return &ScheduleService{repo: repo}
}

// CreateSchedule creates a new sendpost schedule with the specified sendpost ID and planned time.
// It saves the schedule using the repository and returns the created schedule or an error if saving fails.
//
// Parameters:
//
//	ctx - The context for managing request-scoped values, cancellation, and deadlines.
//	sendopstID - The unique identifier for the sendpost.
//	plannedAt - The time when the sendpost is planned to be sent.
//
// Returns:
//
//	*entity.SendpostSchedule - The created sendpost schedule.
//	error - An error if the schedule could not be saved.
func (ss *ScheduleService) CreateSchedule(
	ctx context.Context,
	sendopstID uint,
	plannedAt time.Time,

) (*entity.SendpostSchedule, error) {
	schedule := entity.SendpostSchedule{
		SendpostID: sendopstID,
		PlannedAt:  plannedAt,
	}
	if err := ss.repo.SaveSchedule(ctx, &schedule); err != nil {
		return nil, fmt.Errorf("[ScheduleService] error create schedule: %s", err)
	}
	return &schedule, nil
}

// UpdateSchedule updates the given SendpostSchedule in the repository.
// It takes a context and a pointer to a SendpostSchedule entity as parameters.
// Returns an error if the update operation fails.
func (ss *ScheduleService) UpdateSchedule(ctx context.Context, schedule *entity.SendpostSchedule) error {
	err := ss.repo.SaveSchedule(ctx, schedule)
	if err != nil {
		return fmt.Errorf("[ScheduleService] error updating schedule: %s", err)
	}
	return nil
}

// UpdateScheduleStartedAt updates the StartedAt timestamp of a SendpostSchedule
// entity to the current time and saves the updated schedule to the repository.
// Returns an error if the save operation fails.
//
// Parameters:
//   - ctx: The context for managing request-scoped values, cancellation signals,
//     and deadlines.
//   - schedule: A pointer to the SendpostSchedule entity to be updated.
//
// Returns:
//   - error: An error if the schedule could not be saved, otherwise nil.
func (ss *ScheduleService) UpdateScheduleStartedAt(ctx context.Context, schedule *entity.SendpostSchedule) error {
	currentTime := time.Now()
	schedule.StartedAt = &currentTime
	err := ss.repo.SaveSchedule(ctx, schedule)
	if err != nil {
		return fmt.Errorf("[ScheduleService] error updating schedule StartedAt: %s", err)
	}
	return nil
}

// UpdateScheduleCompletedAt updates the CompletedAt field of a SendpostSchedule
// entity to the current time and saves the updated schedule to the repository.
// Returns an error if the save operation fails.
//
// Parameters:
//
//	ctx - The context for managing request-scoped values, cancellation signals, and deadlines.
//	schedule - A pointer to the SendpostSchedule entity to be updated.
//
// Returns:
//
//	An error if the schedule could not be saved, otherwise nil.
func (ss *ScheduleService) UpdateScheduleCompletedAt(ctx context.Context, schedule *entity.SendpostSchedule) error {
	currentTime := time.Now()
	schedule.CompletedAt = &currentTime
	err := ss.repo.SaveSchedule(ctx, schedule)
	if err != nil {
		return fmt.Errorf("[ScheduleService] error updating schedule CompletedAt: %s", err)
	}
	return nil
}
