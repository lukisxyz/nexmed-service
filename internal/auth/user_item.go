package auth

import (
	"errors"
	"time"

	"github.com/lukisxyz/nexmed-service/lib/db/model"
	"github.com/oklog/ulid/v2"
)

var (
	ErrEmailAlreadyUsed = errors.New("account: email already used")
	ErrAccountNotFound = errors.New("account: not found")
	ErrWrongPassword = errors.New("account: password not match")
	ErrResetPasswirdCannotSame = errors.New("account: new password cannot same")
	ErrEmailNotValid = errors.New("account: email not valid")
)

func NewAccount(email, password string) (model.Account, error) {
	hash, err := generateHash([]byte(password))
	if err != nil {
		return model.Account{}, err
	}
	account := model.Account{
		Id: ulid.Make(),
		Email: email,
		PasswordHash: hash,
		CreatedAt: time.Now(),
	}
	return account, nil
}