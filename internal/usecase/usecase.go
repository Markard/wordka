package usecase

import (
	"github.com/Markard/wordka/internal/usecase/auth"
)

type UseCases struct {
	AuthUseCase *auth.UseCase
	GameUseCase *GameUseCase
}
