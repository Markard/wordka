package entity

import (
	"github.com/uptrace/bun"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type User struct {
	bun.BaseModel `bun:"table:users"`

	Id              int64        `bun:"id,pk,autoincrement"`
	Name            string       `bun:"name,notnull"`
	Email           string       `bun:"email,notnull,unique"`
	EmailVerifiedAt bun.NullTime `bun:"email_verified_at"`
	Password        string       `bun:"password,notnull"`
	CreatedAt       time.Time    `bun:"created_at,notnull"`
	UpdatedAt       time.Time    `bun:"updated_at,notnull"`
}

func NewUser(name string, email string, rawPassword string) (*User, error) {
	now := time.Now()
	password, err := hashPassword(rawPassword)
	if err != nil {
		return nil, err
	}

	return &User{
		Name:      name,
		Email:     email,
		Password:  password,
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}

func hashPassword(rawPassword string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(rawPassword), 12)
	return string(bytes), err
}
