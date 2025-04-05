package auth

import (
	"errors"
	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
)

// ErrorHandler handles JWT unsuccessful request
func ErrorHandler(_ *fiber.Ctx, err error) error {
	if errors.Is(err, jwtware.ErrJWTMissingOrMalformed) {
		return fiber.NewError(fiber.StatusBadRequest, "invalid credentials")
	}
	return fiber.NewError(fiber.StatusUnauthorized, "unauthorized")
}
