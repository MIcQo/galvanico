package user

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

var ErrDuplicateEntry = errors.New("duplicate entry")

type Repository interface {
	GetByUsername(ctx context.Context, username string) (*User, error)
	GetByID(ctx context.Context, id uuid.UUID) (*User, error)
	Create(ctx context.Context, user *User) error
	AddFeature(ctx context.Context, feature *Feature) error
	RemoveFeature(ctx context.Context, feature *Feature) error
	UpdateLastLogin(ctx context.Context, user *User, ip string) error
	ChangeUsername(ctx context.Context, user *User) error
	ChangePassword(ctx context.Context, user *User) error
}

type RepositoryImpl struct {
	db *bun.DB
}

func (r *RepositoryImpl) GetByUsername(ctx context.Context, username string) (*User, error) {
	var user User
	var query = r.db.NewSelect().Model(&user).
		Relation("Features").
		Relation("Resources").
		Where("username = ? OR email = ?", username, username)

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

func (r *RepositoryImpl) UpdateLastLogin(ctx context.Context, user *User, ip string) error {
	user.LastLoginIP.Valid = true
	user.LastLoginIP.String = ip
	user.LastLogin.Valid = true
	user.LastLogin.Time = time.Now().UTC()

	if _, err := r.db.NewUpdate().Model(user).Column("last_login", "last_login_ip").WherePK().Exec(ctx); err != nil {
		return err
	}
	return nil
}

func (r *RepositoryImpl) Create(ctx context.Context, user *User) error {
	if _, err := r.db.NewInsert().Model(user).Exec(ctx); err != nil {
		if strings.Contains(err.Error(), "duplicate") {
			return ErrDuplicateEntry
		}
		return err
	}

	return nil
}

func (r *RepositoryImpl) ChangeUsername(ctx context.Context, user *User) error {
	if _, err := r.db.NewUpdate().Model(user).Column("username").Exec(ctx); err != nil {
		return err
	}
	return nil
}

func (r *RepositoryImpl) ChangePassword(ctx context.Context, user *User) error {
	if _, err := r.db.NewUpdate().Model(user).Column("password").Exec(ctx); err != nil {
		return err
	}
	return nil
}

func NewUserRepository(db *bun.DB) Repository {
	return &RepositoryImpl{db: db}
}

type FakerUserRepository struct {
	data map[string]*User
}

func (f *FakerUserRepository) GetByUsername(_ context.Context, username string) (*User, error) {
	if usr, ok := f.data[username]; ok {
		return usr, nil
	}
	return nil, sql.ErrNoRows
}

func (f *FakerUserRepository) GetByID(_ context.Context, id uuid.UUID) (*User, error) {
	for _, usr := range f.data {
		if usr.ID.String() == id.String() {
			return usr, nil
		}
	}

	return nil, sql.ErrNoRows
}

func (f *FakerUserRepository) AddFeature(_ context.Context, _ *Feature) error {
	panic("implement me")
}

func (f *FakerUserRepository) RemoveFeature(_ context.Context, _ *Feature) error {
	panic("implement me")
}

func (f *FakerUserRepository) UpdateLastLogin(_ context.Context, _ *User, _ string) error {
	return nil
}

func (f *FakerUserRepository) ChangeUsername(_ context.Context, _ *User) error {
	return nil
}

func (f *FakerUserRepository) Create(_ context.Context, usr *User) error {
	if _, ok := f.data[usr.Username]; ok {
		return errors.New("username already exists")
	}

	f.data[usr.Username] = usr

	return nil
}

func (f *FakerUserRepository) ChangePassword(_ context.Context, _ *User) error {
	return nil
}

func NewFakeUserRepository(data map[string]*User) Repository {
	return &FakerUserRepository{data: data}
}
