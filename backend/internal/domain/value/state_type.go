package value

import (
	"encoding/json"
	"fmt"
)

// StateType представляет возможные состояния задачи
type StateType string

const (
	Scheduled    StateType = "SCHEDULED"
	Pending      StateType = "PENDING"
	Running      StateType = "RUNNING"
	Completed    StateType = "COMPLETED"
	Failed       StateType = "FAILED"
	Cancelled    StateType = "CANCELLED"
	Crashed      StateType = "CRASHED"
	Paused       StateType = "PAUSED"
	Cancelling   StateType = "CANCELLING"
	NeverRunning StateType = "NEVERRUNNING"
	Updated      StateType = "UPDATED"
)

func (st StateType) IsValid() bool {
	switch st {
	case Scheduled, Pending, Running, Completed, Failed, Cancelled, Crashed, Paused, Cancelling, Updated, NeverRunning:
		return true
	default:
		return false
	}
}

func (st *StateType) UnmarshalJSON(data []byte) error {
	var state string
	if err := json.Unmarshal(data, &state); err != nil {
		return err
	}

	*st = StateType(state)
	if !st.IsValid() {
		return fmt.Errorf("unknown state type: %s", state)
	}

	return nil
}

func (st *StateType) String() string {
	return string(*st)
}
