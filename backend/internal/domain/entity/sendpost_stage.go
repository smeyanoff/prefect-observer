package entity

import (
	"crm-uplift-ii24-backend/internal/domain/value"
	"errors"

	"gorm.io/gorm"
)

type Stage struct {
	gorm.Model
	SendpostID uint      `gorm:"not_null;index"`
	Sendpost   *Sendpost `gorm:"foreignkey:SendpostID;references:ID;constraint:OnDelete:CASCADE;"`

	State value.StateType `gorm:"size:20;default:NEVERRUNNING;not null"`
	Type  value.StageType `gorm:"size:20;default:SEQUENTIAL;not null"`

	DeploymnentID   string       `gorm:"size:255"`
	FlowRunID       *string      `gorm:"size:255"`
	StageParameters *value.JSONB `gorm:"type:jsonb"`

	NextStageID *uint  `gorm:"index"`
	NextStage   *Stage `gorm:"foreignKey:NextStageID"`

	ParentStageID *uint    `gorm:"index"`
	SubStages     []*Stage `gorm:"foreignKey:ParentStageID;constraint:OnDelete:CASCADE;"`

	IsBlocked bool `gorm:"default:false"`
}

func (s *Stage) IsParallel() bool {
	return s.Type == value.ParallelStage
}

func (s *Stage) UpdateState(state value.StateType) {
	s.State = state
}

func (s *Stage) UpdateNextStageID(nextStageID *uint) error {
	if s.ParentStageID != nil {
		return errors.New("parent stage id and next stage id couldn't both be fullfield")
	}
	s.NextStageID = nextStageID
	return nil
}

func (s *Stage) Copy(newSendpostId uint) *Stage {
	return &Stage{
		SendpostID:      newSendpostId,
		Type:            s.Type,
		DeploymnentID:   s.DeploymnentID,
		StageParameters: s.StageParameters,
		IsBlocked:       s.IsBlocked,
	}
}
