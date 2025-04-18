package city

import (
	"database/sql"
	"galvanico/internal/auth"
	"galvanico/internal/config"
	"galvanico/internal/game/user"
	"github.com/goccy/go-json"
	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
	"io"
	"net/http"
	"testing"
	"time"
)

const signingKey = "AToAQz1ZtiDFPd6S5O4lyPCixPpo5I58"

var testUserID = uuid.Must(uuid.NewRandom())

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
			ID:            uuid.Must(uuid.NewRandom()),
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
	})
}
