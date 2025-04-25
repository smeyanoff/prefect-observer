package requests

import (
	"crm-uplift-ii24-backend/internal/domain/value"
)

type Stage struct {
	StageType       value.StageType `json:"type" binding:"required"`
	DeploymentID    string          `json:"deployment_id" binding:"required"`
	StageParameters *value.JSONB    `json:"stage_parameters"`
	PreviousStageID *uint           `json:"previous_stage_id"`
}
