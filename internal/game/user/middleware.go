package user

import (
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type Middleware struct {
	repo Repository
}

func NewUserMiddleware(repo Repository) *Middleware {
	return &Middleware{repo: repo}
}

// CheckNotBanned Middleware to check if the user is banned
func (m *Middleware) CheckNotBanned() fiber.Handler {
	return func(c *fiber.Ctx) error {
		token, ok := c.Locals("user").(*jwt.Token)
		if !ok || token == nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid token",
			})
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid claims",
			})
		}

		userID, ok := claims["sub"].(string)
		if !ok || userID == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Missing user ID in token",
			})
		}

		parsedID, err := uuid.Parse(userID)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid user ID format",
			})
		}

		user, err := m.repo.GetByID(c.Context(), parsedID)
		if err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "User not found",
			})
		}

		if user.BanExpiration.Valid {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "User is banned",
			})
		}

		c.Locals("user_data", user)
		return c.Next()
	}
}
