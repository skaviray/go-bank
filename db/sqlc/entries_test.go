package db

import (
	"context"
	"simple-bank/utils"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCreateEntries(t *testing.T) {
	// var id sql.NullInt64

	args := CreateEntryParams{
		AccountID: utils.RandomInt(1, 10),
		Amount:    10,
	}
	entry, err := testQueries.CreateEntry(context.TODO(), args)
	require.NoError(t, err)
	require.NotEmpty(t, entry)
	require.Equal(t, args.AccountID, entry.AccountID)
	require.Equal(t, args.Amount, entry.Amount)
	require.NotZero(t, entry.ID)
	require.NotZero(t, entry.CreatedAt)
}
