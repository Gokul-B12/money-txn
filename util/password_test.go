package util

import (
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func TestHashPassword(t *testing.T) {

	password := RandomString(6)

	hashedPass1, err := HashPassword(password)
	require.NoError(t, err)
	require.NotEmpty(t, hashedPass1)

	input_pass1 := password
	err = CheckPassword(input_pass1, hashedPass1)
	require.NoError(t, err)

	wrongPassword := RandomString(6)
	err = CheckPassword(wrongPassword, hashedPass1)
	require.EqualError(t, err, bcrypt.ErrMismatchedHashAndPassword.Error())

	// we are checking whether for same passwords, HASH value is generating same or not.
	input_pass2 := password
	hashedPass2, err := HashPassword(input_pass2)
	require.NoError(t, err)
	require.NotEmpty(t, hashedPass2)
	require.NotEqual(t, hashedPass1, hashedPass2)

	//ANS: tests are passing... so, it is not generating same hash for same passwords. Because, it takes random SALT value. Ref: https://heynode.com/blog/2020-04/salt-and-hash-passwords-bcrypt/ and https://www.geeksforgeeks.org/what-is-salted-password-hashing/ or lec 17.

}
