package profile

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/lukisxyz/nexmed-service/lib/db/model"
	"github.com/oklog/ulid/v2"
	"github.com/rs/zerolog/log"
)


func saveProfile(
	ctx context.Context,
	tx pgx.Tx,
	profile model.Profile,
) error {
	q := `
		INSERT INTO profiles (
			id,
			user_id,
			full_name,
			bio,
			phone_number,
			address,
			birth_date,
			created_at
		) VALUES (
			$1,
			$2,
			$3,
			$4,
			$5,
			$6,
			$7,
			NOW()
		) ON CONFLICT (user_id) DO UPDATE SET
			full_name = EXCLUDED.full_name,
			bio = EXCLUDED.bio,
			phone_number = EXCLUDED.phone_number,
			address = EXCLUDED.address,
			birth_date = EXCLUDED.birth_date,
			updated_at = NOW()
	`
	_, err := tx.Exec(
		ctx,
		q,
		ulid.MustParse(profile.Id.String()),
		ulid.MustParse(profile.AccountId.String()),
		profile.FullName,
		profile.Bio,
		profile.PhoneNumber,
		profile.Address,
		profile.BirthDate,
	)
	if err != nil {
		return err
	}
	return nil
}

func getProfileByUserId(
	ctx context.Context,
	tx pgx.Tx,
	id string,
) (model.Profile, error) {
	q := `
		SELECT id, user_id, full_name, phone_number, address, bio, avatar, birth_date, created_at, updated_at 
		FROM profiles WHERE user_id = $1
	`
	accountId, err := ulid.Parse(id)
	if err != nil {
		log.Debug().Err(err).Msg("failed parse ulid")
		return model.Profile{}, ErrInternalServer
	}
	row := tx.QueryRow(ctx, q, accountId)
	var profile model.Profile
	if err := row.Scan(
		&profile.Id,
		&profile.AccountId,
		&profile.FullName,
		&profile.PhoneNumber,
		&profile.Address,
		&profile.Bio,
		&profile.Avatar,
		&profile.BirthDate,
		&profile.CreatedAt,
		&profile.UpdatedAt,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			log.Debug().Err(err).Msg("can't find any item")
			return model.Profile{}, ErrProfileNotFound
		}
		return model.Profile{}, err
	}
	return profile, nil
}