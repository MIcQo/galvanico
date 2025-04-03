package user

import (
	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type Feature struct {
	bun.BaseModel `bun:"table:user_features,alias:uf"`

	UserID  uuid.UUID `bun:"type:uuid,pk,default:gen_random_uuid()"`
	Feature string
}
