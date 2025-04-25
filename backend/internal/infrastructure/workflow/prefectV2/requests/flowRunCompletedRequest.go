package prefectV2

import "time"

type FlowRunCompletedRequest struct {
	HistoryStart           time.Time `json:"history_start"`
	HistoryEnd             time.Time `json:"history_end"`
	HistoryIntervalSeconds int       `json:"history_interval_seconds"`
	FlowRuns               FlowRuns  `json:"flow_runs"`
	Sort                   string    `json:"sort"`
	Limit                  int       `json:"limit"`
}

type FlowRuns struct {
	DeploymentID DeploymentFilter `json:"deployment_id"`
	State        StateFilter      `json:"state"`
	StartTime    TimeFilter       `json:"start_time"`
}

type DeploymentFilter struct {
	Any []string `json:"any_"`
}

type StateFilter struct {
	Type StateType `json:"type"`
}

type StateType struct {
	Any []string `json:"any_"`
}

type TimeFilter struct {
	After  time.Time `json:"after_"`
	Before time.Time `json:"before_"`
}
