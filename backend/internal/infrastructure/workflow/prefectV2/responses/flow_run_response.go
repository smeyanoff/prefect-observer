package prefectV2

import "crm-uplift-ii24-backend/internal/domain/value"

type FlowRunResponse struct {
	FlowID        string          `json:"id"`
	Name          string          `json:"name"`
	DeploymnentID string          `json:"deployment_id"`
	StateType     value.StateType `json:"state_type"`
}
