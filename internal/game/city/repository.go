package city

import (
	"context"
	"errors"
	"galvanico/internal/game/building"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

// Repository provides methods for working with City data
type Repository interface {
	// GetCitiesByUser retrieves all cities associated with a given user ID.
	GetCitiesByUser(ctx context.Context, userID uuid.UUID) ([]*City, error)
	// CreateCity persists a new city with all its related entities.
	CreateCity(ctx context.Context, city *City) error
	// GetCityByID returns city by their id
	GetCityByID(ctx context.Context, id uuid.UUID) (*City, error)
}

type RepositoryImpl struct {
	db *bun.DB
}

func NewRepository(db *bun.DB) *RepositoryImpl {
	return &RepositoryImpl{db: db}
}

func (r *RepositoryImpl) GetCitiesByUser(ctx context.Context, userID uuid.UUID) ([]*City, error) {
	var cities = make([]*City, 0)
	var err = r.db.NewSelect().Model(&cities).
		Relation("UserCity").
		Relation("Buildings").
		Relation("Resources").
		Where("user_id = ?", userID).
		Scan(ctx)
	if err != nil {
		return nil, err
	}
	return cities, nil
}

func (r *RepositoryImpl) CreateCity(ctx context.Context, city *City) error {
	var tx, err = r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	var _, cityErr = tx.NewInsert().Model(city).Exec(ctx)
	if cityErr != nil {
		return cityErr
	}

	var _, buildingsErr = tx.NewInsert().Model(&city.Buildings).Exec(ctx)
	if buildingsErr != nil {
		return buildingsErr
	}

	var _, resourcesErr = tx.NewInsert().Model(&city.Resources).Exec(ctx)
	if resourcesErr != nil {
		return resourcesErr
	}
	var _, userErr = tx.NewInsert().Model(&city.UserCity).Exec(ctx)
	if userErr != nil {
		return userErr
	}

	return tx.Commit()
}

func (r *RepositoryImpl) GetCityByID(ctx context.Context, id uuid.UUID) (*City, error) {
	if id == uuid.Nil {
		return nil, errors.New("city not found")
	}

	var city = new(City)
	var err = r.db.NewSelect().Model(&city).
		Relation("UserCity").
		Relation("Buildings").
		Relation("Resources").
		Where("id = ?", id).
		Scan(ctx)

	return city, err
}

type FakeRepository struct {
	data map[string][]*City
}

func (f *FakeRepository) GetCitiesByUser(_ context.Context, userID uuid.UUID) ([]*City, error) {
	return f.data[userID.String()], nil
}

func (f *FakeRepository) CreateCity(_ context.Context, city *City) error {
	f.data[city.UserCity.UserID.String()] = append(f.data[city.UserCity.UserID.String()], city)
	return nil
}

func (f *FakeRepository) GetCityByID(_ context.Context, id uuid.UUID) (*City, error) {
	for _, cities := range f.data {
		for _, city := range cities {
			if city.ID.String() == id.String() {
				return city, nil
			}
		}
	}

	return nil, errors.New("city not found")
}

func NewFakeRepository() Repository {
	var cityID = uuid.MustParse("c46d4d06-05c8-46a0-945c-ad02dfe73e8b")
	var secondCityID = uuid.MustParse("7d5c12c7-a09f-4b0f-b912-a65a0cf2c997")
	var userID = uuid.MustParse("9efa2461-e40a-423a-a734-ce29f302437b")

	return &FakeRepository{data: map[string][]*City{
		"9efa2461-e40a-423a-a734-ce29f302437b": {
			{
				ID:        cityID,
				Name:      "City",
				PositionX: 1,
				PositionY: 1,
				UserCity: UserCity{
					UserID: userID,
					CityID: cityID,
				},
				Buildings: []Building{
					{
						CityID:   cityID,
						Building: building.CityHall,
						Level:    1,
						Position: 0,
					},
				},
				Resources: Resources{
					CityID: cityID,
				},
			},
			{
				ID:        secondCityID,
				Name:      "City",
				PositionX: 1,
				PositionY: 1,
				UserCity: UserCity{
					UserID: userID,
					CityID: secondCityID,
				},
				Buildings: []Building{
					{
						CityID:   secondCityID,
						Building: building.CityHall,
						Level:    1,
						Position: 0,
					},
					{
						CityID:   secondCityID,
						Building: building.IndustrialPort,
						Level:    1,
						Position: PortSlot2,
					},
					{
						CityID:   secondCityID,
						Building: building.Cathedral,
						Level:    1,
						Position: 1,
					},
				},
				Resources: Resources{
					CityID: cityID,
				},
			},
		},
	}}
}
