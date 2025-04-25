package services

import (
	"context"
	"crm-uplift-ii24-backend/internal/domain/entity"
	"crm-uplift-ii24-backend/internal/domain/value"
	"crm-uplift-ii24-backend/pkg/logging"
	"errors"
	"fmt"
	"time"

	"go.uber.org/zap"
)

const (
	StageCompleted string = "[StageRunnerService] Stage completed"
	StageFailed    string = "[StageRunnerService] Stage failed"
)

type StageRunnerService struct {
	stageService              *StageService
	executor                  entity.StageExecutor
	stageExecutorQueryTimeout int
}

func NewStageRunnerService(executor entity.StageExecutor, stageService *StageService, stageExecutorQueryTimeout int) *StageRunnerService {
	return &StageRunnerService{executor: executor, stageService: stageService, stageExecutorQueryTimeout: stageExecutorQueryTimeout}
}

// Start initiates the execution of a stage by running it through the executor.
// It updates the stage's FlowRunID and state upon successful execution.
// If any error occurs during execution or state update, it returns the error.
//
// Parameters:
//
//	ctx - The context for managing request-scoped values, cancellation, and timeouts.
//	stage - A pointer to the Stage entity containing deployment ID and parameters.
//
// Returns:
//
//	An error if the execution or state update fails, otherwise nil.
func (bsr *StageRunnerService) Start(ctx context.Context, stage *entity.Stage) error {
	flowRunID, state, err := bsr.executor.Run(ctx, stage.DeploymnentID, (*map[string]interface{})(stage.StageParameters))
	if err != nil {
		bsr.HandleFailedStage(ctx, stage, err)
		return fmt.Errorf("[StageRunnerService] error starting stage: %s", err)
	}

	stage.FlowRunID = flowRunID
	stage.State = *state
	if err := bsr.stageService.UpdateStageState(ctx, stage, *state); err != nil {
		bsr.HandleFailedStage(ctx, stage, err)
		return fmt.Errorf("[StageRunnerService] error starting stage: %s", err)
	}

	return nil
}

// CheckState monitors the state of a given stage until it completes, fails, or the context is done.
// It periodically checks the status of the stage using a ticker and handles different states accordingly.
// If the stage fails, it invokes HandleFailedStage. If the stage completes, it updates the stage state.
// Logs errors and completion status using the internal logging package.
//
// Parameters:
//
//	ctx - The context to control cancellation and timeout.
//	stage - The stage entity to monitor.
//
// Returns:
//
//	An error if the context is done or if there is an issue checking the stage status.
func (bsr *StageRunnerService) CheckState(ctx context.Context, stage *entity.Stage) error {

	ticker := time.NewTicker(time.Duration(bsr.stageExecutorQueryTimeout) * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			logging.Error("[StageRunnerService] Timeout exceeded while waiting for stage completion", zap.Uint("stage_id", stage.ID))
			bsr.HandleFailedStage(ctx, stage, errors.New("timeout exceeded while waiting for stage completion"))
			return ctx.Err()
		case <-ticker.C:
			state, err := bsr.executor.Status(ctx, *stage.FlowRunID)
			if err != nil {
				bsr.HandleFailedStage(ctx, stage, err)
				return fmt.Errorf("[StageRunnerService] error geting stage status: %s", err)
			}

			if bsr.IsStageFailed(state) {
				return bsr.HandleFailedStage(ctx, stage, errors.New("stage completed with failed state"))
			}

			if *state == value.Completed {
				logging.Info(StageCompleted, zap.Uint("stage_id", stage.ID))
				return bsr.stageService.UpdateStageState(ctx, stage, *state)
			}
		}
	}
}

// CheckStageCompledSuccesfullyInPeriod verifies if a stage has completed successfully
// within a specified time period. It checks the flow run completion using the stage's
// DeploymentID and updates the stage state accordingly. If the check fails, it handles
// the stage as failed.
//
// Parameters:
//
//	ctx - The context for managing request-scoped values, cancelation, and deadlines.
//	stage - The stage entity to be checked.
//	start - The start time of the period to check.
//	end - The end time of the period to check.
//
// Returns:
//
//	An error if the stage did not complete successfully or if there was an issue
//	updating the stage state.
func (s *StageRunnerService) CheckStageCompledSuccesfullyInPeriod(ctx context.Context, stage *entity.Stage, start time.Time, end time.Time) error {
	if err := s.executor.CheckFlowRunCompletionByDeploymentID(ctx, start, end, stage.DeploymnentID); err != nil {
		return s.HandleFailedStage(ctx, stage, err)
	}
	logging.Info(StageCompleted, zap.Uint("stage_id", stage.ID))
	return s.stageService.UpdateStageState(ctx, stage, value.Completed)
}

// IsStageFailed checks if the given stage state indicates a failure.
// It returns true if the state is Cancelled, Cancelling, Failed, or Crashed.
func (s *StageRunnerService) IsStageFailed(state *value.StateType) bool {
	return *state == value.Cancelled || *state == value.Cancelling || *state == value.Failed || *state == value.Crashed
}

// HandleFailedStage logs a warning for a failed stage and updates its state.
// It takes a context, a pointer to a Stage entity, and a StateType value.
// If updating the stage state fails, it returns the error.
// Otherwise, it returns a generic error indicating stage execution failure.
func (s *StageRunnerService) HandleFailedStage(ctx context.Context, stage *entity.Stage, err error) error {
	logging.Warn(StageFailed, zap.Uint("stage_id", stage.ID), zap.Error(err))
	if err := s.stageService.UpdateStageState(ctx, stage, value.Failed); err != nil {
		return err
	}
	return errors.New(StageFailed)
}
