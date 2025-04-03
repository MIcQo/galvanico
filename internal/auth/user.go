package auth

import (
	"galvanico/internal/database"
	user2 "galvanico/internal/game/user"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func GetUser(c *fiber.Ctx) (*user2.User, error) {
	var user = c.Locals("user").(*jwt.Token)
	var claims = user.Claims.(jwt.MapClaims)
	var sub = uuid.MustParse(claims["sub"].(string))
	var repo = user2.NewUserRepository(database.Connection())
	var usr, err = repo.GetByID(c.Context(), sub)
	if err != nil {
		return nil, err
	}

	return usr, nil
}
