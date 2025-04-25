package responses

import (
	"crm-uplift-ii24-backend/internal/domain/value"
)

type StageDetailed struct {
	ID              uint            `json:"id" validate:"required"`
	Type            value.StageType `json:"type" validate:"required"`
	State           value.StateType `json:"state" validate:"required"`
	ParentStageID   *uint           `json:"parent_stage_id"`
	DeploymentID    string          `json:"deployment_id" validate:"required"`
	StageParameters *value.JSONB    `json:"stage_parameters"`
}

type SendpostStages []*Stage

type Stage struct {
	ID        uint            `json:"id" validate:"required"`
	Type      value.StageType `json:"type" validate:"required"`
	State     value.StateType `json:"state" validate:"required"`
	IsBlocked bool            `json:"is_blocked" validate:"required"`
}
