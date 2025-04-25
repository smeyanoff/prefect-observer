package prefectV2

import "time"

type Interval struct {
	IntervalStart time.Time `json:"interval_start"`
	IntervalEnd   time.Time `json:"interval_end"`
	States        []State   `json:"states"`
}

// State представляет информацию о состоянии запуска потока
type State struct {
	StateType            string  `json:"state_type"`
	StateName            string  `json:"state_name"`
	CountRuns            int     `json:"count_runs"`
	SumEstimatedRunTime  float64 `json:"sum_estimated_run_time"`
	SumEstimatedLateness float64 `json:"sum_estimated_lateness"`
}
