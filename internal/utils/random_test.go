package utils

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRandomIndex(t *testing.T) {
	var s = []string{
		"foo",
		"bar",
		"baz",
	}

	var r, err = RandomIndex(s)
	require.NoError(t, err)
	require.NotEmpty(t, r)
}

func TestRandomNumber(t *testing.T) {
	var n, err = RandomNumber(10, 50)
	require.NoError(t, err)
	require.Greater(t, n, int64(10))
	require.Less(t, n, int64(50))
}

func TestRandomString(t *testing.T) {
	var str = RandomString(10)
	require.NotEmpty(t, str)
	require.Len(t, str, 10)
}
