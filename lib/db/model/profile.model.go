package model

import (
	"time"

	"github.com/oklog/ulid/v2"
	"gopkg.in/guregu/null.v4"
)

type Profile struct {
	Id          ulid.ULID  `json:"id"`
	AccountId   ulid.ULID  `json:"account_id"`
	FullName   null.String     `json:"full_name"`
	PhoneNumber null.String     `json:"phone_number"`
	Address     null.String     `json:"address"`
	Bio        null.String     `json:"bio"`
	Avatar        null.String     `json:"avatar"`
	BirthDate null.Time	`json:"birth_date"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   null.Time  `json:"updated_at"`
}
