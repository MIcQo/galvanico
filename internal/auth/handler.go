package auth

import (
	"database/sql"
	"errors"
	"galvanico/internal/config"
	"galvanico/internal/database"
	"galvanico/internal/game/user"

	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
)

type authRequest struct {
	Username string `json:"username"`
}

// LoginHandler handles login request
func LoginHandler(ctx *fiber.Ctx) error {
	var cfg, cfgErr = config.Load()
	if cfgErr != nil {
		return fiber.NewError(fiber.StatusInternalServerError, cfgErr.Error())
	}

	var req authRequest
	if err := ctx.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	var repo = user.NewUserRepository(database.Connection())
	var usr, err = repo.GetByUsername(ctx.Context(), req.Username)
	if err != nil {
		// if errors.Is(err, pgx.ErrNoRows) {
		if errors.Is(err, sql.ErrNoRows) {
			return fiber.NewError(fiber.StatusUnauthorized, "invalid credentials")
		}
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	var jwt, jwtErr = GenerateJWT(cfg, usr.ID)
	if jwtErr != nil {
		return fiber.NewError(fiber.StatusInternalServerError, jwtErr.Error())
	}

	return ctx.JSON(fiber.Map{
		"token": jwt,
	})
}

// ErrorHandler handles JWT unsuccessful request
func ErrorHandler(_ *fiber.Ctx, err error) error {
	if errors.Is(err, jwtware.ErrJWTMissingOrMalformed) {
		return fiber.NewError(fiber.StatusBadRequest, "invalid credentials")
	}
	return fiber.NewError(fiber.StatusUnauthorized, "unauthorized")
}
