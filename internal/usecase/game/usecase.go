package game

import (
	"github.com/Markard/wordka/internal/entity"
)

type IGameRepository interface {
	FindCurrentGame(currentUser *entity.User) (*entity.Game, error)
	IsCurrentGameExists(currentUser *entity.User) (bool, error)
	CreateGame(word *entity.Word, currentUser *entity.User) (*entity.Game, error)
	FindRandomWord() (*entity.Word, error)
	FindWord(word string) (*entity.Word, error)
	AddGuessForCurrentGame(user *entity.User, word *entity.Word) (*entity.Game, error)
}

type UseCase struct {
	repository IGameRepository
}

func NewGameUseCase(repository IGameRepository) *UseCase {
	return &UseCase{repository: repository}
}

func (p *UseCase) FindCurrentGame(user *entity.User) (*entity.Game, error) {
	game, err := p.repository.FindCurrentGame(user)
	if err != nil {
		if game == nil {
			return nil, ErrCurrentGameNotFound{}
		} else {
			return nil, err
		}
	}

	return game, nil
}

func (p *UseCase) CreateGame(user *entity.User) (*entity.Game, error) {
	isExists, err := p.repository.IsCurrentGameExists(user)
	if err != nil {
		return nil, err
	}

	if isExists {
		return nil, ErrCurrentGameAlreadyExists{}
	}

	randomWord, err := p.repository.FindRandomWord()
	if err != nil {
		return nil, err
	}

	game, err := p.repository.CreateGame(randomWord, user)
	if err != nil {
		return nil, err
	}

	return game, nil
}

func (p *UseCase) Guess(user *entity.User, wordStr string) (*entity.Game, error) {
	word, _ := p.repository.FindWord(wordStr)
	if word == nil {
		return nil, ErrIncorrectWord{}
	}

	game, errAddGuess := p.repository.AddGuessForCurrentGame(user, word)
	if errAddGuess != nil {
		return nil, errAddGuess
	}

	return game, nil
}

func (p *UseCase) is5LetterNoun(word string) bool {
	w, _ := p.repository.FindWord(word)

	return w != nil
}
