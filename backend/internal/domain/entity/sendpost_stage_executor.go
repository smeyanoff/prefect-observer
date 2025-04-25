package entity

import (
	"context"
	"crm-uplift-ii24-backend/internal/domain/value"
	"time"
)

type StageExecutor interface {
	Run(ctx context.Context, deploymentID string, parameters *map[string]interface{}) (flowRunID *string, flowRunState *value.StateType, err error)
	Status(ctx context.Context, flowRunID string) (stateType *value.StateType, err error)
	CheckFlowRunCompletionByDeploymentID(ctx context.Context, hisoryStart time.Time, historyEnd time.Time, deploymentID string) error
	GetDeploymentParameters(ctx context.Context, deploymentID string) (map[string]interface{}, error)
}
