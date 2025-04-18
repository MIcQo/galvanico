package city

import (
	"galvanico/internal/game/building"
	"time"

	"github.com/uptrace/bun"

	"github.com/google/uuid"
)

const (
	DefaultCityName          = "City"
	DefaultTownHallLevel     = 1
	DefaultTownHallPosition  = 0
	DefaultWood              = 1000
	DefaultWarehouseCapacity = 2500
	DefaultPopulation        = 40
	DefaultMaxPopulation     = 100

	MinPosition = 0
	MaxPosition = 100
)

type City struct {
	ID        uuid.UUID `bun:"type:uuid,pk,default:gen_random_uuid()"`
	Name      string
	PositionX int64 `bun:"column:position_x"`
	PositionY int64 `bun:"column:position_y"`

	UserCity  UserCity   `bun:"rel:has-one,join:id=city_id"`
	Buildings []Building `bun:"rel:has-many,join:id=city_id"`
	Resources Resources  `bun:"rel:has-one,join:id=city_id"`
}

type UserCity struct {
	UserID uuid.UUID
	CityID uuid.UUID
}

type Resources struct {
	bun.BaseModel     `bun:"table:city_resources,alias:cr"`
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
	bun.BaseModel `bun:"table:city_buildings,alias:cb"`
	CityID        uuid.UUID
	Building      building.Building
	Level         uint
	Position      uint
}
