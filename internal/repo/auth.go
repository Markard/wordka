package repo

import (
	"context"
	"github.com/Markard/wordka/internal/entity"
	"github.com/uptrace/bun"
)

type AuthRepository struct {
	pgDb *bun.DB
}

func NewAuthRepository(pgDb *bun.DB) *AuthRepository {
	return &AuthRepository{pgDb: pgDb}
}

func (r AuthRepository) Create(user *entity.User) error {
	_, err := r.pgDb.NewInsert().Model(user).Returning("id").Exec(context.Background())
	if err != nil {
		return err
	}

	return nil
}
