package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTransferTx(t *testing.T) {

	store := NewStore(testDB)

	account1 := CreateRandomAccount(t)
	account2 := CreateRandomAccount(t)

	//run n cuncurrent transfer transactions

	n := 5

	amount := int64(10)

	errs := make(chan error)
	results := make(chan TransferTxResult)

	for i := 0; i < n; i++ {
		go func() {
			result, err := store.TransferTx(context.Background(), TransferTxParams{
				FromAccountID: account1.ID,
				ToAccountID:   account2.ID,
				Amount:        amount,
			})

			errs <- err
			results <- result

		}()
	}

	// check results

	for i := 0; i < n; i++ {

		err := <-errs

		require.NoError(t, err)

		result := <-results

		require.NotEmpty(t, result)
		//checking  Transfer data
		transfer := result.Transfer

		require.NotEmpty(t, transfer)
		require.Equal(t, transfer.FromAccountID, account1.ID)
		require.Equal(t, transfer.ToAccountID, account2.ID)
		require.Equal(t, transfer.Amount, amount)
		require.NotEmpty(t, transfer.CreatedAt)
		require.NotEmpty(t, transfer.ID)

		//to check the transfer data is created or not... we will use GetTransfer()

		crtd_transfer, err := store.GetTransfer(context.Background(), transfer.ID)
		require.NotEmpty(t, crtd_transfer)
		require.NoError(t, err)

		//checking FromEntry data
		fromEntry := result.FromEntry
		require.NotEmpty(t, fromEntry)
		require.Equal(t, fromEntry.AccountID, account1.ID)
		require.NotZero(t, fromEntry.CreatedAt)
		require.NotEmpty(t, fromEntry.ID)
		require.Equal(t, fromEntry.Amount, -amount)

		//to check the FromEntry data is created or not... we will use GetEntry()

		crtdFrom_entry, err := store.GetEntry(context.Background(), fromEntry.ID)
		require.NotEmpty(t, crtdFrom_entry)
		require.NoError(t, err)

		//checking ToEntry data
		toEntry := result.ToEntry
		require.NotEmpty(t, toEntry)
		require.Equal(t, toEntry.AccountID, account2.ID)
		require.Equal(t, toEntry.Amount, amount)
		require.NotZero(t, toEntry.CreatedAt)
		require.NotEmpty(t, toEntry.ID)

		//to check the FromEntry data is created or not... we will use GetEntry()

		crtdTo_entry, err := store.GetEntry(context.Background(), toEntry.ID)
		require.NotEmpty(t, crtdTo_entry)
		require.NoError(t, err)

	}

}
