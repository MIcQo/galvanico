package user

import (
	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type Resources struct {
	bun.BaseModel `bun:"table:user_resources,alias:ur"`

	UserID             uuid.UUID `bun:"type:uuid,pk,default:gen_random_uuid()"`
	Gold               int
	TrainCars          int
	AvailableTrainCars int
	Electricity        float64
	MaxElectricity     float64
	Waste              float64
	MaxWaste           float64
}
