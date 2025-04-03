package auth

import (
	"galvanico/internal/config"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

const tokenExpiration = time.Hour * 24

func GenerateJWT(cfg *config.Config, id uuid.UUID) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"sub": id.String(),
			"exp": time.Now().Add(tokenExpiration).Unix(),
			"iat": time.Now().Unix(),
		})

	return token.SignedString(cfg.Auth.GetJWTKey())
}
