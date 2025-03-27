package profile

import (
	"context"
	"errors"

	"github.com/lukisxyz/nexmed-service/lib/db/model"
	"github.com/oklog/ulid/v2"
	"github.com/rs/zerolog/log"
	"gopkg.in/guregu/null.v4"
)

func getProfile(
	ctx context.Context,
	id string,
) (model.Profile, error) {
	tx, err := pool.Begin(ctx)
	if err != nil {
		log.Error().Err(err).Msg("failed to pool postgres")
		return model.Profile{}, err
	}
	profile, err := getProfileByUserId(ctx, tx, id)
	tx.Commit(ctx)
	return profile, err
}

func updateProfile(
	ctx context.Context,
	id string,
	fullName string,
	address string,
	phoneNumber string,
	birthDate null.Time,
	bio string,
) (model.Profile, error) {
	tx, err := pool.Begin(ctx)
	if err != nil {
		log.Error().Err(err).Msg("failed to pool postgres")
		return model.Profile{}, err
	}
	profile, err := getProfileByUserId(
		ctx,
		tx,
		id,
	)

	if (err == nil) {
		accountId, err := ulid.Parse(id)
		if (err != nil) {
			return profile, ErrInternalServer
		}
		profile.Address = null.StringFrom(address)
		profile.Bio = null.StringFrom(bio)
		profile.BirthDate = null.Time(birthDate)
		profile.FullName = null.StringFrom(fullName)
		profile.AccountId = accountId
		profile.PhoneNumber = null.StringFrom(phoneNumber)
		err = saveProfile(ctx, tx, profile); 
		if err != nil {
			return profile, err
		}
		tx.Commit(ctx)
		return profile, nil
	}

	if (errors.Is(err, ErrProfileNotFound)) {
		tx.Rollback(ctx)
		return profile, ErrProfileNotFound
	}

	tx.Rollback(ctx)
	return profile, err
}

func createProfile(
	ctx context.Context,
	id string,
	fullName string,
	address string,
	phoneNumber string,
	birthDate null.Time,
	bio string,
) (model.Profile, error) {
	tx, err := pool.Begin(ctx)
	if err != nil {
		log.Error().Err(err).Msg("failed to pool postgres")
		return model.Profile{}, err
	}
	newProfile, err := getProfileByUserId(
		ctx,
		tx,
		id,
	)

	if (err == nil) {
		tx.Commit(ctx)
		return newProfile, ErrProfileConflict
	}

	if (errors.Is(err, ErrProfileNotFound)) {
		newProfile, err := NewProfile(
			id,
			fullName,
			address,
			phoneNumber,
			birthDate,
			bio,
		)
		if err != nil {
			tx.Rollback(ctx)
			return newProfile, err
		}
		err = saveProfile(ctx, tx, newProfile);
		if err != nil {
			return newProfile, err
		}
		tx.Commit(ctx)
		return newProfile, nil
	}

	return newProfile, err
}