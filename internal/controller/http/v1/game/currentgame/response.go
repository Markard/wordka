package currentgame

import (
	"github.com/Markard/wordka/internal/entity"
)

type Letter struct {
	Letter            string `json:"letter"`
	IsInWord          bool   `json:"is_in_word"`
	IsCorrectPosition bool   `json:"is_correct_position"`
}

type Guess struct {
	Letters []*Letter `json:"letters"`
}

type Response struct {
	IsWon   *bool    `json:"is_won"`
	Guesses []*Guess `json:"guesses"`
}

func NewResponse(game *entity.Game) *Response {
	guesses := make([]*Guess, 0)
	for _, g := range game.Guesses {
		guess := &Guess{}
		for wPos, wr := range g.Word.AsRunes() {
			l := &Letter{Letter: string(wr)}
			for swPos, swr := range game.Word.AsRunes() {
				if swr == wr {
					l.IsInWord = true
					if swPos == wPos {
						l.IsCorrectPosition = true
					}
				}
			}
			guess.Letters = append(guess.Letters, l)
		}
		guesses = append(guesses, guess)
	}

	resp := &Response{Guesses: guesses}
	if game.IsWon.Valid {
		resp.IsWon = &game.IsWon.Bool
	}

	return resp
}
