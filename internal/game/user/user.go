package user

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type User struct {
	bun.BaseModel `bun:"table:users,alias:u"`

	ID            uuid.UUID
	ExternalID    sql.NullInt64
	Status        string
	Username      string
	LastLogin     sql.NullTime
	LastLoginIP   sql.Null[[]byte]
	Language      string
	BanExpiration sql.NullTime
	BanReason     sql.NullString
	CreatedAt     time.Time
	UpdatedAt     sql.NullTime
	DeletedAt     sql.NullTime
}
