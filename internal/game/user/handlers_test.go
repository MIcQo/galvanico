package user

import (
	"bytes"
	"context"
	"database/sql"
	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
	"io"
	"net/http"
	"testing"
	"time"
)

type fakerUserRepository struct {
	mock.Mock
	data map[string]*User
}

func (f *fakerUserRepository) GetByUsername(_ context.Context, username string) (*User, error) {
	if usr, ok := f.data[username]; ok {
		return usr, nil
	}
	return nil, sql.ErrNoRows
}

func (f *fakerUserRepository) GetByID(_ context.Context, _ uuid.UUID) (*User, error) {
	panic("implement me")
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

func setup() (*fiber.App, *Handler) {
	var pass, err = bcrypt.GenerateFromPassword([]byte("test"), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}

	var banTime, parseErr = time.Parse("2006-01-02 15:04:05", "2030-01-01 00:00:00")
	if parseErr != nil {
		panic(parseErr)
	}

	var app = fiber.New()
	var handler = NewHandler(&fakerUserRepository{data: map[string]*User{
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
	}})

	return app, handler
}

func TestHandler_LoginHandler(t *testing.T) {
	var app, handler = setup()

	app.Post("/auth/login", handler.LoginHandler)

	noArgsReq, _ := http.NewRequest(
		"POST",
		"/auth/login",
		nil,
	)

	noArgsRes, err := app.Test(noArgsReq, -1)

	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, noArgsRes.StatusCode)

	reqBody, err := json.Marshal(authRequest{
		Username: "test",
		Password: "test",
	})
	assert.NoError(t, err)

	req, _ := http.NewRequest(
		"POST",
		"/auth/login",
		bytes.NewReader(reqBody),
	)
	req.Header.Add("Content-Type", "application/json")

	res, err := app.Test(req, -1)

	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, res.StatusCode)

	bodyBytes, err := io.ReadAll(res.Body)
	assert.Nil(t, err)

	var body map[string]any
	err = json.Unmarshal(bodyBytes, &body)
	assert.Nil(t, err)

	assert.NotEmpty(t, body["token"])

	// //

	notFoundReqBody, err := json.Marshal(authRequest{
		Username: "notfound",
		Password: "test",
	})
	assert.NoError(t, err)

	notFoundReq, _ := http.NewRequest(
		"POST",
		"/auth/login",
		bytes.NewReader(notFoundReqBody),
	)
	notFoundReq.Header.Add("Content-Type", "application/json")
	notFoundRes, err := app.Test(notFoundReq, -1)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusUnauthorized, notFoundRes.StatusCode)

	notFoundBodyBytes, err := io.ReadAll(notFoundRes.Body)
	assert.Nil(t, err)

	var notFoundbody = string(notFoundBodyBytes)
	assert.Equal(t, "invalid credentials", notFoundbody)

	// //

	bannedReqBody, err := json.Marshal(authRequest{
		Username: "banned",
		Password: "test",
	})
	assert.NoError(t, err)

	bannedReq, _ := http.NewRequest(
		"POST",
		"/auth/login",
		bytes.NewReader(bannedReqBody),
	)
	bannedReq.Header.Add("Content-Type", "application/json")
	bannedRes, err := app.Test(bannedReq, -1)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusUnprocessableEntity, bannedRes.StatusCode)

	bannedBodyBytes, err := io.ReadAll(bannedRes.Body)
	assert.Nil(t, err)

	var bannedBody map[string]any
	err = json.Unmarshal(bannedBodyBytes, &bannedBody)
	assert.Nil(t, err)

	assert.NotEmpty(t, bannedBody["message"])
	assert.NotEmpty(t, bannedBody["reason"])
	assert.Equal(t, "banned", bannedBody["reason"])
	assert.Equal(t, "user is banned", bannedBody["message"])
}
