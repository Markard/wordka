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

	Id        int64     `bun:"id,pk,autoincrement"`
	UserId    int64     `bun:"user_id,notnull"`
	WordId    int       `bun:"word_id,notnull"`
	IsPlaying bool      `bun:"is_playing,notnull,default=true"`
	IsWon     bool      `bun:"is_won,notnull"`
	CreatedAt time.Time `bun:"created_at,notnull"`
	UpdatedAt time.Time `bun:"updated_at,notnull"`

	User *User `bun:"rel:belongs-to,join:user_id=id"`
	Word *Word `bun:"rel:belongs-to,join:word_id=id"`
}

type Guess struct {
	bun.BaseModel `bun:"table:guesses"`

	Id        int64     `bun:"id,pk,autoincrement"`
	GameId    int64     `bun:"game_id,notnull"`
	CreatedAt time.Time `bun:"created_at,notnull"`

	Game *Game `bun:"rel:belongs-to,join:game_id=id"`
}

type Letter struct {
	bun.BaseModel `bun:"table:letters"`

	Id                int64     `bun:"id,pk,autoincrement"`
	GuessId           int64     `bun:"guess_id,notnull"`
	IsInWord          bool      `bun:"is_in_word,notnull"`
	IsCorrectPosition bool      `bun:"is_correct_position,notnull"`
	CreatedAt         time.Time `bun:"created_at,notnull"`

	Guess *Guess `bun:"rel:belongs-to,join:guess_id=id"`
}
