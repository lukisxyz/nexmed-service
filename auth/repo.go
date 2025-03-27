package auth

import (
	"context"
	"errors"
	"nexmedis-services/model"

	"github.com/jackc/pgx/v5"
	"github.com/rs/zerolog/log"
)

func saveAccount(
	ctx context.Context,
	tx pgx.Tx,
	account model.Account,
) error {
	q := `
		INSERT INTO users(
			id,
			email,
			password_hash,
			created_at
		) VALUES (
			$1,
			$2,
			$3,
			$4
		)
		ON CONFLICT(email) DO UPDATE SET
			password_hash = $3,
			updated_at = $4
	`
	_, err := tx.Exec(
		ctx,
		q,
		account.Id,
		account.Email,
		account.PasswordHash,
		account.CreatedAt,
	)
	if err != nil {
		return err
	}
	return nil
}

func getUserByEmail(
	ctx context.Context,
	tx pgx.Tx,
	email string,
) (model.Account, error) {
	q := `
		FROM users
		|> WHERE email = $1
	`
	row := tx.QueryRow(ctx, q, email)
	var account model.Account
	if err := row.Scan(
		&account.Id,
		&account.Email,
		&account.PasswordHash,
		&account.CreatedAt,
		&account.UpdatedAt,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			log.Debug().Err(err).Msg("can't find any item")
			return model.Account{}, ErrAccountNotFound
		}
		return model.Account{}, err
	}
	return account, nil
}