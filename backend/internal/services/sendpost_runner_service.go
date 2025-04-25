package services

import (
	"context"
	"crm-uplift-ii24-backend/internal/domain/entity"
	"crm-uplift-ii24-backend/internal/domain/value"
	"crm-uplift-ii24-backend/pkg/logging"

	"go.uber.org/zap"
)

const (
	RunningStageError                               string = "[SendpostRunnerService] error running sendpost"
	ProcessError                                    string = "[SendpostRunnerService] error processing stage"
	ErrorNotifyRunSendpost                          string = "[SendpostRunnerService] error notifying run sendpost"
	ReplaceStageParametersWithSendpostParametersErr string = "[SendpostRunnerService] error replacing stage parameters with sendpost parameters"
)

type SendpostRunnerService struct {
	sendpostService            *SendpostService
	stageService               *StageService
	senpostNotificationService *SenpostRunNotificationService
	stageRunnerFactory         entity.StageRunnerFactory
}

func NewSendpostRunService(sendpostService *SendpostService, stageService *StageService, senpostNotificationService *SenpostRunNotificationService, stageRunnerFactory entity.StageRunnerFactory) *SendpostRunnerService {
	return &SendpostRunnerService{
		stageRunnerFactory:         stageRunnerFactory,
		stageService:               stageService,
		sendpostService:            sendpostService,
		senpostNotificationService: senpostNotificationService,
	}
}

// Start initiates the sendpost process for a given sendpost ID.
// It retrieves the sendpost and its associated stages from the repository,
// then sequentially processes each stage using the appropriate runner.
// If any error occurs during retrieval or processing, it returns an error.
//
// Parameters:
//
//	ctx - The context for managing request-scoped values and cancellation.
//	sendpostID - The unique identifier of the sendpost to be processed.
//
// Returns:
//
//	An error if any step in the process fails, otherwise nil.
func (srs *SendpostRunnerService) Start(ctx context.Context, sendpostID uint) {
	go srs.runStages(ctx, sendpostID)
}

func (srs *SendpostRunnerService) runStages(ctx context.Context, sendpostID uint) {
	srs.senpostNotificationService.AddRunSendpostToNotify(sendpostID)
	defer srs.senpostNotificationService.RemoveRunSendpostToNotify(sendpostID)

	stage, err := srs.sendpostService.GetFirstStage(ctx, sendpostID)
	if err != nil {
		srs.notifyRunErr(ctx, sendpostID, err)
		return
	}
	if stage == nil {
		srs.notifyRunErr(ctx, sendpostID, err)
		return
	}

	if err := srs.sendpostService.UpdateSendpostState(ctx, sendpostID, value.Running); err != nil {
		srs.notifyRunErr(ctx, sendpostID, err)
		return
	}

	if err := srs.processStage(ctx, stage); err != nil {
		srs.notifyRunErr(ctx, sendpostID, err)
		return
	}

	for stage.NextStageID != nil {
		nextStage, err := srs.stageService.GetStage(ctx, *stage.NextStageID)
		if err != nil {
			srs.notifyRunErr(ctx, sendpostID, err)
			return
		}
		stage = nextStage
		if err := srs.processStage(ctx, stage); err != nil {
			srs.notifyRunErr(ctx, sendpostID, err)
			return
		}
	}
	if err := srs.sendpostService.UpdateSendpostState(ctx, sendpostID, value.Completed); err != nil {
		srs.notifyRunErr(ctx, sendpostID, err)
		return
	}
	if err := srs.senpostNotificationService.NotifyRunSendpost(sendpostID, value.Completed); err != nil {
		logging.Warn(ErrorNotifyRunSendpost)
	}
}

// processStage processes a given stage by replacing its parameters, creating a runner,
// starting the runner, and notifying the sendpost service of updates. It handles errors
// by wrapping them with a process error and logs warnings if notification fails.
//
// Parameters:
//
//	ctx - The context for managing request-scoped values, cancelation signals, and deadlines.
//	stage - The stage entity to be processed.
//
// Returns:
//
//	An error if any step in the process fails, otherwise nil.
func (srs *SendpostRunnerService) processStage(ctx context.Context, stage *entity.Stage) error {
	logging.Debug("[SendpostRunnerService] processStage", zap.Any("stage", stage))
	if stage.Type == value.ParallelStage {
		stages, err := srs.stageService.GetSubStages(ctx, stage.ID)
		if err != nil {
			return logging.WrapError(ProcessError, err)
		}
		for _, subStage := range stages {
			if err := srs.ReplaceStageParametersWithSendpostParameters(ctx, subStage); err != nil {
				return logging.WrapError(ProcessError, err)
			}
		}
	} else {
		if err := srs.ReplaceStageParametersWithSendpostParameters(ctx, stage); err != nil {
			return logging.WrapError(ProcessError, err)
		}
	}
	runner := srs.stageRunnerFactory.CreateRunner(stage.Type)
	if err := runner.Start(ctx, stage); err != nil {
		return logging.WrapError(ProcessError, err)
	}
	if err := srs.senpostNotificationService.NotifyRunSendpost(stage.SendpostID, value.Updated); err != nil {
		logging.Warn(ErrorNotifyRunSendpost)
	}

	if err := runner.CheckState(ctx, stage); err != nil {
		return logging.WrapError(ProcessError, err)
	}
	if err := srs.senpostNotificationService.NotifyRunSendpost(stage.SendpostID, value.Updated); err != nil {
		logging.Warn(ErrorNotifyRunSendpost)
	}

	return nil
}

// notifyRunErr handles errors during the execution of a sendpost operation.
// It updates the sendpost state to 'Failed' and sends a notification about the failure.
// Additionally, it logs the error using the organization's logging package.
//
// Parameters:
//
//	ctx - The context for managing request-scoped values.
//	sendpostID - The unique identifier of the sendpost operation.
//	err - The error encountered during the sendpost execution.
func (srs *SendpostRunnerService) notifyRunErr(ctx context.Context, sendpostID uint, err error) {
	srs.sendpostService.UpdateSendpostState(ctx, sendpostID, value.Failed)
	srs.senpostNotificationService.NotifyRunSendpost(sendpostID, value.Failed)
	logging.Error(RunningStageError, zap.Error(err))
}

// replaceStageParametersWithSendpostParameters updates the stage parameters with
// the corresponding sendpost parameters. It retrieves global parameters using the
// sendpost ID and replaces matching keys in the stage parameters. If successful,
// the updated stage is saved. Returns an error if any operation fails.
//
// Parameters:
//
//	ctx - The context for managing request-scoped values and cancellation.
//	stage - The stage entity whose parameters are to be updated.
//
// Returns:
//
//	An error if retrieving sendpost parameters or saving the stage fails.
func (src *SendpostRunnerService) ReplaceStageParametersWithSendpostParameters(ctx context.Context, stage *entity.Stage) error {
	globalParams, err := src.sendpostService.GetSendpostParameters(ctx, stage.SendpostID)
	if err != nil {
		return logging.WrapError(ReplaceStageParametersWithSendpostParametersErr, err)
	}
	stageParams := *stage.StageParameters
	for k, v := range *globalParams {
		if _, ok := (*stage.StageParameters)[k]; ok {
			stageParams[k] = v
		}
	}
	stage.StageParameters = &stageParams
	if err := src.stageService.saveStage(ctx, stage); err != nil {
		return logging.WrapError(ReplaceStageParametersWithSendpostParametersErr, err)
	}
	return nil
}
