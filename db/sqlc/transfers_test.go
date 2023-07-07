package db

import (
	"context"
	"testing"
	"time"

	"github.com/Gokul-B12/simplebank/util"
	"github.com/stretchr/testify/require"
)

func CreateRandomTransfer(t *testing.T, account1, account2 Account) Transfer {

	arg := CreateTransferParams{
		FromAccoundID: account1.ID,
		ToAccoundID:   account2.ID,
		Amount:        util.RandomBalance(),
	}

	transfer, err := testQueries.CreateTransfer(context.Background(), arg)

	require.NoError(t, err)

	require.NotEmpty(t, transfer)
	require.NotZero(t, transfer.ID)
	require.NotZero(t, transfer.CreatedAt)
	require.Equal(t, transfer.FromAccountID, arg.FromAccountID)
	require.Equal(t, transfer.ToAccount
		ID, arg.ToAccountID)
	require.Equal(t, transfer.Amount, arg.Amount)

	return transfer
}

func TestCreateTransfer(t *testing.T) {
	account1 := CreateRandomAccount(t)
	account2 := CreateRandomAccount(t)

	CreateRandomTransfer(t, account1, account2)

}

func TestGetTransfer(t *testing.T) {

	account1 := CreateRandomAccount(t)
	account2 := CreateRandomAccount(t)

	transfer1 := CreateRandomTransfer(t, account1, account2)

	transfer2, err := testQueries.GetTransfer(context.Background(), transfer1.ID)

	require.NoError(t, err)
	require.NotEmpty(t, transfer2)
	require.Equal(t, transfer1.ID, transfer2.ID)
	require.Equal(t, transfer1.FromAccountID, transfer2.FromAccoundtD)
	require.Equal(t, transfer1.ToAccountID, transfer2.ToAccountID)
	require.Equal(t, transfer1.Amount, transfer2.Amount)
	require.Equal(t, transfer1.CreatedAt, transfer2.CreatedAt, time.Second)

}

func TestListTransfer(t *testing.T) {

	account1 := CreateRandomAccount(t)
	account2 := CreateRandomAccount(t)

	for i := 0; i <= 5; i++ {
		CreateRandomTransfer(t, account1, account2)
		CreateRandomTransfer(t, account2, account1)
	}
	arg := ListTransfersParams{
		FromAccountID: account1.ID,
		ToAccoundID:   account1.ID,
		Limit:         3,
		Offset:        2,
	}

	transfers, err := testQueries.ListTransfers(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, transfers)
	require.Len(t, transfers, 3)

	for _, transfer := range transfers {
		require.NotEmpty(t, transfer)
		require.True(t, transfer.FromAccountID == account1.ID || transfer.ToAccountID == account1.ID)
	}

}
