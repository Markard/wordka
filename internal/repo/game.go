package repo

import (
	"github.com/Markard/wordka/internal/entity"
	"github.com/uptrace/bun"
)

type GameRepository struct {
	pgDb *bun.DB
}

func NewGameRepository(pgDb *bun.DB) *GameRepository {
	return &GameRepository{pgDb: pgDb}
}

func (r *GameRepository) FindCurrentGame() (*entity.Game, error) {
	return nil, nil
}
