package user

import (
	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type Resource struct {
	bun.BaseModel `bun:"table:user_resources,alias:ur"`

	UserID         uuid.UUID
	Gold           float64
	Ships          int
	AvailableShips int
	Electricity    float64
	MaxElectricity float64
	Waste          float64
	MaxWaste       float64
}
