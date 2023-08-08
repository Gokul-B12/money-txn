package db

import (
	"context"
	"testing"
	"time"

	"github.com/Gokul-B12/money-txn/util"
	"github.com/stretchr/testify/require"
)

func CreateRandomUser(t *testing.T) User {
	name := util.RandomOwner()
	arg := CreateUserParams{
		Username:       name,
		HashedPassword: "secret",
		FullName:       name,
		Email:          util.RandomEmail(),
	}

	user, err := testQueries.CreateUser(context.Background(), arg)

	//To check the test result, we are using the testify package..It’s more concise than just using the standard if else statements.
	// link: https://pkg.go.dev/github.com/stretchr/testify@v1.8.4/require and more details on https://pkg.go.dev/github.com/stretchr/testify

	require.NoError(t, err)

	require.Equal(t, arg.Username, user.Username)
	require.Equal(t, arg.HashedPassword, user.HashedPassword)
	require.Equal(t, arg.FullName, user.FullName)
	require.Equal(t, arg.Email, user.Email)

	require.NotEmpty(t, user.Username)
	require.NotZero(t, user.CreatedAt)
	require.True(t, user.PasswordChangedAt.IsZero())

	return user

}
func TestCreateUser(t *testing.T) {
	CreateRandomUser(t)
}

func TestGetUser(t *testing.T) {
	user1 := CreateRandomUser(t)

	user2, err := testQueries.GetUser(context.Background(), user1.Username)

	require.NoError(t, err)
	require.NotEmpty(t, user2)
	require.Equal(t, user1.Username, user2.Username)
	require.Equal(t, user1.HashedPassword, user2.HashedPassword)
	require.Equal(t, user1.FullName, user2.FullName)
	require.WithinDuration(t, user1.PasswordChangedAt, user2.PasswordChangedAt, time.Second)
	require.WithinDuration(t, user1.CreatedAt, user2.CreatedAt, time.Second)

}
