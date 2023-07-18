package db

import (
	"context"
	"fmt"
	"runtime"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTransferTx(t *testing.T) {

	store := NewStore(testDB)

	account1 := CreateRandomAccount(t)
	account2 := CreateRandomAccount(t)

	fmt.Println("Go-routines:", runtime.NumGoroutine())
	fmt.Println(">>before tx:", account1.Balance, " and ", account2.Balance)

	//run n cuncurrent transfer transactions

	n := 2

	amount := int64(10)

	errs := make(chan error)
	results := make(chan TransferTxResult)

	for i := 0; i < n; i++ {
		txName := fmt.Sprintf("tx %d ", i+1)
		go func() {
			ctx := context.WithValue(context.Background(), txKey, txName)
			result, err := store.TransferTx(ctx, TransferTxParams{
				FromAccountID: account1.ID,
				ToAccountID:   account2.ID,
				Amount:        amount,
			})

			errs <- err
			results <- result

		}()
		fmt.Println("Go-routines:", runtime.NumGoroutine())

	}
	fmt.Println("Go-routines:", runtime.NumGoroutine())

	// check results
	existed := make(map[int]bool)
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
		require.Equal(t, account2.ID, toEntry.AccountID)
		require.Equal(t, toEntry.Amount, amount)
		require.NotZero(t, toEntry.CreatedAt)
		require.NotEmpty(t, toEntry.ID)

		//to check the ToEntry data is created or not... we will use GetEntry()

		crtdTo_entry, err := store.GetEntry(context.Background(), toEntry.ID)
		require.NotEmpty(t, crtdTo_entry)
		require.NoError(t, err)

		//checking accounts

		fromAccount := result.FromAccount
		require.NotEmpty(t, fromAccount)
		require.Equal(t, account1.ID, fromAccount.ID)

		toAccount := result.ToAccount
		require.NotEmpty(t, toAccount)
		require.Equal(t, account2.ID, toAccount.ID)

		//check accounts balance
		fmt.Println(">> tx:", fromAccount.Balance, " and ", toAccount.Balance)
		diff1 := account1.Balance - fromAccount.Balance
		diff2 := toAccount.Balance - account2.Balance
		require.Equal(t, diff1, diff2)
		require.True(t, diff1 > 0)
		require.True(t, diff1%amount == 0) //amount. 2*amount. 3*amount, ..., n*amount

		k := int(diff1 / amount)

		require.True(t, k >= 1 && k <= n) //in order to check this we have used "existed" var
		existed[k] = true

		// check the final updated balance of two accounts

		updatedAccount1, err := store.GetAccount(context.Background(), account1.ID)
		require.NoError(t, err)
		updatedAccount2, err := store.GetAccount(context.Background(), account2.ID)
		require.NoError(t, err)

		require.Equal(t, account1.Balance-int64(n)*amount, updatedAccount1.Balance)
		require.Equal(t, account1.Balance+int64(n)*amount, updatedAccount2.Balance)

		fmt.Println(">> tx:", fromAccount.Balance, " and ", toAccount.Balance)
	}

}
