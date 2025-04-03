package user

import (
	"context"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type Repository interface {
	GetByUsername(ctx context.Context, username string) (*User, error)
	GetByID(ctx context.Context, id uuid.UUID) (*User, error)
	AddFeature(ctx context.Context, feature *Feature) error
	RemoveFeature(ctx context.Context, feature *Feature) error
}

type RepositoryImpl struct {
	db *bun.DB
}

func (r *RepositoryImpl) GetByUsername(ctx context.Context, username string) (*User, error) {
	var user User
	var query = r.db.NewSelect().Model(&user).
		Relation("Features").
		Relation("Resources").
		Where("username = ?", username)

	if err := query.Scan(ctx); err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *RepositoryImpl) GetByID(ctx context.Context, id uuid.UUID) (*User, error) {
	var user User
	var query = r.db.NewSelect().Model(&user).
		Relation("Features").
		Relation("Resources").
		Where("id = ?", id)

	if err := query.Scan(ctx); err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *RepositoryImpl) AddFeature(ctx context.Context, feature *Feature) error {
	if _, err := r.db.NewInsert().Model(feature).Exec(ctx); err != nil {
		return err
	}
	return nil
}

func (r *RepositoryImpl) RemoveFeature(ctx context.Context, feature *Feature) error {
	if _, err := r.db.NewDelete().Model(feature).Exec(ctx); err != nil {
		return err
	}
	return nil
}

func NewUserRepository(db *bun.DB) Repository {
	return &RepositoryImpl{db: db}
}
