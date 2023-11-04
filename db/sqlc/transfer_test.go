package db

import (
	"context"
	"database/sql"
	"simplebank/util"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func CreateRandomTransfer(t *testing.T) Transfer {
	account1 := CreateRandomAccount(t)
	account2 := CreateRandomAccount(t)

	arg := CreateTransferParams{
		FromAccountsID: account1.ID,
		ToAccountsID:   account2.ID,
		Amount:         util.RandomMoney(),
	}

	transfer1, err := testQueries.CreateTransfer(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, transfer1)

	require.Equal(t, arg.FromAccountsID, transfer1.FromAccountsID)
	require.Equal(t, arg.ToAccountsID, transfer1.ToAccountsID)
	require.Equal(t, arg.Amount, transfer1.Amount)
	require.NotZero(t, transfer1.ID)
	require.NotZero(t, transfer1.CreatedAt)
	return transfer1
}

func TestCreateTransfer(t *testing.T) {
	CreateRandomTransfer(t)
}

func TestGetTransfer(t *testing.T) {
	transfer1 := CreateRandomTransfer(t)
	transfer2, err := testQueries.GetTransfer(context.Background(), transfer1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, transfer2)

	require.Equal(t, transfer1.ID, transfer2.ID)
	require.Equal(t, transfer1.FromAccountsID, transfer2.FromAccountsID)
	require.Equal(t, transfer1.ToAccountsID, transfer2.ToAccountsID)
	require.Equal(t, transfer1.Amount, transfer2.Amount)
	require.WithinDuration(t, transfer1.CreatedAt, transfer2.CreatedAt, time.Second)
}

func TestUpdateTransfer(t *testing.T) {
	transfer1 := CreateRandomTransfer(t)
	arg := UpdateTransferParams{
		ID:     transfer1.ID,
		Amount: util.RandomMoney(),
	}
	transfer2, err := testQueries.UpdateTransfer(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, transfer2)

	require.Equal(t, arg.ID, transfer2.ID)
	require.Equal(t, transfer1.FromAccountsID, transfer2.FromAccountsID)
	require.Equal(t, transfer1.ToAccountsID, transfer2.ToAccountsID)
	require.Equal(t, arg.Amount, transfer2.Amount)
	require.WithinDuration(t, transfer1.CreatedAt, transfer2.CreatedAt, time.Second)
}

func TestDeleteTransfer(t *testing.T) {
	transfer1 := CreateRandomTransfer(t)
	err := testQueries.DeleteTransfer(context.Background(), transfer1.ID)
	require.NoError(t, err)

	transfer2, err := testQueries.GetTransfer(context.Background(), transfer1.ID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())

	require.Empty(t, transfer2)
}

func TestListTransfers(t *testing.T) {
	for i := 0; i < 10; i++ {
		CreateRandomTransfer(t)
	}
	arg := ListTransfersParams{
		Limit:  5,
		Offset: 5,
	}
	transfers, err := testQueries.ListTransfers(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, transfers)
	for _, transfer := range transfers {
		require.NotEmpty(t, transfer)
	}
}
