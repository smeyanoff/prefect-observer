package requests

import "crm-uplift-ii24-backend/internal/domain/value"

type Senpost struct {
	SendpostName     string       `binding:"required" json:"sendpost_name"`
	Description      *string      `json:"description"`
	GlobalParameters *value.JSONB `json:"global_parameters"`
}
