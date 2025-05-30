package repo

import (
	"context"
	"database/sql"
	"errors"
	"github.com/Markard/wordka/internal/entity"
	"github.com/uptrace/bun"
	"time"
)

var (
	ErrCurrentGameNotFound = errors.New("current game not found")
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
	sq := tx.NewSelect()
	errSelect := getSelectQueryFindCurrentGame(sq, game, currentUser.Id).Scan(ctx)

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

func (r *GameRepository) AddGuessForCurrentGame(currentUser *entity.User, word *entity.Word) (*entity.Game, error) {
	ctx := context.Background()
	tx, err := r.pgDb.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelRepeatableRead})
	if err != nil {
		return nil, err
	}

	currentGame := &entity.Game{}
	sq := tx.NewSelect()
	errSelect := getSelectQueryFindCurrentGame(sq, currentGame, currentUser.Id).Scan(ctx)
	if errSelect != nil {
		_ = tx.Rollback()
		if currentGame.Id == 0 {
			return nil, ErrCurrentGameNotFound
		}
		return nil, errSelect
	}

	guess := currentGame.AddGuess(word)
	_, errInsert := tx.NewInsert().Model(guess).Returning("id").Exec(ctx)
	if errInsert != nil {
		_ = tx.Rollback()
		return nil, errInsert
	}

	if currentGame.IsPlaying == false {
		_, errUpdate := tx.NewUpdate().Model(currentGame).Where("id = ?", currentGame.Id).Exec(ctx)
		if errUpdate != nil {
			_ = tx.Rollback()
			return nil, errUpdate
		}
	}

	_ = tx.Commit()

	return currentGame, nil
}

func (r *GameRepository) SaveWords(words []string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	var wordEntities []*entity.Word
	for _, word := range words {
		wordEntities = append(wordEntities, entity.NewWord(word))
	}
	_, err := r.pgDb.NewInsert().
		Model(&wordEntities).
		Ignore().
		Exec(ctx)
	if err != nil {
		return err
	}

	return nil
}

func getSelectQueryFindCurrentGame(sq *bun.SelectQuery, model *entity.Game, userId int64) *bun.SelectQuery {
	sq.
		Model(model).
		Relation("Word").
		Relation("Guesses").
		Relation("Guesses.Word").
		Where("user_id = ?", userId).
		Where("is_playing = ?", true)
	return sq
}
