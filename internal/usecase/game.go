package usecase

import (
	"github.com/Markard/wordka/internal/entity"
)

type IGameRepository interface {
	FindCurrentGame() (*entity.Game, error)
}

type GameUseCase struct {
	repository IGameRepository
}

func NewGameUseCase(repository IGameRepository) *GameUseCase {
	return &GameUseCase{repository: repository}
}

func (p *GameUseCase) FindCurrentGame(userId int64) (*entity.Game, error) {
	return p.repository.FindCurrentGame()
}
