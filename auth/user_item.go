package auth

import (
	"errors"
	"nexmedis-services/model"
	"time"

	"github.com/oklog/ulid/v2"
)

var (
	ErrEmailAlreadyUsed = errors.New("account: email already used")
	ErrAccountNotFound = errors.New("account: not found")
	ErrFailedCreateAccount = errors.New("account: failed create account, contact customer service")
)

func NewAccount(email, password string) (model.Account, error) {
	hash, err := generateHash([]byte(password), nil)
	if err != nil {
		return model.Account{}, err
	}
	account := model.Account{
		Id: ulid.Make(),
		Email: email,
		PasswordHash: string(hash.Hash),
		CreatedAt: time.Now(),
	}
	return account, nil
}