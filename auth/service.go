package auth

import (
	"context"
	"errors"
)

func createAccount(
	ctx context.Context,
	email string,
	password string,
) error {
	account, err := NewAccount(email, password)
	if err != nil {
		return err
	}
	tx, err := pool.Begin(ctx)
	_, err = getUserByEmail(ctx, tx, account.Email)
	if (err != nil) {
		tx.Rollback(ctx)
		return ErrFailedCreateAccount
	}
	if !(errors.Is(err, ErrAccountNotFound)) {
		tx.Rollback(ctx)
		return err
	}
	err = saveAccount(ctx, tx, account)
	return nil
}