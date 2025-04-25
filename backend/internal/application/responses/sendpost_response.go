package responses

import "crm-uplift-ii24-backend/internal/domain/value"

type Sendposts []*Sendpost

type Sendpost struct {
	ID               uint         `json:"id" validate:"required"`
	Name             string       `json:"name" validate:"required"`
	Description      *string      `json:"description"`
	State            string       `json:"state" validate:"required"`
	GlobalParameters *value.JSONB `json:"global_parameters"`
}
