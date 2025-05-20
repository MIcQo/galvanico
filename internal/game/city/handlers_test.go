package city

import (
	"database/sql"
	"galvanico/internal/auth"
	"galvanico/internal/config"
	"galvanico/internal/game/building"
	"galvanico/internal/game/user"
	"galvanico/internal/utils"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/goccy/go-json"
	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

var signingKey = utils.RandomString(32)
var testUserID = uuid.MustParse("9efa2461-e40a-423a-a734-ce29f302437b")
var bannedUserID = uuid.MustParse("0c06a637-73cb-4234-a1db-77034215bc6f")

func setup() (*fiber.App, *Handler) {
	var pass, err = bcrypt.GenerateFromPassword([]byte("test"), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}

	var banTime, parseErr = time.Parse("2006-01-02 15:04:05", "2030-01-01 00:00:00")
	if parseErr != nil {
		panic(parseErr)
	}

	var userRepo = user.NewFakeUserRepository(map[string]*user.User{
		"test": {
			Username: "test",
			Password: sql.NullString{Valid: true, String: string(pass)},
			ID:       testUserID,
		},
		"banned": {
			Username:      "banned",
			Password:      sql.NullString{Valid: true, String: string(pass)},
			ID:            bannedUserID,
			BanExpiration: sql.NullTime{Time: banTime, Valid: true},
			BanReason:     sql.NullString{Valid: true, String: "banned"},
		},
	})
	var userSvc = user.NewFakeService(userRepo)
	var repo = NewFakeRepository()
	var svc = NewFakeService(repo)

	var app = fiber.New()
	var handler = NewHandler(repo, svc, userSvc)

	return app, handler
}

func TestHandler_HandleGetUserCities(t *testing.T) {
	var app, handler = setup()
	var cfg = config.NewDefaultConfig()
	cfg.Auth.Settings["key"] = signingKey
	app.Use(jwtware.New(jwtware.Config{
		SigningKey: jwtware.SigningKey{Key: []byte(signingKey), JWTAlg: jwtware.HS256},
	}))

	app.Get("/", handler.HandleGetUserCities)

	t.Run("should return 200 OK", func(t *testing.T) {
		req, err := http.NewRequest(
			http.MethodGet,
			"/",
			nil,
		)
		require.NoError(t, err)

		token, err := auth.GenerateJWT(cfg, testUserID)
		require.NoError(t, err)

		req.Header.Add("Content-Type", "application/json")
		req.Header.Add("Authorization", "Bearer "+token)
		res, err := app.Test(req, -1)
		require.NoError(t, err)

		require.Equal(t, fiber.StatusOK, res.StatusCode)

		bodyBytes, err := io.ReadAll(res.Body)
		require.NoError(t, err)
		var cities []*City
		err = json.Unmarshal(bodyBytes, &cities)
		require.NoError(t, err)

		// Verify city data structure
		require.NotEmpty(t, cities, "Expected at least one city in the response")

		// If you expect specific attributes in the response, verify them
		for _, city := range cities {
			require.NotEmpty(t, city.ID, "City ID should not be empty")
			// Add more assertions based on your city model
		}
	})

	t.Run("should return 401 Unauthorized without token", func(t *testing.T) {
		req, err := http.NewRequest(
			http.MethodGet,
			"/",
			nil,
		)
		require.NoError(t, err)

		req.Header.Add("Content-Type", "application/json")
		req.Header.Add("Authorization", "Bearer "+"invalid")

		// No Authorization header added
		res, err := app.Test(req, -1)
		require.NoError(t, err)

		require.Equal(t, fiber.StatusUnauthorized, res.StatusCode)
	})

	t.Run("should forbid banned user", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/", nil)
		token, err := auth.GenerateJWT(cfg, bannedUserID)
		require.NoError(t, err)
		req.Header.Set("Authorization", "Bearer "+token)
		res, err := app.Test(req, -1)
		require.NoError(t, err)
		require.Equal(t, fiber.StatusForbidden, res.StatusCode)
	})
}

func TestHandler_HandleAvailableSlotBuildings(t *testing.T) {
	var app, handler = setup()
	var cfg = config.NewDefaultConfig()
	cfg.Auth.Settings["key"] = signingKey
	app.Use(jwtware.New(jwtware.Config{
		SigningKey: jwtware.SigningKey{Key: []byte(signingKey), JWTAlg: jwtware.HS256},
	}))

	app.Get("/:city/:slot/available", handler.HandleAvailableSlotBuildings)

	type availableBuildingsResponse struct {
		Buildings []building.Building `json:"buildings"`
		Slot      int                 `json:"slot"`
	}

	t.Run("should return 200 OK standard building", func(t *testing.T) {
		req, err := http.NewRequest(
			http.MethodGet,
			"/c46d4d06-05c8-46a0-945c-ad02dfe73e8b/1/available",
			nil,
		)
		require.NoError(t, err)

		token, err := auth.GenerateJWT(cfg, testUserID)
		require.NoError(t, err)

		req.Header.Add("Content-Type", "application/json")
		req.Header.Add("Authorization", "Bearer "+token)
		res, err := app.Test(req, -1)
		require.NoError(t, err)

		require.Equal(t, fiber.StatusOK, res.StatusCode)

		bodyBytes, err := io.ReadAll(res.Body)
		require.NoError(t, err)

		var response availableBuildingsResponse
		err = json.Unmarshal(bodyBytes, &response)
		require.NoError(t, err)
		require.Len(t, response.Buildings, 24, "expected 24 buildings")
	})

	t.Run("should return 200 OK port building", func(t *testing.T) {
		req, err := http.NewRequest(
			http.MethodGet,
			"/c46d4d06-05c8-46a0-945c-ad02dfe73e8b/5/available",
			nil,
		)
		require.NoError(t, err)

		token, err := auth.GenerateJWT(cfg, testUserID)
		require.NoError(t, err)

		req.Header.Add("Content-Type", "application/json")
		req.Header.Add("Authorization", "Bearer "+token)
		res, err := app.Test(req, -1)
		require.NoError(t, err)

		require.Equal(t, fiber.StatusOK, res.StatusCode)

		bodyBytes, err := io.ReadAll(res.Body)
		require.NoError(t, err)

		var response availableBuildingsResponse
		err = json.Unmarshal(bodyBytes, &response)
		require.NoError(t, err)
		require.Len(t, response.Buildings, 2, "expected 2 buildings")
	})

	t.Run("should return 200 OK slot with standard building that is already built", func(t *testing.T) {
		req, err := http.NewRequest(
			http.MethodGet,
			"/c46d4d06-05c8-46a0-945c-ad02dfe73e8b/8/available",
			nil,
		)
		require.NoError(t, err)

		token, err := auth.GenerateJWT(cfg, testUserID)
		require.NoError(t, err)

		req.Header.Add("Content-Type", "application/json")
		req.Header.Add("Authorization", "Bearer "+token)
		res, err := app.Test(req, -1)
		require.NoError(t, err)

		require.Equal(t, fiber.StatusOK, res.StatusCode)

		bodyBytes, err := io.ReadAll(res.Body)
		require.NoError(t, err)

		var response availableBuildingsResponse
		err = json.Unmarshal(bodyBytes, &response)
		require.NoError(t, err)
		require.Len(t, response.Buildings, 1, "expected 1 buildings")
	})

	t.Run("should return 200 OK already bult cathedral", func(t *testing.T) {
		req, err := http.NewRequest(
			http.MethodGet,
			"/7d5c12c7-a09f-4b0f-b912-a65a0cf2c997/1/available",
			nil,
		)
		require.NoError(t, err)

		token, err := auth.GenerateJWT(cfg, testUserID)
		require.NoError(t, err)

		req.Header.Add("Content-Type", "application/json")
		req.Header.Add("Authorization", "Bearer "+token)
		res, err := app.Test(req, -1)
		require.NoError(t, err)

		require.Equal(t, fiber.StatusOK, res.StatusCode)

		bodyBytes, err := io.ReadAll(res.Body)
		require.NoError(t, err)

		var response availableBuildingsResponse
		err = json.Unmarshal(bodyBytes, &response)
		require.NoError(t, err)
		require.Len(t, response.Buildings, 23, "expected 23 buildings")
	})

	t.Run("should return 200 OK port building filtered", func(t *testing.T) {
		req, err := http.NewRequest(
			http.MethodGet,
			"/7d5c12c7-a09f-4b0f-b912-a65a0cf2c997/5/available",
			nil,
		)
		require.NoError(t, err)

		token, err := auth.GenerateJWT(cfg, testUserID)
		require.NoError(t, err)

		req.Header.Add("Content-Type", "application/json")
		req.Header.Add("Authorization", "Bearer "+token)
		res, err := app.Test(req, -1)
		require.NoError(t, err)

		require.Equal(t, fiber.StatusOK, res.StatusCode)

		bodyBytes, err := io.ReadAll(res.Body)
		require.NoError(t, err)

		var response availableBuildingsResponse
		err = json.Unmarshal(bodyBytes, &response)
		require.NoError(t, err)
		require.Len(t, response.Buildings, 1, "expected 1 buildings")
	})

	t.Run("should return 400 Bad Request slot out of range positive ", func(t *testing.T) {
		req, err := http.NewRequest(
			http.MethodGet,
			"/7d5c12c7-a09f-4b0f-b912-a65a0cf2c997/50/available",
			nil,
		)
		require.NoError(t, err)

		token, err := auth.GenerateJWT(cfg, testUserID)
		require.NoError(t, err)

		req.Header.Add("Content-Type", "application/json")
		req.Header.Add("Authorization", "Bearer "+token)
		res, err := app.Test(req, -1)
		require.NoError(t, err)

		require.Equal(t, fiber.StatusBadRequest, res.StatusCode)
	})

	t.Run("should return 400 Bad Request slot out of range negative ", func(t *testing.T) {
		req, err := http.NewRequest(
			http.MethodGet,
			"/7d5c12c7-a09f-4b0f-b912-a65a0cf2c997/-50/available",
			nil,
		)
		require.NoError(t, err)

		token, err := auth.GenerateJWT(cfg, testUserID)
		require.NoError(t, err)

		req.Header.Add("Content-Type", "application/json")
		req.Header.Add("Authorization", "Bearer "+token)
		res, err := app.Test(req, -1)
		require.NoError(t, err)

		require.Equal(t, fiber.StatusBadRequest, res.StatusCode)
	})

	t.Run("should return 401 Unauthorized without token", func(t *testing.T) {
		req, err := http.NewRequest(
			http.MethodGet,
			"/abc/1/available",
			nil,
		)
		require.NoError(t, err)

		req.Header.Add("Content-Type", "application/json")
		req.Header.Add("Authorization", "Bearer "+"invalid")

		// No Authorization header added
		res, err := app.Test(req, -1)
		require.NoError(t, err)

		require.Equal(t, fiber.StatusUnauthorized, res.StatusCode)
	})
}
