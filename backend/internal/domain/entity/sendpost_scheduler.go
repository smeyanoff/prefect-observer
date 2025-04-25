package entity

import "context"

type SendpostScheduler interface {
	AddSchedule(ctx context.Context, schedule SendpostSchedule)
	RemoveSchedule(sendpostID uint)
	GetQueue() <-chan SendpostSchedule
}
