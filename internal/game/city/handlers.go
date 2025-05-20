package city

import (
	"galvanico/internal/game/building"
	"galvanico/internal/game/user"
	"galvanico/internal/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

const PortSlot1 = 4
const PortSlot2 = 5
const DefenseSlot1 = 8

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
		return fiber.NewError(fiber.StatusForbidden, usrErr.Error())
	}

	if usr.BanExpiration.Valid {
		return fiber.NewError(fiber.StatusForbidden, "User banned")
	}

	var cities, err = h.repository.GetCitiesByUser(c.Context(), usr.ID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	if len(cities) == 0 {
		var city, cityErr = h.service.InitCity(c.Context(), usr.ID)
		if cityErr != nil {
			return fiber.NewError(fiber.StatusInternalServerError, cityErr.Error())
		}

		cities = append(cities, city)
	}

	return c.JSON(cities)
}

func (h *Handler) HandleAvailableSlotBuildings(c *fiber.Ctx) error {
	var token, ok = c.Locals("user").(*jwt.Token)
	if !ok {
		return fiber.ErrUnauthorized
	}

	var cityID = c.Params("city", "")
	if cityID == "" {
		return fiber.NewError(fiber.StatusBadRequest, "city ID is required")
	}

	var slot, err = c.ParamsInt("slot", 0)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	// Validate slot number is within valid range (assuming valid slots are 1-10)
	if slot < 1 || slot > 10 {
		return fiber.NewError(fiber.StatusBadRequest, "invalid slot number")
	}

	var city, cityErr = h.repository.GetCityByID(c.Context(), uuid.MustParse(cityID))
	if cityErr != nil {
		return fiber.NewError(fiber.StatusInternalServerError, cityErr.Error())
	}

	var usr, _ = h.userService.GetUser(c.Context(), token)
	if city.UserCity.UserID != usr.ID {
		return fiber.ErrForbidden
	}

	var alreadyBuilt []building.Building
	for _, b := range city.Buildings {
		alreadyBuilt = append(alreadyBuilt, b.Building)
	}

	var buildings []building.Building

	switch slot {
	case PortSlot1, PortSlot2:
		buildings = building.GetPortBuildings()

	case DefenseSlot1:
		buildings = building.GetDefenseBuildings()
	default:
		buildings = building.GetStandardBuildings()
	}

	return c.JSON(fiber.Map{
		"slot":      slot,
		"buildings": utils.FilterNotIn(buildings, alreadyBuilt),
	})
}
