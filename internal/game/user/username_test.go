package user

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUserNameGenerator(t *testing.T) {
	var username, err = UsernameGenerator()
	require.NoError(t, err)
	assert.NotEmpty(t, username)
}
