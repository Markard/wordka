package repo

import (
	"context"
	"github.com/Markard/wordka/internal/entity"
	"github.com/uptrace/bun"
)

type GameRepository struct {
	pgDb *bun.DB
}

func NewGameRepository(pgDb *bun.DB) *GameRepository {
	return &GameRepository{pgDb: pgDb}
}

func (r *GameRepository) FindCurrentGame(currentUser *entity.User) (*entity.Game, error) {
	game := &entity.Game{}
	err := r.pgDb.NewSelect().
		Model(game).
		Relation("Guesses").
		Relation("Guesses.Letters").
		Where("user_id = ?", currentUser.Id).
		Where("is_playing = ?", true).
		Scan(context.Background())

	if err != nil {
		return nil, err
	}

	return game, nil
}
