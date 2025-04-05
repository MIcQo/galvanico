package user

import (
	"database/sql"
	"errors"
	"galvanico/internal/auth"
	"galvanico/internal/config"
	"galvanico/internal/database"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"time"
)

type Handler struct {
	UserRepository Repository
}

func NewHandler(userRepository Repository) *Handler {
	return &Handler{UserRepository: userRepository}
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

// LoginHandler handles login request
func (h *Handler) LoginHandler(ctx *fiber.Ctx) error {
	var cfg, cfgErr = config.Load()
	if cfgErr != nil {
		return fiber.NewError(fiber.StatusInternalServerError, cfgErr.Error())
	}

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
		return ctx.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"message": "user is banned", "reason": usr.BanReason.String})
	}

	var updateErr = h.UserRepository.UpdateLastLogin(ctx.Context(), usr, ctx.IP())
	if updateErr != nil {
		return fiber.NewError(fiber.StatusInternalServerError, updateErr.Error())
	}

	var token, jwtErr = auth.GenerateJWT(cfg, usr.ID)
	if jwtErr != nil {
		return fiber.NewError(fiber.StatusInternalServerError, jwtErr.Error())
	}

	return ctx.JSON(fiber.Map{
		"token": token,
	})
}

func (h *Handler) RegisterHandler(_ *fiber.Ctx) error {
	return nil
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
