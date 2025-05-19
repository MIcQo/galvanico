package user

import (
	"database/sql"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func generateTestToken(userID uuid.UUID, secret string) string {
	claims := jwt.MapClaims{
		"sub": userID.String(),
		"exp": time.Now().Add(time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	s, _ := token.SignedString([]byte(secret))
	return s
}

func setupAppWithMiddleware(secret string, middleware *Middleware) *fiber.App {
	app := fiber.New()

	app.Use(jwtware.New(jwtware.Config{
		SigningKey:  jwtware.SigningKey{Key: []byte(secret)},
		ContextKey:  "user",
		TokenLookup: "header:Authorization",
		AuthScheme:  "Bearer",
	}))

	app.Use(middleware.CheckNotBanned())

	app.Get("/", func(c *fiber.Ctx) error {
		user := c.Locals("user_data").(*User)
		return c.SendString("Welcome " + user.Username + "!")
	})

	return app
}

func TestMiddleware_WithStruct(t *testing.T) {
	secret := "testsecret"

	invalidID := uuid.MustParse("99999999-9999-9999-9999-999999999999")
	validID := uuid.MustParse("11111111-1111-1111-1111-111111111111")
	bannedID := uuid.MustParse("22222222-2222-2222-2222-222222222222")
	repo := NewFakeUserRepository(map[string]*User{
		validID.String(): {
			ID:       validID,
			Username: "Alice",
		},
		bannedID.String(): {
			ID:            bannedID,
			Username:      "Bob",
			BanExpiration: sql.NullTime{Valid: true, Time: time.Now().Add(time.Hour * 2)},
		},
	})
	middleware := NewUserMiddleware(repo)

	t.Run("User not banned", func(t *testing.T) {
		app := setupAppWithMiddleware(secret, middleware)
		token := generateTestToken(validID, secret)

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set("Authorization", "Bearer "+token)

		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, 200, resp.StatusCode)
	})

	t.Run("User is banned", func(t *testing.T) {
		app := setupAppWithMiddleware(secret, middleware)
		token := generateTestToken(bannedID, secret)

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set("Authorization", "Bearer "+token)

		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, 403, resp.StatusCode)
	})

	t.Run("User not found", func(t *testing.T) {
		app := setupAppWithMiddleware(secret, middleware)
		token := generateTestToken(invalidID, secret)

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set("Authorization", "Bearer "+token)

		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, 404, resp.StatusCode)
	})
}
