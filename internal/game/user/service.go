package user

import (
	"context"
	"errors"
	"galvanico/internal/broker"
	"galvanico/internal/notifications"

	"golang.org/x/crypto/bcrypt"

	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

var (
	ErrAlreadyExists = errors.New("user already exists")
	ErrEmptyPassword = errors.New("empty password")
)

type Service interface {
	GetUser(ctx context.Context, token *jwt.Token) (*User, error)
	Register(ctx context.Context, usr *User) error
	SendActivationEmail(email *notifications.ActivationEmail) error
	SendPasswordWasChangedEmail(email *notifications.PasswordWasChanged) error
}

type ServiceIml struct {
	UserRepository Repository
	Publisher      broker.Publisher
}

func NewService(userRepository Repository, publisher broker.Publisher) Service {
	return &ServiceIml{UserRepository: userRepository, Publisher: publisher}
}

func (s *ServiceIml) GetUser(ctx context.Context, token *jwt.Token) (*User, error) {
	var claims, claimOk = token.Claims.(jwt.MapClaims)
	if !claimOk {
		return nil, errors.New("invalid user claims")
	}

	var sub, ok = claims["sub"].(string)
	if !ok {
		return nil, errors.New("invalid user sub")
	}
	var uid = uuid.MustParse(sub)
	var usr, err = s.UserRepository.GetByID(ctx, uid)
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

	if pubErr := s.Publisher.Publish("channels.email", u); pubErr != nil {
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

func (s *ServiceIml) Register(ctx context.Context, usr *User) error {
	var username, genErr = UsernameGenerator()
	if genErr != nil {
		return genErr
	}

	usr.Username = username

	if usr.Password.String == "" {
		return ErrEmptyPassword
	}

	var pass, cryptErr = bcrypt.GenerateFromPassword([]byte(usr.Password.String), bcrypt.DefaultCost)
	if cryptErr != nil {
		return cryptErr
	}

	usr.Password.String = string(pass)
	usr.Password.Valid = true

	if err := s.UserRepository.Create(ctx, usr); err != nil {
		if errors.Is(err, ErrDuplicateEntry) {
			return ErrAlreadyExists
		}
		return err
	}

	if err := s.SendActivationEmail(notifications.NewActivationEmail(usr.Email, usr.Username)); err != nil {
		return err
	}

	return nil
}

type FakeService struct {
	repo Repository
}

// GetUser this faked function mirrors real-one, cause of jwt
//
// Consider: Probably we can you other "verification" in the future
func (f *FakeService) GetUser(ctx context.Context, token *jwt.Token) (*User, error) {
	var claims, claimOk = token.Claims.(jwt.MapClaims)
	if !claimOk {
		return nil, errors.New("invalid user claims")
	}

	var sub, ok = claims["sub"].(string)
	if !ok {
		return nil, errors.New("invalid user sub")
	}

	var uid = uuid.MustParse(sub)
	var usr, err = f.repo.GetByID(ctx, uid)
	if err != nil {
		return nil, err
	}

	return usr, nil
}

// SendActivationEmail is not what we want in testing mode
func (f *FakeService) SendActivationEmail(_ *notifications.ActivationEmail) error {
	return nil
}

// SendPasswordWasChangedEmail is not what we want in testing mode
func (f *FakeService) SendPasswordWasChangedEmail(_ *notifications.PasswordWasChanged) error {
	return nil
}

func (f *FakeService) Register(_ context.Context, _ *User) error {
	return nil
}

func NewFakeService(repo Repository) Service {
	return &FakeService{repo: repo}
}
