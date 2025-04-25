package repository

import (
	"context"
	"crm-uplift-ii24-backend/internal/domain/entity"
)

type SendpostScheduleRepository interface {
	SaveSchedule(ctx context.Context, schedule *entity.SendpostSchedule) error
	DeleteSchedule(ctx context.Context, scheduleID uint) error
	GetScheduleByID(ctx context.Context, scheduleID uint) (*entity.SendpostSchedule, error)
	GetScheduleBySendpostID(ctx context.Context, sendpostID uint) (*entity.SendpostSchedule, error)
}
