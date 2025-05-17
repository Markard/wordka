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
		return nil, errSelect
	}
	_ = tx.Commit()

	return game, nil
}

func (r *GameRepository) IsCurrentGameExists(currentUser *entity.User) (bool, error) {
	ctx := context.Background()

	isExists, errSelect := r.pgDb.NewSelect().
		Table("games").
		Where("user_id = ?", currentUser.Id).
		Where("is_playing = ?", true).
		Exists(ctx)

	return isExists, errSelect
}

func (r *GameRepository) CreateGame(word *entity.Word, currentUser *entity.User) (*entity.Game, error) {
	ctx := context.Background()
	game := entity.NewGame(word, currentUser)

	_, errInsert := r.pgDb.NewInsert().Model(game).Returning("id").Exec(ctx)
	if errInsert != nil {
		return nil, errInsert
	}

	return game, nil
}

func (r *GameRepository) FindRandomWord() (*entity.Word, error) {
	ctx := context.Background()
	word := entity.Word{}

	errSelect := r.pgDb.
		NewRaw("SELECT * FROM words ORDER BY RANDOM() LIMIT ?", 1).
		Scan(ctx, &word.Id, &word.Word, &word.CreatedAt)
	if errSelect != nil {
		return nil, errSelect
	}

	return &word, nil
}

func (r *GameRepository) FindWord(word string) (*entity.Word, error) {
	ctx := context.Background()
	w := entity.Word{}

	errSelect := r.pgDb.
		NewSelect().
		Model(&w).
		Where("word = ?", word).
		Scan(ctx)
	if errSelect != nil {
		return nil, errSelect
	}

	return &w, nil
}
