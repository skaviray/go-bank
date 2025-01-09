package db

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func executeTransaction(store Store, txName string, transferArgs TransferTXParams, results chan TransferTXResponse, errs chan error) {
	ctx := context.WithValue(context.Background(), TxKey("txname"), txName)
	transferResponse, err := store.TransferTX(ctx, transferArgs)
	errs <- err
	results <- transferResponse
}

func TestCreateTXN(t *testing.T) {
	// log.Println(testDb)
	store := NewStore(testDb)
	from_account := createRandomAccount(t)
	to_account := createRandomAccount(t)
	// fmt.Printf(">> Before Tx: from_account: %d, to_account: %d\n", from_account.Balance, to_account.Balance)
	transferArgs := TransferTXParams{
		FromAccount: from_account.ID,
		ToAccount:   to_account.ID,
		Amount:      int64(10),
	}
	n := 5
	errs := make(chan error)
	results := make(chan TransferTXResponse)

	for i := 0; i < n; i++ {
		txName := fmt.Sprintf("txn:%d", i+1)
		go executeTransaction(store, txName, transferArgs, results, errs)
	}
	existed_map := make(map[int]bool)
	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)
		result := <-results
		require.NotEmpty(t, result)
		transfer := result.Transfer
		require.NotEmpty(t, transfer)
		require.Equal(t, from_account.ID, transfer.FromAccount)
		require.Equal(t, to_account.ID, transfer.ToAccount)
		require.Equal(t, int64(10), transfer.Amount)
		require.NotZero(t, transfer.ID)
		require.NotZero(t, transfer.CreatedAt)
		_, err = store.GetTransfer(context.Background(), transfer.ID)
		require.NoError(t, err)

		// Check from entry
		fromEntry, err := store.GetEntry(context.Background(), result.FromEntry.ID)
		require.NoError(t, err)
		require.NotEmpty(t, fromEntry)
		require.Equal(t, from_account.ID, fromEntry.AccountID)
		require.Equal(t, -int64(10), fromEntry.Amount)
		require.NotZero(t, fromEntry.ID)
		require.NotZero(t, fromEntry.CreatedAt)

		// Check from entry
		toEntry, err := store.GetEntry(context.Background(), result.ToEntry.ID)
		require.NoError(t, err)
		require.NotEmpty(t, toEntry)
		require.Equal(t, to_account.ID, toEntry.AccountID)
		require.Equal(t, int64(10), toEntry.Amount)
		require.NotZero(t, toEntry.ID)
		require.NotZero(t, toEntry.CreatedAt)

		//check accounts

		require.NotEmpty(t, result.FromAccount)
		require.Equal(t, from_account.ID, result.FromAccount.ID)
		require.NotEmpty(t, result.ToAccount)
		require.Equal(t, to_account.ID, result.ToAccount.ID)

		// //check balances
		// fmt.Printf(">> After Tx: from_account: %d, to_account: %d\n", result.FromAccount.Balance, result.ToAccount.Balance)
		diff := from_account.Balance - result.FromAccount.Balance
		diff2 := result.ToAccount.Balance - to_account.Balance

		require.Equal(t, diff, diff2)
		require.True(t, diff > 0)
		require.True(t, diff%int64(10) == 0)

		k := int(diff / 10)
		// fmt.Println(k, n)
		require.True(t, k >= 1 && k <= n)
		require.NotContains(t, existed_map, k)
		existed_map[int(k)] = true
	}
	// check the final updated balance
	updatedBalance1, err := testQueries.GetAccount(context.Background(), from_account.ID)
	require.NoError(t, err)
	updatedBalance2, err := testQueries.GetAccount(context.Background(), to_account.ID)
	require.NoError(t, err)
	require.Equal(t, from_account.Balance-int64(n)*int64(10), updatedBalance1.Balance)
	require.Equal(t, to_account.Balance+int64(n)*int64(10), updatedBalance2.Balance)
}
