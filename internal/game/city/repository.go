package city

import (
	"context"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

// Repository provides methods for working with City data
type Repository interface {
	// GetCitiesByUser retrieves all cities associated with a given user ID.
	GetCitiesByUser(ctx context.Context, userID uuid.UUID) ([]*City, error)
	// CreateCity persists a new city with all its related entities.
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
