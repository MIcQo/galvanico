package city

import (
	"context"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type Repository interface {
	GetCitiesByUser(ctx context.Context, userID uuid.UUID) ([]*City, error)
	CreateCity(ctx context.Context, city *City) error
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
	var _, cityErr = r.db.NewInsert().Model(city).Exec(ctx)
	if cityErr != nil {
		return cityErr
	}

	var _, buildingsErr = r.db.NewInsert().Model(&city.Buildings).Exec(ctx)
	if buildingsErr != nil {
		return buildingsErr
	}

	var _, resourcesErr = r.db.NewInsert().Model(&city.Resources).Exec(ctx)
	if resourcesErr != nil {
		return resourcesErr
	}
	var _, userErr = r.db.NewInsert().Model(&city.UserCity).Exec(ctx)
	if userErr != nil {
		return userErr
	}

	return nil
}
