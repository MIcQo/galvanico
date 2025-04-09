package city

import (
	"galvanico/internal/game/building"
	"time"

	"github.com/google/uuid"
)

type City struct {
	ID        uuid.UUID `bun:"type:uuid,pk,default:gen_random_uuid()"`
	Name      string
	PositionX int `bun:"column:position_x"`
	PositionY int `bun:"column:position_y"`

	UserCity  UserCity   `bun:"rel:has-one,join:id=city_id"`
	Buildings []Building `bun:"rel:has-many,join:id=city_id"`
	Resources Resources  `bun:"rel:has-one,join:id=city_id"`
}

type UserCity struct {
	UserID uuid.UUID
	CityID uuid.UUID
}

type Resources struct {
	CityID            uuid.UUID
	Wood              float64
	Water             float64
	Iron              float64
	Oil               float64
	Cotton            float64
	WarehouseCapacity float64
	UpdatedAt         time.Time
	Population        float64
	MaxPopulation     float64
}

type Building struct {
	CityID   uuid.UUID
	Building building.Building
	Level    uint
}
