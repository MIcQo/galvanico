package city

import (
	"context"
	"galvanico/internal/game/building"
	"galvanico/internal/utils"
	"time"

	"github.com/google/uuid"
)

type Service interface {
	InitCity(context context.Context, id uuid.UUID) (*City, error)
}

type ServiceImpl struct {
	repository Repository
}

func (s *ServiceImpl) randomCoordinate() (int64, error) {
	return utils.RandomNumber(MinPosition, MaxPosition)
}

func (s *ServiceImpl) InitCity(ctx context.Context, id uuid.UUID) (*City, error) {
	var randomX, errX = s.randomCoordinate()
	if errX != nil {
		return nil, errX
	}

	var randomY, errY = s.randomCoordinate()
	if errY != nil {
		return nil, errY
	}

	var cityID = uuid.Must(uuid.NewRandom())
	var city = &City{
		ID:        cityID,
		Name:      DefaultCityName,
		PositionX: randomX,
		PositionY: randomY,
		UserCity: UserCity{
			CityID: cityID,
			UserID: id,
		},
		Buildings: []Building{
			{
				CityID:   cityID,
				Building: building.CityHall,
				Level:    DefaultTownHallLevel,
				Position: DefaultTownHallPosition,
			},
		},
		Resources: Resources{
			CityID:            cityID,
			Wood:              DefaultWood,
			WarehouseCapacity: DefaultWarehouseCapacity,
			Population:        DefaultPopulation,
			MaxPopulation:     DefaultMaxPopulation,
			UpdatedAt:         time.Now(),
		},
	}

	var err = s.repository.CreateCity(ctx, city)
	if err != nil {
		return nil, err
	}

	return city, nil
}

func NewService(repository Repository) *ServiceImpl {
	return &ServiceImpl{repository: repository}
}

type FakeService struct {
	repo Repository
}

func (f *FakeService) InitCity(_ context.Context, id uuid.UUID) (*City, error) {
	return &City{ID: id}, nil
}

func NewFakeService(repo Repository) *FakeService {
	return &FakeService{repo: repo}
}
