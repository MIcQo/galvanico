package city

import (
	"context"
	"galvanico/internal/game/building"
	"galvanico/internal/game/user"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
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

	if len(cities) == 0 {
		var city, cityErr = h.service.InitCity(c.Context(), usr.ID)
		if cityErr != nil {
			return fiber.NewError(fiber.StatusInternalServerError, cityErr.Error())
		}

		cities = append(cities, city)
	}

	return c.JSON(cities)
}

type FakeRepository struct {
	data map[string]*City
}

func (f *FakeRepository) GetCitiesByUser(_ context.Context, userID uuid.UUID) ([]*City, error) {
	var cityID = uuid.Must(uuid.NewRandom())
	return []*City{
		{
			ID:        cityID,
			Name:      "City",
			PositionX: 1,
			PositionY: 1,
			UserCity: UserCity{
				UserID: userID,
				CityID: cityID,
			},
			Buildings: []Building{
				{
					CityID:   cityID,
					Building: building.CityHall,
					Level:    1,
					Position: 0,
				},
			},
			Resources: Resources{
				CityID: cityID,
			},
		},
	}, nil
}

func (f *FakeRepository) CreateCity(_ context.Context, city *City) error {
	cityID := city.ID.String()
	f.data[cityID] = city
	return nil
}

func NewFakeRepository() Repository {
	return &FakeRepository{data: make(map[string]*City)}
}
