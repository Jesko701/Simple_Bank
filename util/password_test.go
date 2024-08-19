package util

import (
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func TestGeneratePassword(t *testing.T) {
	password := RandomString(6)

	hashed_password, err := HashPassword(password)
	require.NoError(t, err)
	require.NotEmpty(t, hashed_password)
}

func TestComparePassword(t *testing.T) {
	password := RandomString(6)

	// Generate Password
	hashed_password, err := HashPassword(password)
	require.NoError(t, err)
	require.NotEmpty(t, hashed_password)

	// Check Valid Password
	err = CheckPassword(hashed_password, password)
	require.NoError(t, err)

	// Check Invalid Password
	wrongPassword := RandomString(6)
	err = CheckPassword(hashed_password, wrongPassword)
	require.Error(t, err)
	require.EqualError(t, err, bcrypt.ErrMismatchedHashAndPassword.Error())
	require.NotEmpty(t, wrongPassword)
	require.NotEqual(t, hashed_password, wrongPassword)
}
