package usecase

import (
	"github.com/Markard/wordka/internal/usecase/auth"
	"github.com/Markard/wordka/internal/usecase/game"
)

type UseCases struct {
	AuthUseCase *auth.UseCase
	GameUseCase *game.GameUseCase
}
