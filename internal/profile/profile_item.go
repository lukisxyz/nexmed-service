package profile

import (
	"errors"
	"time"

	"github.com/lukisxyz/nexmed-service/lib/db/model"
	"github.com/oklog/ulid/v2"
	"gopkg.in/guregu/null.v4"
)

var (
	ErrProfileNotFound = errors.New("profile: not found")
	ErrProfileConflict = errors.New("profile: already have profile")
)

func NewProfile(
	id string,
	fullName string,
	address string,
	phoneNumber string,
	birthDate null.Time,
	bio string,
) (model.Profile, error) {
	accountId, err := ulid.Parse(id)
	if (err != nil) {
		return model.Profile{}, err
	}
	profile := model.Profile{
		Id: ulid.Make(),
		AccountId: accountId,
		FullName: null.StringFrom(fullName),
		PhoneNumber: null.StringFrom(phoneNumber),
		BirthDate: birthDate,
		Bio: null.StringFrom(bio),
		Address: null.StringFrom(address),
		CreatedAt: time.Now(),
	}

	return profile, nil
}