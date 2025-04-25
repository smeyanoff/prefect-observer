package repository

import (
	"context"
	"crm-uplift-ii24-backend/internal/domain/entity"
	"crm-uplift-ii24-backend/internal/domain/value"
)

type SendpostRepository interface {
	SaveSendpost(ctx context.Context, sendpost *entity.Sendpost) error
	DeleteSendpost(ctx context.Context, sendpostID uint) error
	GetSendpostByID(ctx context.Context, sendpostID uint) (*entity.Sendpost, error)
	GetSendposts(ctx context.Context) ([]*entity.Sendpost, error)
	GetFirstStage(ctx context.Context, sendpostID uint) (*entity.Stage, error)
	GetSendpostParameters(ctx context.Context, sendpostID uint) (*value.JSONB, error)
	UpdateSendpostParameters(ctx context.Context, sendpostID uint, parameters *map[string]interface{}) error
}
