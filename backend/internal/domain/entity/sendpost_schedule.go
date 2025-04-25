package entity

import (
	"time"

	"gorm.io/gorm"
)

type SendpostSchedule struct {
	gorm.Model

	SendpostID uint      `gorm:"not_null;index"`
	Sendpost   *Sendpost `gorm:"foreignKey:SendpostID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;references:ID"`

	PlannedAt   time.Time
	StartedAt   *time.Time
	CompletedAt *time.Time
}
