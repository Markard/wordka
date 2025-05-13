package game

import (
	"github.com/Markard/wordka/internal/entity"
)

type IGameRepository interface {
	FindCurrentGame(currentUser *entity.User) (*entity.Game, error)
}

type GameUseCase struct {
	repository IGameRepository
}

func NewGameUseCase(repository IGameRepository) *GameUseCase {
	return &GameUseCase{repository: repository}
}

func (p *GameUseCase) FindCurrentGame(user *entity.User) (*entity.Game, error) {
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
