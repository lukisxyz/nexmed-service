package model

import (
	"time"

	"github.com/oklog/ulid/v2"
	"gopkg.in/guregu/null.v4"
)

type Profile struct {
	Id          ulid.ULID  `json:"id"`
	AccountId   ulid.ULID  `json:"account_id"`
	FirstName   string     `json:"first_name"`
	LastName    string     `json:"last_name"`
	PhoneNumber string     `json:"phone_number"`
	Address     string     `json:"address"`
	City        string     `json:"city"`
	Province    string     `json:"province"`
	PostalCode  string     `json:"postal_code"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   null.Time  `json:"updated_at"`
}
