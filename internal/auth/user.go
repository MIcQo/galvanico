package auth

import (
	"errors"
	"galvanico/internal/database"
	user2 "galvanico/internal/game/user"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func GetUser(c *fiber.Ctx) (*user2.User, error) {
	var user, userOk = c.Locals("user").(*jwt.Token)
	if !userOk {
		return nil, errors.New("user not authenticated")
	}

	var claims, claimOk = user.Claims.(jwt.MapClaims)
	if !claimOk {
		return nil, errors.New("invalid user claims")
	}

	var sub, ok = claims["sub"].(string)
	if !ok {
		return nil, errors.New("invalid user sub")
	}
	var uid = uuid.MustParse(sub)
	var repo = user2.NewUserRepository(database.Connection())
	var usr, err = repo.GetByID(c.Context(), uid)
	if err != nil {
		return nil, err
	}

	return usr, nil
}
