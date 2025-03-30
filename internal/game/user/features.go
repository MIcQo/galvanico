package user

import (
	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type Feature struct {
	bun.BaseModel `bun:"table:user_features,alias:uf"`

	UserID  uuid.UUID
	Feature string
}
