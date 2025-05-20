package user

import (
	"bytes"
	"database/sql"
	"galvanico/internal/auth"
	"galvanico/internal/config"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/goccy/go-json"
	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

const signingKey = "AToAQz1ZtiDFPd6S5O4lyPCixPpo5I58"

func setup() (*fiber.App, *Handler) {
	var pass, err = bcrypt.GenerateFromPassword([]byte("test"), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}

	var banTime, parseErr = time.Parse("2006-01-02 15:04:05", "2030-01-01 00:00:00")
	if parseErr != nil {
		panic(parseErr)
	}

	var cfg = config.NewDefaultConfig()
	var repo = &FakerUserRepository{data: map[string]*User{
		"test": {
			Username: "test",
			Password: sql.NullString{Valid: true, String: string(pass)},
			ID:       uuid.New(),
		},
		"banned": {
			Username:      "banned",
			Password:      sql.NullString{Valid: true, String: string(pass)},
			ID:            uuid.New(),
			BanExpiration: sql.NullTime{Time: banTime, Valid: true},
			BanReason:     sql.NullString{Valid: true, String: "banned"},
		},
	}}
	var svc = &FakeService{repo: repo}
	var app = fiber.New()
	var handler = NewHandler(repo, svc, cfg)

	return app, handler
}

func TestHandler_LoginHandler(t *testing.T) {
	var app, handler = setup()

	app.Post("/auth/login", handler.LoginHandler)

	t.Run("login with no args", func(t *testing.T) {
		noArgsReq, _ := http.NewRequest(
			http.MethodPost,
			"/auth/login",
			nil,
		)

		noArgsRes, err := app.Test(noArgsReq, -1)

		require.NoError(t, err)
		assert.Equal(t, fiber.StatusBadRequest, noArgsRes.StatusCode)
	})

	t.Run("success login", func(t *testing.T) {
		reqBody, err := json.Marshal(authRequest{
			Username: "test",
			Password: "test",
		})
		require.NoError(t, err)

		req, _ := http.NewRequest(
			http.MethodPost,
			"/auth/login",
			bytes.NewReader(reqBody),
		)
		req.Header.Add("Content-Type", "application/json")

		res, err := app.Test(req, -1)

		require.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, res.StatusCode)

		bodyBytes, err := io.ReadAll(res.Body)
		require.NoError(t, err)

		var body map[string]any
		err = json.Unmarshal(bodyBytes, &body)
		require.NoError(t, err)

		assert.NotEmpty(t, body["token"])
	})

	t.Run("invalid credentials", func(t *testing.T) {
		notFoundReqBody, err := json.Marshal(authRequest{
			Username: "notfound",
			Password: "test",
		})
		require.NoError(t, err)

		notFoundReq, _ := http.NewRequest(
			http.MethodPost,
			"/auth/login",
			bytes.NewReader(notFoundReqBody),
		)
		notFoundReq.Header.Add("Content-Type", "application/json")
		notFoundRes, err := app.Test(notFoundReq, -1)
		require.NoError(t, err)
		assert.Equal(t, fiber.StatusUnauthorized, notFoundRes.StatusCode)

		notFoundBodyBytes, err := io.ReadAll(notFoundRes.Body)
		require.NoError(t, err)

		var notFoundbody = string(notFoundBodyBytes)
		assert.Equal(t, "invalid credentials", notFoundbody)
	})

	t.Run("banned user", func(t *testing.T) {
		bannedReqBody, err := json.Marshal(authRequest{
			Username: "banned",
			Password: "test",
		})
		require.NoError(t, err)

		bannedReq, _ := http.NewRequest(
			http.MethodPost,
			"/auth/login",
			bytes.NewReader(bannedReqBody),
		)
		bannedReq.Header.Add("Content-Type", "application/json")
		bannedRes, err := app.Test(bannedReq, -1)
		require.NoError(t, err)
		assert.Equal(t, fiber.StatusUnprocessableEntity, bannedRes.StatusCode)

		bannedBodyBytes, err := io.ReadAll(bannedRes.Body)
		require.NoError(t, err)

		var bannedBody map[string]any
		err = json.Unmarshal(bannedBodyBytes, &bannedBody)
		require.NoError(t, err)

		assert.NotEmpty(t, bannedBody["message"])
		assert.NotEmpty(t, bannedBody["reason"])
		assert.Equal(t, "banned", bannedBody["reason"])
		assert.Equal(t, "user is banned", bannedBody["message"])
	})
}

func TestHandler_RegisterHandler(t *testing.T) {
	var app, handler = setup()
	app.Post("/auth/register", handler.RegisterHandler)

	t.Run("register with no args", func(t *testing.T) {
		noArgsReq, _ := http.NewRequest(
			http.MethodPost,
			"/auth/register",
			nil,
		)

		noArgsRes, err := app.Test(noArgsReq, -1)

		require.NoError(t, err)
		assert.Equal(t, fiber.StatusBadRequest, noArgsRes.StatusCode)
	})

	t.Run("register success", func(t *testing.T) {
		reqBody, err := json.Marshal(registerRequest{
			Email:    gofakeit.Email(),
			Password: gofakeit.Password(true, false, false, false, false, 10),
		})
		require.NoError(t, err)
		req, _ := http.NewRequest(
			http.MethodPost,
			"/auth/register",
			bytes.NewReader(reqBody),
		)
		req.Header.Add("Content-Type", "application/json")
		res, err := app.Test(req, -1)

		require.NoError(t, err)
		assert.Equal(t, fiber.StatusCreated, res.StatusCode)
	})

	t.Run("validation does not pass", func(t *testing.T) {
		validationReq, err := json.Marshal(registerRequest{
			Email:    gofakeit.Username(),
			Password: gofakeit.Password(true, false, false, false, false, 10),
		})
		require.NoError(t, err)
		invalidReq, _ := http.NewRequest(
			http.MethodPost,
			"/auth/register",
			bytes.NewReader(validationReq),
		)
		invalidReq.Header.Add("Content-Type", "application/json")
		invalidRes, err := app.Test(invalidReq, -1)

		require.NoError(t, err)
		assert.Equal(t, fiber.StatusBadRequest, invalidRes.StatusCode)
	})
}

func TestHandler_ChangeUsernameHandler(t *testing.T) {
	var app, handler = setup()
	var cfg = config.NewDefaultConfig()
	cfg.Auth.Settings["key"] = signingKey
	app.Use(jwtware.New(jwtware.Config{
		SigningKey: jwtware.SigningKey{Key: []byte(signingKey), JWTAlg: jwtware.HS256},
	}))
	app.Patch("/api/user/username", handler.ChangeUsernameHandler)

	t.Run("change username handler", func(t *testing.T) {
		var reqBody, err = json.Marshal(usernameRequest{Username: gofakeit.Username()})
		require.NoError(t, err)

		req, _ := http.NewRequest(
			http.MethodPatch,
			"/api/user/username",
			bytes.NewReader(reqBody),
		)

		usr, err := handler.UserRepository.GetByUsername(t.Context(), "test")
		require.NoError(t, err)
		token, err := auth.GenerateJWT(cfg, usr.ID)
		require.NoError(t, err)

		req.Header.Add("Content-Type", "application/json")
		req.Header.Add("Authorization", "Bearer "+token)
		res, err := app.Test(req, -1)

		require.NoError(t, err)
		require.Equal(t, fiber.StatusOK, res.StatusCode)

		bodyBytes, err := io.ReadAll(res.Body)
		require.NoError(t, err)

		var body map[string]any
		err = json.Unmarshal(bodyBytes, &body)
		require.NoError(t, err)
	})
}

func TestHandler_GetHandler(t *testing.T) {
	var app, handler = setup()
	var cfg = config.NewDefaultConfig()
	cfg.Auth.Settings["key"] = signingKey
	app.Use(jwtware.New(jwtware.Config{
		SigningKey: jwtware.SigningKey{Key: []byte(signingKey), JWTAlg: jwtware.HS256},
	}))
	app.Get("/api/user", handler.GetHandler)

	t.Run("get user", func(t *testing.T) {
		req, _ := http.NewRequest(
			http.MethodGet,
			"/api/user",
			nil,
		)

		usr, err := handler.UserRepository.GetByUsername(t.Context(), "test")
		require.NoError(t, err)
		token, err := auth.GenerateJWT(cfg, usr.ID)
		require.NoError(t, err)

		req.Header.Add("Content-Type", "application/json")
		req.Header.Add("Authorization", "Bearer "+token)
		res, err := app.Test(req, -1)

		require.NoError(t, err)
		require.Equal(t, fiber.StatusOK, res.StatusCode)

		bodyBytes, err := io.ReadAll(res.Body)
		require.NoError(t, err)

		var body map[string]any
		err = json.Unmarshal(bodyBytes, &body)
		require.NoError(t, err)
		assert.NotEmpty(t, body["user"])
	})
}

func TestHandler_ChangePasswordHandler(t *testing.T) {
	var app, handler = setup()
	var cfg = config.NewDefaultConfig()
	cfg.Auth.Settings["key"] = signingKey
	app.Use(jwtware.New(jwtware.Config{
		SigningKey: jwtware.SigningKey{Key: []byte(signingKey), JWTAlg: jwtware.HS256},
	}))
	app.Patch("/api/user/password", handler.ChangePasswordHandler)

	t.Run("change password handler", func(t *testing.T) {
		var reqBody, err = json.Marshal(changePasswordRequest{
			Password:    "test",
			NewPassword: gofakeit.Password(true, false, false, false, false, 10)},
		)
		require.NoError(t, err)

		req, _ := http.NewRequest(
			http.MethodPatch,
			"/api/user/password",
			bytes.NewReader(reqBody),
		)

		usr, err := handler.UserRepository.GetByUsername(t.Context(), "test")
		require.NoError(t, err)
		token, err := auth.GenerateJWT(cfg, usr.ID)
		require.NoError(t, err)

		req.Header.Add("Content-Type", "application/json")
		req.Header.Add("Authorization", "Bearer "+token)
		res, err := app.Test(req, -1)

		require.NoError(t, err)
		require.Equal(t, fiber.StatusOK, res.StatusCode)

		bodyBytes, err := io.ReadAll(res.Body)
		require.NoError(t, err)

		var body map[string]any
		err = json.Unmarshal(bodyBytes, &body)
		require.NoError(t, err)
	})
}
