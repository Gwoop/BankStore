/*
	Pattern: <filename>_test.go
*/

package db

import (
	"Bankstore/utils"
	"context"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/require"
)

func TestCreateAccount(t *testing.T) {
	createRandomAccount(t)
}

func createRandomAccount(t *testing.T) Account {
	user := createRandomUser(t)

	ra := utils.RandomAccount()
	arg := CreateAccountParams{
		Owner:   user.Username, // change because it's foreign key now
		Balance: utils.RandomInt(0, 1000),
		// Balance: ra.Balance,
		Currency: Currency(ra.Currency),
	}

	account, err := testQueries.CreateAccount(context.Background(), arg)

	// Две проверки на результат работы CreateAccount
	require.NoError(t, err)
	require.NotEmpty(t, account)

	// Дополнительные проверки
	require.Equal(t, arg.Owner, account.Owner)
	require.Equal(t, arg.Balance, account.Balance)
	require.Equal(t, arg.Currency, account.Currency)

	require.NotZero(t, account.ID)
	require.NotZero(t, account.CreatedAt)

	return account
}

func TestGetAccount(t *testing.T) {
	account1 := createRandomAccount(t)
	account2, err := testQueries.GetAccount(context.Background(), account1.ID)

	require.NoError(t, err)
	require.NotEmpty(t, account2)

	// test all fields
	require.Equal(t, account1.ID, account2.ID)
	require.Equal(t, account1.Owner, account2.Owner)
	require.Equal(t, account1.Balance, account2.Balance)
	require.Equal(t, account1.Currency, account2.Currency)
	require.WithinDuration(t, account1.CreatedAt.Time, account2.CreatedAt.Time, time.Second)
}

func TestUpdateAccount(t *testing.T) {
	account1 := createRandomAccount(t)

	arg := UpdateAccountParams{
		ID:      account1.ID,
		Balance: utils.RandomInt(0, 1000),
	}

	account2, err := testQueries.UpdateAccount(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, account2)

	// test all fields
	require.Equal(t, account1.ID, account2.ID)
	require.Equal(t, account1.Owner, account2.Owner)
	require.Equal(t, arg.Balance, account2.Balance)
	require.Equal(t, account1.Currency, account2.Currency)
	require.WithinDuration(t, account1.CreatedAt.Time, account2.CreatedAt.Time, time.Second)
}

func TestDeleteAccount(t *testing.T) {
	account1 := createRandomAccount(t)
	// Delete created account
	count, err := testQueries.DeleteAccount(context.Background(), account1.ID)
	require.NoError(t, err)
	require.Equal(t, count, int64(1))
	account2, err := testQueries.GetAccount(context.Background(), account1.ID)

	require.Error(t, err)
	require.Empty(t, account2)
	// проверка на тип ошибки(нуль строк было обработано)
	require.EqualError(t, err, pgx.ErrNoRows.Error())

}

func TestListAccounts(t *testing.T) {

	for range 10 {
		createRandomAccount(t)
	}
	var n int32 = 5
	arg := ListAccountsParams{
		Limit:  n,
		Offset: 5,
	}

	accounts, err := testQueries.ListAccounts(context.Background(), arg)

	require.NoError(t, err)
	require.Len(t, accounts, int(n))

	for _, acc := range accounts {
		require.NotZero(t, acc)

	}

}
