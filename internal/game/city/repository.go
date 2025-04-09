package city

import (
	"context"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type Repository interface {
	GetCitiesByUser(ctx context.Context, userID uuid.UUID) ([]*City, error)
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
