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
	Guesses []*Guess `json:"guesses"`
}

func NewResponse(game *entity.Game) *Response {
	guesses := make([]*Guess, 0)
	for _, g := range game.Guesses {
		var letters []*Letter
		for _, l := range g.Letters {
			letters = append(letters, &Letter{Letter: l.Letter, IsInWord: l.IsInWord, IsCorrectPosition: l.IsCorrectPosition})
		}
		guesses = append(guesses, &Guess{Letters: letters})
	}
	return &Response{Guesses: guesses}
}
