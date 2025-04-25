package services

import (
	"context"
	"crm-uplift-ii24-backend/internal/domain/entity"
	"crm-uplift-ii24-backend/internal/domain/repository"
	"crm-uplift-ii24-backend/internal/domain/value"
	"crm-uplift-ii24-backend/pkg/logging"
	"fmt"

	"go.uber.org/zap"
)

type SendpostService struct {
	sendpostRepo repository.SendpostRepository
	stageService *StageService
}

const (
	ErrorDeleteSendpost          string = "[SendpostService] error DeleteSendpost"
	ErrorUpdateSendpost          string = "[SendpostService] error UpdateSendpost"
	ErrorCopySendpost            string = "[SendpostService] error CopySendpost"
	ErrorAddSendpostParameter    string = "[SendpostService] error AddSendpostParameter"
	ErrorDeleteSendpostParameter string = "[SendpostService] error DeleteSendpostParameters"
)

func NewSendpostService(sendpostRepo repository.SendpostRepository, stageService *StageService) *SendpostService {
	return &SendpostService{sendpostRepo: sendpostRepo, stageService: stageService}
}

/*
CreateSendpost creates a new sendpost entity with the specified name, date, and optional first stage ID.
It saves the entity using the repository and returns the created sendpost or an error if the operation fails.

Parameters:
  - ctx: The context for managing request-scoped values, cancellation, and deadlines.
  - sendpostName: The name of the sendpost to be created.
  - sendpostDate: The date associated with the sendpost.
  - firstStageID: An optional pointer to the ID of the first stage.

Returns:
  - A pointer to the created Sendpost entity.
  - An error if the sendpost creation or saving fails.
*/
func (ss *SendpostService) CreateSendpost(
	ctx context.Context,
	sendpostName string,
	description *string,
	parameters *value.JSONB,
) (*entity.Sendpost, error) {

	logging.Debug("[SendpostService] CreateSendpost")

	sendpost := entity.Sendpost{}
	sendpost.Update(sendpostName, description, parameters)
	if err := ss.sendpostRepo.SaveSendpost(ctx, &sendpost); err != nil {
		return nil, fmt.Errorf("[SendpostService] error CreateSendpost: %s", err)
	}
	return &sendpost, nil
}

/*
DeleteSendpost deletes a sendpost and its associated stages from the repository.

This method first retrieves all stages associated with the given sendpostID
and attempts to delete each stage. If any stage deletion fails, an error is returned.
After successfully deleting all stages, it deletes the sendpost itself. If the
sendpost deletion fails, an error is returned.

Parameters:
  - ctx: The context for managing request-scoped values, cancellation signals, and deadlines.
  - sendpostID: The unique identifier of the sendpost to be deleted.

Returns:
  - error: An error if any operation fails, otherwise nil.
*/
func (ss *SendpostService) DeleteSendpost(ctx context.Context, sendpostID uint) error {

	logging.Debug("[SendpostService] DeleteSendpost")

	stages, err := ss.stageService.GetSendpostStages(ctx, sendpostID)

	if err != nil {
		return logging.WrapError(ErrorDeleteSendpost, err)
	}
	if err := ss.sendpostRepo.DeleteSendpost(ctx, sendpostID); err != nil {
		return logging.WrapError(ErrorDeleteSendpost, err)
	}
	for _, stage := range stages { // cascade delete stages
		if err := ss.stageService.DeleteStage(ctx, stage.ID); err != nil {
			return logging.WrapError(ErrorDeleteSendpost, err)
		}
	}
	return nil
}

// GetSendpost retrieves a Sendpost entity by its ID from the repository.
// It returns the Sendpost entity if found, or an error if the retrieval fails.
//
// Parameters:
//
//	ctx - The context for managing request-scoped values, cancellation signals, and deadlines.
//	sendpostID - The unique identifier of the Sendpost to be retrieved.
//
// Returns:
//
//	A pointer to the Sendpost entity if successful, or an error if the retrieval fails.
func (ss *SendpostService) GetSendpost(ctx context.Context, sendpostID uint) (*entity.Sendpost, error) {
	logging.Debug("[SendpostService] GetSendpost")
	sendpost, err := ss.sendpostRepo.GetSendpostByID(ctx, sendpostID)
	if err != nil {
		return nil, fmt.Errorf("[SendpostService] error GetSendpost: %s", err)
	}
	return sendpost, nil
}

// GetSendposts retrieves a list of Sendpost entities from the repository.
// It returns a slice of Sendpost pointers and an error if the retrieval fails.
// The context parameter is used for request-scoped values, cancellation, and deadlines.
func (ss *SendpostService) GetSendposts(ctx context.Context) ([]*entity.Sendpost, error) {
	logging.Debug("[SendpostService] GetSendposts")
	sendposts, err := ss.sendpostRepo.GetSendposts(ctx)
	if err != nil {
		return nil, fmt.Errorf("[SendpostService] error GetSendpost: %s", err)
	}
	return sendposts, nil
}

// UpdateSendpostState updates the state of a given Sendpost entity and persists the change.
// It logs the operation and any errors encountered during the update process.
//
// Parameters:
//
//	ctx - The context for managing request-scoped values, cancellation signals, and deadlines.
//	sendpost - The Sendpost entity whose state is to be updated.
//	state - The new state to be assigned to the Sendpost entity.
//
// Returns:
//
//	An error if the update operation fails, otherwise nil.
func (srs *SendpostService) UpdateSendpostState(ctx context.Context, id uint, state value.StateType) error {
	sendpost, err := srs.sendpostRepo.GetSendpostByID(ctx, id)
	if err != nil {
		return logging.WrapError("[SendpostRunnerService] error UpdateSendpostState", err)
	}
	logging.Debug("[SendpostRunnerService] UpdateSendpostState", zap.Any("sendspost", sendpost), zap.Any("state", state))
	sendpost.State = state
	if err := srs.sendpostRepo.SaveSendpost(ctx, sendpost); err != nil {
		return logging.WrapError("[SendpostRunnerService] error UpdateSendpostState", err)
	}
	return nil
}

// GetFirstStage retrieves the first stage of a sendpost by its ID.
// It logs the operation and returns the stage entity or an error if the retrieval fails.
//
// Parameters:
//
//	ctx - The context for controlling cancellation and deadlines.
//	sendpostID - The unique identifier of the sendpost.
//
// Returns:
//
//	*entity.Stage - The first stage entity of the sendpost.
//	error - An error if the operation fails, otherwise nil.
func (src *SendpostService) GetFirstStage(ctx context.Context, sendpostID uint) (*entity.Stage, error) {
	logging.Debug("[SendpostService] GetFirstStage")
	stage, err := src.sendpostRepo.GetFirstStage(ctx, sendpostID)
	if err != nil {
		return nil, fmt.Errorf("[SendpostService] error GetFirstStage: %s", err)
	}
	return stage, nil
}

/*
CopySendpost duplicates an existing sendpost identified by sendpostID, assigning it a new name and optional description.
It retrieves the original sendpost, creates a copy with the new attributes, and saves the new sendpost to the repository.
Additionally, it copies associated stages from the original sendpost to the new one.

Parameters:
  - ctx: The context for managing request-scoped values, cancellation signals, and deadlines.
  - sendpostID: The unique identifier of the sendpost to be copied.
  - NewName: The new name for the copied sendpost.
  - NewDescription: An optional new description for the copied sendpost.

Returns:
  - A pointer to the newly created sendpost entity.
  - An error if the operation fails at any step.
*/
func (src *SendpostService) CopySendpost(ctx context.Context, sendpostID uint, NewName string, NewDescription *string, NewParams *value.JSONB) (*entity.Sendpost, error) {
	logging.Debug("[SendpostService] CopySendpost")
	sendpost, err := src.sendpostRepo.GetSendpostByID(ctx, sendpostID)
	if err != nil {
		return nil, logging.WrapError(ErrorCopySendpost, err)
	}
	newSendpost := sendpost.Copy(NewName, NewDescription)
	if err := src.sendpostRepo.SaveSendpost(ctx, newSendpost); err != nil {
		return nil, logging.WrapError(ErrorCopySendpost, err)
	}
	for k, v := range *NewParams {
		if err := src.AddUpdateSendpostParameter(ctx, newSendpost.ID, k, v); err != nil {
			return nil, logging.WrapError(ErrorCopySendpost, err)
		}
	}
	logging.Debug("[SendpostService] CopySendpost", zap.Any("newSendpost", newSendpost))

	if err := src.stageService.CopyStages(ctx, sendpostID, newSendpost.ID); err != nil {
		return nil, logging.WrapError(ErrorCopySendpost, err)
	}
	return newSendpost, nil
}

// AddUpdateSendpostParameter updates or adds a parameter to a sendpost entity.
// It retrieves the sendpost by ID, updates its parameters, and saves the changes.
// Logs the operation and wraps any errors encountered during the process.
//
// Parameters:
//
//	ctx - The context for managing request-scoped values and cancellation.
//	sendpostID - The unique identifier of the sendpost to update.
//	key - The parameter key to add or update.
//	value - The value to associate with the specified key.
//
// Returns:
//
//	An error if the operation fails, otherwise nil.
func (src *SendpostService) AddUpdateSendpostParameter(ctx context.Context, sendpostID uint, key string, value any) error {
	logging.Debug("[SendpostService] AddUpdateSendpostParameter")
	sendpost, err := src.sendpostRepo.GetSendpostByID(ctx, sendpostID)
	if err != nil {
		return logging.WrapError(ErrorAddSendpostParameter, err)
	}
	parameters, err := src.GetSendpostParameters(ctx, sendpostID)
	if err != nil {
		return logging.WrapError(ErrorAddSendpostParameter, err)
	}
	(*parameters)[key] = value
	sendpost.GlobalParameters = parameters
	logging.Debug("[SendpostService] AddUpdateSendpostParameter", zap.Any("sendpost", sendpost))
	if err := src.sendpostRepo.SaveSendpost(ctx, sendpost); err != nil {
		return logging.WrapError(ErrorAddSendpostParameter, err)
	}
	return nil
}

// DeleteSendpostParameter removes a specific parameter from a sendpost's global parameters.
// It retrieves the sendpost by its ID, deletes the specified key from its parameters,
// and saves the updated sendpost back to the repository. If any error occurs during
// these operations, it wraps and returns the error.
func (src *SendpostService) DeleteSendpostParameter(ctx context.Context, sendpostID uint, key string) error {
	sendpost, err := src.sendpostRepo.GetSendpostByID(ctx, sendpostID)
	if err != nil {
		return logging.WrapError(ErrorDeleteSendpostParameter, err)
	}
	parameters, err := src.sendpostRepo.GetSendpostParameters(ctx, sendpostID)
	if err != nil {
		return logging.WrapError(ErrorDeleteSendpostParameter, err)
	}
	delete(*parameters, key)
	sendpost.GlobalParameters = parameters
	logging.Debug("[SendpostService] DeleteSendpostParameter", zap.Any("sendpost", sendpost))
	if err := src.sendpostRepo.SaveSendpost(ctx, sendpost); err != nil {
		return logging.WrapError(ErrorDeleteSendpostParameter, err)
	}
	return nil
}

// GetSendpostParamters retrieves the parameters for a given sendpost ID from the repository.
// It returns the parameters as a JSONB value or an error if the retrieval fails.
// The context is used for request-scoped values, cancellation, and deadlines.
func (src *SendpostService) GetSendpostParameters(ctx context.Context, sendpostID uint) (*value.JSONB, error) {
	logging.Debug("[Sendpost Service] GetSendpostParameters")
	parameters, err := src.sendpostRepo.GetSendpostParameters(ctx, sendpostID)
	if err != nil {
		return nil, logging.WrapError("[Sendpost Service] Error GetSendpostParameters", err)
	}
	if parameters == nil {
		logging.Info("[Sendpost Service] GetSendpostParameters: parameters is nil")
		parameters = &value.JSONB{}
	}
	return parameters, nil
}
