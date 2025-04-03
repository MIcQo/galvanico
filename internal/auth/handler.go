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

func LoginHandler(ctx *fiber.Ctx) error {
	var cfg, cfgErr = config.Load()
	if cfgErr != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": cfgErr.Error()})
	}

	var req authRequest
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	var repo = user.NewUserRepository(database.Connection())
	var usr, err = repo.GetByUsername(ctx.Context(), req.Username)
	if err != nil {
		// if errors.Is(err, pgx.ErrNoRows) {
		if errors.Is(err, sql.ErrNoRows) {
			return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid credentials"})
		}
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	var jwt, jwtErr = GenerateJWT(cfg, usr.ID)
	if jwtErr != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": jwtErr.Error()})
	}

	return ctx.JSON(fiber.Map{
		"token": jwt,
	})
}

func ErrorHandler(c *fiber.Ctx, err error) error {
	if errors.Is(err, jwtware.ErrJWTMissingOrMalformed) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
	}
	return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "unauthorized"})
}
