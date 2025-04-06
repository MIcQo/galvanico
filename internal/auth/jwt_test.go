package auth

import (
	"galvanico/internal/config"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestGenerateJWT(t *testing.T) {
	var cfg = config.NewDefaultConfig()
	cfg.Auth.Settings["key"] = "token"

	var fixedUuid = uuid.MustParse("ebfc76b7-7ace-4034-a8b6-cc369afa8fb8")
	var token, err = GenerateJWT(cfg, fixedUuid)
	require.NoError(t, err)
	require.NotEmpty(t, token)
}
