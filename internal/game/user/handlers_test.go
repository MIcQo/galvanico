package user

import (
	"bytes"
	"context"
	"database/sql"
	"errors"
	"galvanico/internal/auth"
	"galvanico/internal/config"
	"io"
	"net/http"
	"testing"
	"time"

	jwtware "github.com/gofiber/contrib/jwt"

	"github.com/stretchr/testify/require"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

const signingKey = "AToAQz1ZtiDFPd6S5O4lyPCixPpo5I58"

type fakerUserRepository struct {
	data map[string]*User
}

func (f *fakerUserRepository) GetByUsername(_ context.Context, username string) (*User, error) {
	if usr, ok := f.data[username]; ok {
		return usr, nil
	}
	return nil, sql.ErrNoRows
}

func (f *fakerUserRepository) GetByID(_ context.Context, id uuid.UUID) (*User, error) {
	for _, usr := range f.data {
		if usr.ID.String() == id.String() {
			return usr, nil
		}
	}

	return nil, sql.ErrNoRows
}

func (f *fakerUserRepository) AddFeature(_ context.Context, _ *Feature) error {
	panic("implement me")
}

func (f *fakerUserRepository) RemoveFeature(_ context.Context, _ *Feature) error {
	panic("implement me")
}

func (f *fakerUserRepository) UpdateLastLogin(_ context.Context, _ *User, _ string) error {
	return nil
}

func (f *fakerUserRepository) ChangeUsername(_ context.Context, _ *User) error {
	return nil
}

func (f *fakerUserRepository) Create(_ context.Context, usr *User) error {
	if _, ok := f.data[usr.Username]; ok {
		return errors.New("username already exists")
	}

	f.data[usr.Username] = usr

	return nil
}

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
	var repo = &fakerUserRepository{data: map[string]*User{
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
	var svc = NewService(repo)
	var app = fiber.New()
	var handler = NewHandler(repo, svc, cfg)

	return app, handler
}

func TestHandler_LoginHandler(t *testing.T) {
	var app, handler = setup()

	app.Post("/auth/login", handler.LoginHandler)

	noArgsReq, _ := http.NewRequest(
		http.MethodPost,
		"/auth/login",
		nil,
	)

	noArgsRes, err := app.Test(noArgsReq, -1)

	require.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, noArgsRes.StatusCode)

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

	// //

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

	// //

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
}

func TestHandler_RegisterHandler(t *testing.T) {
	var app, handler = setup()
	app.Post("/auth/register", handler.RegisterHandler)

	noArgsReq, _ := http.NewRequest(
		http.MethodPost,
		"/auth/register",
		nil,
	)

	noArgsRes, err := app.Test(noArgsReq, -1)

	require.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, noArgsRes.StatusCode)

	// //
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

	// //
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
}

func TestHandler_ChangeUsernameHandler(t *testing.T) {
	var app, handler = setup()
	var cfg = config.NewDefaultConfig()
	cfg.Auth.Settings["key"] = signingKey
	app.Use(jwtware.New(jwtware.Config{
		SigningKey: jwtware.SigningKey{Key: []byte(signingKey), JWTAlg: jwtware.HS256},
	}))
	app.Patch("/api/user/username", handler.ChangeUsernameHandler)

	var reqBody, err = json.Marshal(usernameRequest{Username: gofakeit.Username()})
	require.NoError(t, err)

	req, _ := http.NewRequest(
		http.MethodPatch,
		"/api/user/username",
		bytes.NewReader(reqBody),
	)

	usr, err := handler.UserRepository.GetByUsername(t.Context(), "test")
	require.NoError(t, err)
	jwt, err := auth.GenerateJWT(cfg, usr.ID)
	require.NoError(t, err)

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+jwt)
	res, err := app.Test(req, -1)

	require.NoError(t, err)
	require.Equal(t, fiber.StatusOK, res.StatusCode)

	bodyBytes, err := io.ReadAll(res.Body)
	require.NoError(t, err)

	var body map[string]any
	err = json.Unmarshal(bodyBytes, &body)
	require.NoError(t, err)
}

func TestHandler_GetHandler(t *testing.T) {
	var app, handler = setup()
	var cfg = config.NewDefaultConfig()
	cfg.Auth.Settings["key"] = signingKey
	app.Use(jwtware.New(jwtware.Config{
		SigningKey: jwtware.SigningKey{Key: []byte(signingKey), JWTAlg: jwtware.HS256},
	}))
	app.Get("/api/user", handler.GetHandler)

	req, _ := http.NewRequest(
		http.MethodGet,
		"/api/user",
		nil,
	)

	usr, err := handler.UserRepository.GetByUsername(t.Context(), "test")
	require.NoError(t, err)
	jwt, err := auth.GenerateJWT(cfg, usr.ID)
	require.NoError(t, err)

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+jwt)
	res, err := app.Test(req, -1)

	require.NoError(t, err)
	require.Equal(t, fiber.StatusOK, res.StatusCode)

	bodyBytes, err := io.ReadAll(res.Body)
	require.NoError(t, err)

	var body map[string]any
	err = json.Unmarshal(bodyBytes, &body)
	require.NoError(t, err)
	assert.NotEmpty(t, body["user"])
}
