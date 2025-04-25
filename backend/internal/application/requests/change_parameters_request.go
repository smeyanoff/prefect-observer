package requests

import "crm-uplift-ii24-backend/internal/domain/value"

type Parameters struct {
	Parameters value.JSONB `json:"parameters" binding:"required"`
}
