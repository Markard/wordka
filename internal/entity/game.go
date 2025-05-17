package entity

import (
	"database/sql"
	"github.com/uptrace/bun"
	"time"
)

type Word struct {
	bun.BaseModel `bun:"table:words"`

	Id        int       `bun:"id,pk,autoincrement"`
	Word      string    `bun:"word,notnull"`
	CreatedAt time.Time `bun:"created_at,notnull"`
}

type Game struct {
	bun.BaseModel `bun:"table:games"`

	Id         int64        `bun:"id,pk,autoincrement"`
	UserId     int64        `bun:"user_id,notnull"`
	WordId     int          `bun:"word_id,notnull"`
	GuessLimit int8         `bun:"guess_limit,notnull"`
	IsPlaying  bool         `bun:"is_playing,notnull,default:true"`
	IsWon      sql.NullBool `bun:"is_won"`
	CreatedAt  time.Time    `bun:"created_at,notnull"`
	UpdatedAt  time.Time    `bun:"updated_at,notnull"`

	Guesses []*Guess `bun:"rel:has-many,join:id=game_id"`
	Word    *Word    `bun:"rel:belongs-to,join:word_id=id"`
}

type Guess struct {
	bun.BaseModel `bun:"table:guesses"`

	Id        int64     `bun:"id,pk,autoincrement"`
	GameId    int64     `bun:"game_id,notnull"`
	WordId    int       `bun:"word_id,notnull"`
	CreatedAt time.Time `bun:"created_at,notnull"`

	Word *Word `bun:"rel:belongs-to,join:word_id=id"`
}

func NewGame(word *Word, currentUser *User) *Game {
	now := time.Now()
	const guessLimit = 6

	return &Game{
		UserId:     currentUser.Id,
		WordId:     word.Id,
		GuessLimit: guessLimit,
		IsPlaying:  true,
		CreatedAt:  now,
		UpdatedAt:  now,
	}
}

func (g *Game) AddGuess(word *Word) *Guess {
	guess := &Guess{
		GameId:    g.Id,
		CreatedAt: time.Now(),
		WordId:    word.Id,
		Word:      word,
	}
	g.Guesses = append(g.Guesses, guess)

	if guess.WordId == g.WordId {
		g.IsPlaying = false
		g.IsWon.Bool = true
		g.IsWon.Valid = true
	} else if len(g.Guesses) >= int(g.GuessLimit) {
		g.IsPlaying = false
		g.IsWon.Bool = false
		g.IsWon.Valid = true
	}

	return guess
}

func (w *Word) AsRunes() []rune {
	return []rune(w.Word)
}
