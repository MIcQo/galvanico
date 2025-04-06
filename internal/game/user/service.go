package user

import (
	"errors"
	"galvanico/internal/broker"
	"galvanico/internal/notifications"

	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type Service interface {
	GetUser(c *fiber.Ctx) (*User, error)
	SendActivationEmail(email *notifications.ActivationEmail) error
	SendPasswordWasChangedEmail(email *notifications.PasswordWasChanged) error
}

type ServiceIml struct {
	UserRepository Repository
}

func NewService(userRepository Repository) Service {
	return &ServiceIml{UserRepository: userRepository}
}

func (s *ServiceIml) GetUser(c *fiber.Ctx) (*User, error) {
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
	var usr, err = s.UserRepository.GetByID(c.Context(), uid)
	if err != nil {
		return nil, err
	}

	return usr, nil
}

func (s *ServiceIml) SendActivationEmail(email *notifications.ActivationEmail) error {
	var msg = notifications.NewMessage(notifications.ChannelEmail, notifications.TypeActivationEmail, email)
	var u, err = json.Marshal(msg)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	if pubErr := broker.Connection().Publish("channels.email", u); pubErr != nil {
		return fiber.NewError(fiber.StatusInternalServerError, pubErr.Error())
	}

	return nil
}

func (s *ServiceIml) SendPasswordWasChangedEmail(email *notifications.PasswordWasChanged) error {
	var msg = notifications.NewMessage(notifications.ChannelEmail, notifications.TypePasswordChanged, email)
	var u, err = json.Marshal(msg)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	if pubErr := broker.Connection().Publish("channels.email", u); pubErr != nil {
		return fiber.NewError(fiber.StatusInternalServerError, pubErr.Error())
	}

	return nil
}
