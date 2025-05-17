package entity

import (
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

	Id         int64     `bun:"id,pk,autoincrement"`
	UserId     int64     `bun:"user_id,notnull"`
	WordId     int       `bun:"word_id,notnull"`
	GuessLimit int8      `bun:"guess_limit,notnull"`
	IsPlaying  bool      `bun:"is_playing,notnull,default:true"`
	IsWon      bool      `bun:"is_won,notnull"`
	CreatedAt  time.Time `bun:"created_at,notnull"`
	UpdatedAt  time.Time `bun:"updated_at,notnull"`

	Guesses []*Guess `bun:"rel:has-many,join:id=game_id"`
	Word    Word     `bun:"rel:belongs-to,join:word_id=id"`
}

type Guess struct {
	bun.BaseModel `bun:"table:guesses"`

	Id        int64     `bun:"id,pk,autoincrement"`
	GameId    int64     `bun:"game_id,notnull"`
	CreatedAt time.Time `bun:"created_at,notnull"`

	Letters []*Letter `bun:"rel:has-many,join:id=guess_id"`
}

type Letter struct {
	bun.BaseModel `bun:"table:letters"`

	Id                int64     `bun:"id,pk,autoincrement"`
	GuessId           int64     `bun:"guess_id,notnull"`
	Letter            string    `bun:"letter,notnull"`
	IsInWord          bool      `bun:"is_in_word,notnull"`
	IsCorrectPosition bool      `bun:"is_correct_position,notnull"`
	CreatedAt         time.Time `bun:"created_at,notnull"`
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

func (g *Game) Guess(guessWord string) *Guess {
	now := time.Now()
	guess := &Guess{GameId: g.Id, CreatedAt: now}
	for rPos, r := range []rune(guessWord) {
		l := &Letter{Letter: string(r), CreatedAt: now}
		for rwPos, rw := range []rune(g.Word.Word) {
			if rw == r {
				l.IsInWord = true
				if rwPos == rPos {
					l.IsCorrectPosition = true
				}
			}
		}

		guess.Letters = append(guess.Letters, l)
	}
	g.Guesses = append(g.Guesses, guess)

	return guess
}
