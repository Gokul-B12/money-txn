package db

import (
	"context"
	"testing"
	"time"

	"github.com/Gokul-B12/money-txn/util"
	"github.com/stretchr/testify/require"
)

func CreateRandomEntry(t *testing.T, account1 Account) Entry {

	arg := CreateEntryParams{
		AccountID: account1.ID,
		Amount:    util.RandomBalance(),
	}

	entry, err := testQueries.CreateEntry(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, entry)
	require.NotEmpty(t, entry.ID)
	require.NotEmpty(t, entry.CreatedAt)
	require.Equal(t, entry.AccountID, arg.AccountID)
	require.Equal(t, entry.Amount, arg.Amount)

	return entry

}

func TestCreateEntry(t *testing.T) {
	account1 := CreateRandomAccount(t)
	CreateRandomEntry(t, account1)
}

func TestGetEntry(t *testing.T) {
	account1 := CreateRandomAccount(t)
	entry1 := CreateRandomEntry(t, account1)

	entry2, err := testQueries.GetEntry(context.Background(), entry1.ID)

	require.NoError(t, err)
	require.NotEmpty(t, entry2)
	require.Equal(t, entry1.AccountID, entry2.AccountID)
	require.Equal(t, entry1.ID, entry2.ID)
	require.Equal(t, entry1.Amount, entry2.Amount)
	require.Equal(t, entry1.CreatedAt, entry2.CreatedAt, time.Second)
	require.Equal(t, entry1, entry2)
}

func TestListEntries(t *testing.T) {
	account1 := CreateRandomAccount(t)
	for i := 0; i < 5; i++ {
		CreateRandomEntry(t, account1)
	}

	arg := ListEntriesParams{
		AccountID: account1.ID,
		Limit:     2,
		Offset:    3,
	}

	entries, err := testQueries.ListEntries(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, entries)
	require.Len(t, entries, 2)
	for _, entry := range entries {
		require.NotEmpty(t, entry)
		require.Equal(t, entry.AccountID, arg.AccountID)
	}

}
