package repo

import (
	"context"
	"errors"
	"github.com/Markard/wordka/internal/entity"
	"github.com/jackc/pgerrcode"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/driver/pgdriver"
)

var ErrEmailUniqConstraint = errors.New("email already exists")

type AuthRepository struct {
	pgDb *bun.DB
}

func NewAuthRepository(pgDb *bun.DB) *AuthRepository {
	return &AuthRepository{pgDb: pgDb}
}

func (r AuthRepository) Create(user *entity.User) error {
	_, err := r.pgDb.NewInsert().Model(user).Returning("id").Exec(context.Background())
	if err != nil {
		var pgErr pgdriver.Error
		if errors.As(err, &pgErr) && pgErr.IntegrityViolation() && pgErr.Field('C') == pgerrcode.UniqueViolation {
			return ErrEmailUniqConstraint
		} else {
			return err
		}
	}

	return nil
}

func (r AuthRepository) FindBy(email string) (*entity.User, error) {
	user := &entity.User{}
	err := r.pgDb.NewSelect().
		Model(user).
		Where("email = ?", email).
		Scan(context.Background())

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r AuthRepository) FindById(id int64) (*entity.User, error) {
	user := &entity.User{}
	err := r.pgDb.NewSelect().
		Model(user).
		Where("id = ?", id).
		Scan(context.Background())

	if err != nil {
		return nil, err
	}

	return user, nil
}
