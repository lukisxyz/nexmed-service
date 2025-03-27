package model

import (
	"time"

	"github.com/oklog/ulid/v2"
	"gopkg.in/guregu/null.v4"
)

type Account struct {
	Id ulid.ULID `json:"id"`
	Email string `json:"email"`
	PasswordHash string `json:"password_hash"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt null.Time `json:"updated_at"`
}