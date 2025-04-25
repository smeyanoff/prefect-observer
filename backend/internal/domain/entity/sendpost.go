package entity

import (
	"crm-uplift-ii24-backend/internal/domain/value"
	"errors"

	"gorm.io/gorm"
)

type Sendpost struct {
	gorm.Model
	SendpostName string `gorm:"not null"`

	Description *string

	FirstStageID *uint
	FirstStage   *Stage `gorm:"foreignkey:FirstStageID;references:ID;constraint:OnDelete:SET NULL;"`

	State value.StateType `gorm:"size:20;default:NEVERRUNNING;not null"`

	GlobalParameters *value.JSONB `gorm:"type:jsonb"`
}

func (s *Sendpost) Update(sendpostName string, description *string, parameters *value.JSONB) error {
	if sendpostName == "" {
		return errors.New("sendpost name couldnt be empty")
	}
	s.SendpostName = sendpostName

	if description != nil {
		s.Description = description
	}

	if parameters != nil {
		s.GlobalParameters = parameters
	}

	return nil
}

func (s *Sendpost) Copy(NewName string, NewDescription *string) *Sendpost {
	return &Sendpost{
		SendpostName:     NewName,
		Description:      NewDescription,
		GlobalParameters: s.GlobalParameters,
	}
}
