package services

import (
	"context"
	"crm-uplift-ii24-backend/internal/domain/entity"
	"crm-uplift-ii24-backend/internal/domain/repository"
	"crm-uplift-ii24-backend/internal/domain/value"
	"crm-uplift-ii24-backend/pkg/logging"
	"errors"
	"fmt"

	"go.uber.org/zap"
)

const (
	ErrorDeleteStage      string = "[StageService] error DeleteStage"
	ErrorBlockUnblock     string = "[StageService] error BlockUnblockStage"
	ErrorGetSubStages     string = "[StageService] error GetSubStages"
	ErrorCopyStages       string = "[StageService] error copyStages"
	ErrorCopySubStages    string = "[StageService] error copySubStages"
	ErrorUpdateParameters string = "[StageService] error UpdateParameters"
)

type StageService struct {
	stageRepo    repository.StageRepository
	sendpostRepo repository.SendpostRepository
}

func NewStageService(stageRepo repository.StageRepository, sendpostRepo repository.SendpostRepository) *StageService {
	return &StageService{stageRepo: stageRepo, sendpostRepo: sendpostRepo}
}

// SaveStage updates the given stage or creates new in the stageRepository.
// It takes a context and a pointer to an entity.Stage as parameters.
// Returns an error if the update operation fails.
func (s *StageService) saveStage(ctx context.Context, stage *entity.Stage) error {
	err := s.stageRepo.SaveStage(ctx, stage)
	if err != nil {
		return fmt.Errorf("[StageService] error saveStagee: %s", err)
	}
	return nil
}

// UpdateNextStageID updates the next stage ID for a given stage.
// It retrieves the stage using the provided stage and updates its nextStageID.
// Returns an error if update fails.
func (s *StageService) updateNextStageID(ctx context.Context, stage *entity.Stage, nextStageID *uint) error {
	logging.Debug("[StageService] updateNextStageID")

	if err := stage.UpdateNextStageID(nextStageID); err != nil {
		return fmt.Errorf("[StageService] error updateNextStageID: %s", err)
	}

	if err := s.saveStage(ctx, stage); err != nil {
		return err
	}

	return nil
}

// Handles inserting a stage after a given previous stage
func (s *StageService) handlePreviousStage(ctx context.Context, stage *entity.Stage, previousStageID uint) error {
	logging.Debug("[StageService] updateNextStageID")

	previousStage, err := s.GetStage(ctx, previousStageID)
	if err != nil {
		return err
	}

	if previousStage.NextStageID != nil {
		if err := s.updateNextStageID(ctx, stage, previousStage.NextStageID); err != nil {
			return err
		}
	}

	return s.updateNextStageID(ctx, previousStage, &stage.ID)
}

// Handles inserting a stage as the first stage in the sendpost sequence
func (s *StageService) handleFirstStage(ctx context.Context, sendpostID uint, stage *entity.Stage) error {

	logging.Debug("[Stage Service] handleFirstStage")

	firstStage, err := s.sendpostRepo.GetFirstStage(ctx, sendpostID)
	if err != nil {
		return err
	}

	// update sendpost's first stage
	sendpost, err := s.sendpostRepo.GetSendpostByID(ctx, sendpostID)
	if err != nil {
		return err
	}
	sendpost.FirstStageID = &stage.ID
	if err := s.sendpostRepo.SaveSendpost(ctx, sendpost); err != nil {
		return err
	}

	if firstStage != nil {
		return s.updateNextStageID(ctx, stage, &firstStage.ID)
	}
	return nil
}

// AddStage adds a new stage to the system with the specified parameters.
// It saves the stage to the stageRepository and updates the next stage ID if
// a previous stage is provided. Returns the created stage or an error.
//
// Parameters:
//
//	ctx - The context for managing request-scoped values.
//	sendpostID - The ID associated with the sendpost.
//	stageType - The type of the stage being added.
//	parentStageID - The ID of the parent stage, if applicable.
//	deploymentID - The ID of the deployment, if applicable.
//	stageParameters - Additional parameters for the stage.
//	previousStageID - The ID of the previous stage, if applicable.
//
// Returns:
//
//	*entity.Stage - The newly created stage entity.
//	error - An error if the stage could not be added.
func (s *StageService) AddStage(
	ctx context.Context,
	stage *entity.Stage,
	previousStageID *uint,
) error {

	logging.Debug("[StageService] AddStage", zap.Any("stage", stage), zap.Any("previousStageID", previousStageID))

	if err := s.saveStage(ctx, stage); err != nil {
		return err
	}

	if previousStageID != nil {
		if err := s.handlePreviousStage(ctx, stage, *previousStageID); err != nil {
			return err
		}
		return nil
	}

	if err := s.handleFirstStage(ctx, stage.SendpostID, stage); err != nil {
		return err
	}

	return nil
}

// AddSubStage adds a sub-stage to a parent stage if the parent stage is of type ParallelStage.
// It retrieves the parent stage using the provided parentStageID and checks its type.
// If the parent stage is not of type ParallelStage, an error is returned.
// Otherwise, it creates a new sub-stage with the specified parameters and saves it.
// Returns the newly created sub-stage or an error if the operation fails.
func (s *StageService) AddSubStage(
	ctx context.Context,
	parentStageID uint,
	stage *entity.Stage,
) error {

	parentStage, err := s.GetStage(ctx, parentStageID)
	if err != nil {
		return fmt.Errorf("[StageService] error AddSubStage: %s", err)
	}
	if parentStage.Type != value.ParallelStage {
		return errors.New("[StageService] error AddSubStage: try to add sub-stage to a non parallel stage")
	}

	stage.ParentStageID = &parentStageID
	if err := s.saveStage(ctx, stage); err != nil {
		return err
	}
	return nil
}

// cascadeDeleteSubStages deletes all sub-stages of a given stage in parallel.
// It retrieves the sub-stages of the specified stage and attempts to delete each one.
// If an error occurs during the deletion of any sub-stage, it returns an error with a descriptive message.
// Parameters:
//
//	ctx - The context for managing request-scoped values, cancellation, and deadlines.
//	stage - The stage entity for which sub-stages need to be deleted.
//
// Returns:
//
//	An error if any sub-stage deletion fails, otherwise nil.
func (s *StageService) cascadeDeleteSubStages(ctx context.Context, stageId uint) error {
	stages, err := s.GetSubStages(ctx, stageId)
	if err != nil {
		return err
	}
	for _, stage := range stages {
		if err := s.DeleteStage(ctx, stage.ID); err != nil {
			return fmt.Errorf("[StageService] error cascadeDeleteSubStages: %s", err)
		}
	}
	return nil
}

// DeleteStage removes a stage identified by stageID from the stageRepository.
// It returns an error if the deletion fails.
func (s *StageService) DeleteStage(ctx context.Context, stageID uint) error {
	stage, err := s.GetStage(ctx, stageID)
	if err != nil {
		return logging.WrapError(ErrorDeleteStage, err)
	}
	// update previous stage's next stage ID
	previousStage, err := s.stageRepo.GetPreviousStage(ctx, stageID)
	if err != nil {
		return logging.WrapError(ErrorDeleteStage, err)
	}
	if previousStage != nil {
		// change previous stage's next stage ID to stage's next stage ID
		if err := s.updateNextStageID(ctx, previousStage, stage.NextStageID); err != nil {
			return logging.WrapError(ErrorDeleteStage, err)
		}
	}
	if stage.Type == value.ParallelStage { //cascade delete for parallel stages
		if err := s.cascadeDeleteSubStages(ctx, stage.ID); err != nil {
			return logging.WrapError(ErrorDeleteStage, err)
		}
	}
	if err := s.stageRepo.DeleteStage(ctx, stageID); err != nil {
		return logging.WrapError(ErrorDeleteStage, err)
	}
	return nil
}

// GetStage retrieves a stage entity by its ID from the stageRepository.
// It returns the stage entity if found, or an error if the retrieval fails.
//
// Parameters:
//
//	ctx - The context for managing request-scoped values and cancellation.
//	stageID - The unique identifier of the stage to retrieve.
//
// Returns:
//
//	*entity.Stage - The retrieved stage entity.
//	error - An error if the stage could not be retrieved.
func (s *StageService) GetStage(ctx context.Context, stageID uint) (*entity.Stage, error) {
	stage, err := s.stageRepo.GetStageByID(ctx, stageID)
	if err != nil {
		return nil, fmt.Errorf("[StageService] error GetStage: %s", err)
	}
	return stage, nil
}

// UpdateStageState updates the state of a stage identified by stageID.
// It retrieves the stage, updates its state, and saves the changes.
// Returns an error if the stage cannot be retrieved or saved.
func (s *StageService) UpdateStageState(ctx context.Context, stage *entity.Stage, state value.StateType) error {
	stage.UpdateState(state)
	if err := s.saveStage(ctx, stage); err != nil {
		return fmt.Errorf("[StageService] error UpdateStageState: %s", err)
	}
	return nil
}

// GetSubStages retrieves the sub-stages of a given stage by its ID.
// It first checks if the stage is of type ParallelStage, returning an error if not.
// If the stage is valid, it fetches the sub-stages from the repository.
// Returns a slice of Stage entities or an error if the operation fails.
func (s *StageService) GetSubStages(ctx context.Context, stageID uint) ([]*entity.Stage, error) {
	stage, err := s.GetStage(ctx, stageID)
	if err != nil {
		return nil, logging.WrapError(ErrorGetSubStages, err)
	}
	if stage.Type != value.ParallelStage {
		return nil, fmt.Errorf("%s: %s", ErrorGetSubStages, "stage is not a parallel stage")
	}
	stages, err := s.stageRepo.GetSubStages(ctx, stage.ID)
	if err != nil {
		return nil, logging.WrapError(ErrorGetSubStages, err)
	}
	return stages, nil
}

// GetSendpostStages retrieves the stages associated with a given sendpost ID.
// It logs the operation and returns a slice of Stage entities or an error if the retrieval fails.
//
// Parameters:
//
//	ctx - The context for managing request-scoped values, cancellation signals, and deadlines.
//	sendpostID - The unique identifier of the sendpost whose stages are to be retrieved.
//
// Returns:
//
//	A slice of pointers to Stage entities if successful, or an error if the operation fails.
func (s *StageService) GetSendpostStages(ctx context.Context, sendpostID uint) ([]*entity.Stage, error) {
	var stages []*entity.Stage
	firstStage, err := s.sendpostRepo.GetFirstStage(ctx, sendpostID)
	if err != nil {
		return nil, fmt.Errorf("[StageService] error GetSendpostStages: %s", err)
	}
	if firstStage == nil {
		return nil, nil
	}
	stages = append(stages, firstStage)
	nextStageID := firstStage.NextStageID
	for nextStageID != nil {
		stage, err := s.GetStage(ctx, *nextStageID)
		if err != nil {
			return nil, fmt.Errorf("[StageService] error GetSendpostStages: %s", err)
		}
		stages = append(stages, stage)
		nextStageID = stage.NextStageID
	}
	return stages, nil
}

// BlockUnblockStage toggles the blocked status of a stage identified by stageID.
// It retrieves the stage, inverts its blocked status, and saves the updated stage.
// Returns an error if any operation fails during the process.
func (s *StageService) BlockUnblockStage(ctx context.Context, stageID uint) error {
	stage, err := s.GetStage(ctx, stageID)
	if err != nil {
		return logging.WrapError(ErrorBlockUnblock, err)
	}
	stage.IsBlocked = !stage.IsBlocked
	if err := s.stageRepo.SaveStage(ctx, stage); err != nil {
		return logging.WrapError(ErrorBlockUnblock, err)
	}
	return nil
}

// CopyStages duplicates the stages associated with a given sendpost ID to a new sendpost ID.
// It retrieves the stages of the original sendpost, copies each stage to the new sendpost,
// and maintains the order of stages by linking them appropriately.
// If any error occurs during the process, it wraps and returns the error.
func (src *StageService) CopyStages(ctx context.Context, sendpostID uint, newSendpostId uint) error {
	stages, err := src.GetSendpostStages(ctx, sendpostID)
	if err != nil {
		return logging.WrapError(ErrorCopyStages, err)
	}
	var previousStageId *uint
	previousStageId = nil
	for _, stage := range stages {
		newStage := stage.Copy(newSendpostId)
		if err := src.AddStage(ctx, newStage, previousStageId); err != nil {
			return logging.WrapError(ErrorCopyStages, err)
		}
		if err := src.copySubStages(ctx, stage.ID, newStage); err != nil {
			return logging.WrapError(ErrorCopyStages, err)
		}
		previousStageId = &newStage.ID
	}
	return nil
}

/*
copySubStages copies sub-stages from a parent stage to a new parent stage if the new parent stage is parallel.
It retrieves sub-stages of the given parent stage ID and duplicates them under the new parent stage.
If any error occurs during retrieval or addition of sub-stages, it wraps and returns the error.

Parameters:
  - ctx: The context for managing request-scoped values, cancellation, and timeouts.
  - parentStageID: The ID of the parent stage from which sub-stages are copied.
  - newParentStage: The new parent stage to which sub-stages are copied.

Returns:
  - error: An error if the operation fails, otherwise nil.
*/
func (s *StageService) copySubStages(ctx context.Context, parentStageID uint, newParentStage *entity.Stage) error {
	if !newParentStage.IsParallel() {
		return nil
	}
	subStages, err := s.GetSubStages(ctx, parentStageID)
	if err != nil {
		return logging.WrapError(ErrorCopyStages, err)
	}
	for _, subStage := range subStages {
		newSubStage := subStage.Copy(newParentStage.SendpostID)
		if err := s.AddSubStage(ctx, newParentStage.ID, newSubStage); err != nil {
			return logging.WrapError(ErrorCopySendpost, err)
		}
	}
	return nil
}

// UpdateParameters updates the parameters of a stage identified by stageID.
// It retrieves the stage, updates its parameters, and saves the changes.
// Returns an error if the stage cannot be retrieved or saved.
func (s *StageService) UpdateParameters(ctx context.Context, stageID uint, parameters value.JSONB) error {
	stage, err := s.GetStage(ctx, stageID)
	if err != nil {
		return logging.WrapError(ErrorUpdateParameters, err)
	}
	stage.StageParameters = &parameters
	if err := s.stageRepo.SaveStage(ctx, stage); err != nil {
		return logging.WrapError(ErrorUpdateParameters, err)
	}
	return nil
}
