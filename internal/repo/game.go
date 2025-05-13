package repo

import (
	"context"
	"database/sql"
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
	ctx := context.Background()
	tx, err := r.pgDb.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelRepeatableRead})
	if err != nil {
		return nil, err
	}

	game := &entity.Game{}
	errSelect := tx.NewSelect().
		Model(game).
		Relation("Guesses").
		Relation("Guesses.Letters").
		Where("user_id = ?", currentUser.Id).
		Where("is_playing = ?", true).
		Scan(ctx)

	if errSelect != nil {
		_ = tx.Rollback()
		return nil, err
	}
	_ = tx.Commit()

	return game, nil
}
