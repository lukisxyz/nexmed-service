package auth

import (
	"context"
	"errors"

	"github.com/rs/zerolog/log"
)

func changePassword(
	ctx context.Context,
	email string,
	current_password string,
	new_password string,
) error {
	tx, err := pool.Begin(ctx)
	if err != nil {
		log.Error().Err(err).Msg("failed to pool postgres")
		return err
	}
	account, err := getUserByEmail(ctx, tx, email)
	if (err != nil) {
		tx.Rollback(ctx)
		return ErrAccountNotFound
	}

	if !compare([]byte(account.PasswordHash), []byte(current_password)) {
		tx.Rollback(ctx)
		return ErrWrongPassword
	}

	if new_password == current_password {
		tx.Rollback(ctx)
		return ErrResetPasswirdCannotSame
	}

	hash, err := generateHash([]byte(new_password))
	if err != nil {
		tx.Rollback(ctx)
		return ErrInternalServer
	}
	account.PasswordHash = hash

	err = saveAccount(ctx, tx, account)
	tx.Commit(ctx)
	return err
}

func createAccount(
	ctx context.Context,
	email string,
	password string,
) error {
	account, err := NewAccount(email, password)
	if err != nil {
		log.Error().Err(err)
		return err
	}
	tx, err := pool.Begin(ctx)
	if err != nil {
		log.Error().Err(err).Msg("failed to pool postgres")
		return ErrInternalServer
	}
	_, err = getUserByEmail(ctx, tx, account.Email)
	if (err == nil) {
		tx.Rollback(ctx)
		return ErrEmailAlreadyUsed
	}
	if !(errors.Is(err, ErrAccountNotFound)) {
		log.Info().Err(err)
		tx.Rollback(ctx)
		return err
	}
	err = saveAccount(ctx, tx, account)
	tx.Commit(ctx)
	return err
}

func loginAccount(
	ctx context.Context,
	email string,
	password string,
) (string, error) {
	tx, err := pool.Begin(ctx)
	if err != nil {
		log.Error().Err(err).Msg("failed to pool postgres")
		return "", ErrInternalServer
	}
	account, err := getUserByEmail(ctx, tx, email)
	if (err != nil) {
		tx.Rollback(ctx)
		return "", err
	}

	if !compare([]byte(account.PasswordHash), []byte(password)) {
		tx.Rollback(ctx)
		return "", ErrWrongPassword
	}

	tx.Commit(ctx)
	return account.Id.String(), nil
}