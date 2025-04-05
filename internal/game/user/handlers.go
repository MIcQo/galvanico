package user

import (
	"database/sql"
	"errors"
	"galvanico/internal/auth"
	"galvanico/internal/config"
	"galvanico/internal/database"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type Handler struct {
	Config         *config.Config
	UserRepository Repository
}

func NewHandler(userRepository Repository, cfg *config.Config) *Handler {
	return &Handler{UserRepository: userRepository, Config: cfg}
}

func (*Handler) GetHandler(c *fiber.Ctx) error {
	var usr, err = GetUser(c)
	if err != nil {
		return fiber.NewError(fiber.StatusUnauthorized, err.Error())
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"user": usr,
	})
}

type authRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type registerRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// LoginHandler handles login request
func (h *Handler) LoginHandler(ctx *fiber.Ctx) error {
	var req authRequest
	if err := ctx.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	var usr, err = h.UserRepository.GetByUsername(ctx.Context(), req.Username)
	if err != nil {
		// if errors.Is(err, pgx.ErrNoRows) {
		if errors.Is(err, sql.ErrNoRows) {
			return fiber.NewError(fiber.StatusUnauthorized, "invalid credentials")
		}
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	var authorizer, authErr = auth.Authenticate(auth.NewPasswordProvider(usr.Password.String, req.Password))
	if authErr != nil || !authorizer {
		return fiber.NewError(fiber.StatusForbidden, "invalid credentials")
	}

	if usr.BanExpiration.Valid && usr.BanExpiration.Time.After(time.Now().UTC()) {
		return ctx.Status(fiber.StatusUnprocessableEntity).
			JSON(fiber.Map{"message": "user is banned", "reason": usr.BanReason.String})
	}

	var updateErr = h.UserRepository.UpdateLastLogin(ctx.Context(), usr, ctx.IP())
	if updateErr != nil {
		return fiber.NewError(fiber.StatusInternalServerError, updateErr.Error())
	}

	var token, jwtErr = auth.GenerateJWT(h.Config, usr.ID)
	if jwtErr != nil {
		return fiber.NewError(fiber.StatusInternalServerError, jwtErr.Error())
	}

	return ctx.JSON(fiber.Map{
		"token": token,
	})
}

func (h *Handler) RegisterHandler(ctx *fiber.Ctx) error {
	var req registerRequest
	if err := ctx.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	var validate = validator.New()
	if err := validate.Struct(req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	var username, genErr = UsernameGenerator()
	if genErr != nil {
		return fiber.NewError(fiber.StatusInternalServerError, genErr.Error())
	}

	var usr = &User{
		ID:        uuid.New(),
		Email:     req.Email,
		Password:  sql.NullString{String: req.Password, Valid: true},
		Username:  username,
		Status:    "pending",
		Language:  "en",
		CreatedAt: time.Now().UTC(),
		Resources: Resources{
			Gold: DefaultUserGold,
		},
	}

	if err := h.UserRepository.Create(ctx.Context(), usr); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{})
}

func GetUser(c *fiber.Ctx) (*User, error) {
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
	var repo = NewUserRepository(database.Connection())
	var usr, err = repo.GetByID(c.Context(), uid)
	if err != nil {
		return nil, err
	}

	return usr, nil
}
