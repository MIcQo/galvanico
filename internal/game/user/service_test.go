package user

import (
	"database/sql"
	"galvanico/internal/utils"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

type fakePublisher struct{}

func (f fakePublisher) Publish(_ string, _ []byte) error {
	return nil
}

func TestServiceIml_Register(t *testing.T) {
	var svc = NewService(&FakerUserRepository{data: map[string]*User{
		"test": {
			ID:       uuid.MustParse("ebfc76b7-7ace-4034-a8b6-cc369afa8fb8"),
			Username: "test",
		},
	}}, &fakePublisher{})

	t.Run("valid", func(t *testing.T) {
		var err = svc.Register(t.Context(), &User{
			ID:       uuid.New(),
			Password: sql.NullString{String: "password", Valid: true},
		})
		require.NoError(t, err)
	})

	t.Run("empty password", func(t *testing.T) {
		var err = svc.Register(t.Context(), &User{
			ID:       uuid.New(),
			Password: sql.NullString{String: "", Valid: true},
		})
		require.Error(t, err)
		require.ErrorIs(t, err, ErrEmptyPassword)
	})

	t.Run("too long password", func(t *testing.T) {
		var err = svc.Register(t.Context(), &User{
			ID:       uuid.New(),
			Password: sql.NullString{String: utils.RandomString(73), Valid: true},
		})
		require.Error(t, err)
		require.ErrorIs(t, err, bcrypt.ErrPasswordTooLong)
	})
}

func TestServiceIml_GetUser(t *testing.T) {
	var svc = NewService(&FakerUserRepository{data: map[string]*User{
		"test": {
			ID:       uuid.MustParse("ebfc76b7-7ace-4034-a8b6-cc369afa8fb8"),
			Username: "test",
		},
	}}, &fakePublisher{})

	t.Run("valid", func(t *testing.T) {
		var token = jwt.NewWithClaims(
			jwt.SigningMethodHS256,
			jwt.MapClaims{
				"sub": "ebfc76b7-7ace-4034-a8b6-cc369afa8fb8",
				"iss": time.Now().String(),
			},
		)
		var usr, err = svc.GetUser(t.Context(), token)
		require.NoError(t, err)
		require.NotNil(t, usr)
	})
}
