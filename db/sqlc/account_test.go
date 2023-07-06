package db

import (
	"context"
	"testing"

	"github.com/Gokul-B12/simplebank/util"

	"github.com/stretchr/testify/require"
)

func TestCreateAccount(t *testing.T) {
	arg := CreateAccountParams{
		Owner:    util.RandomOwner(),
		Balance:  util.RandomBalance(),
		Currency: util.RandomCurrency(),
	}

	account, err := testQueries.CreateAccount(context.Background(), arg)

	//To check the test result, we are using the testify package..It’s more concise than just using the standard if else statements.
	// link: https://pkg.go.dev/github.com/stretchr/testify@v1.8.4/require and more details on https://pkg.go.dev/github.com/stretchr/testify

	require.NoError(t, err)

	require.Equal(t, arg.Owner, account.Owner)
	require.Equal(t, arg.Balance, account.Balance)
	require.Equal(t, arg.Currency, account.Currency)

	require.NotZero(t, account.ID)
	require.NotZero(t, account.CreatedAt)

}
