package auth

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func TestPasswordProvider(t *testing.T) {
	var pass, err = bcrypt.GenerateFromPassword([]byte("test"), bcrypt.DefaultCost)
	require.NoError(t, err)

	t.Run("authenticate", func(t *testing.T) {
		var provider = NewPasswordProvider(string(pass), "test")
		assert.NotNil(t, provider)
		var auth, authErr = Authenticate(provider)
		require.NoError(t, authErr)
		assert.True(t, auth)
	})

	t.Run("password mismatch", func(t *testing.T) {
		var provider = NewPasswordProvider(string(pass), "invalid pass")
		assert.NotNil(t, provider)
		var auth, authErr = Authenticate(provider)
		require.Error(t, authErr)
		assert.False(t, auth)
	})
}
