package city

import (
	"galvanico/internal/game/user"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

type Handler struct {
	repository  Repository
	service     Service
	userService user.Service
}

func NewHandler(repository Repository, service Service, userService user.Service) *Handler {
	return &Handler{repository: repository, service: service, userService: userService}
}

func (h *Handler) HandleGetUserCities(c *fiber.Ctx) error {
	var token, ok = c.Locals("user").(*jwt.Token)
	if !ok {
		return fiber.ErrUnauthorized
	}

	var usr, usrErr = h.userService.GetUser(c.Context(), token)
	if usrErr != nil {
		return fiber.NewError(fiber.StatusInternalServerError, usrErr.Error())
	}

	var cities, err = h.repository.GetCitiesByUser(c.Context(), usr.ID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(cities)
}
