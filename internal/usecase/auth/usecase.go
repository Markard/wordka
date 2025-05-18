package auth

import (
	"errors"
	"github.com/Markard/wordka/internal/entity"
	serviceJwt "github.com/Markard/wordka/internal/infra/service/jwt"
	"github.com/Markard/wordka/internal/repo"
)

type IAuthRepository interface {
	Create(user *entity.User) error
	FindBy(email string) (*entity.User, error)
}

type UseCase struct {
	repository IAuthRepository
	jwtService *serviceJwt.Service
}

func NewAuth(repository IAuthRepository, tokenService *serviceJwt.Service) *UseCase {
	return &UseCase{repository: repository, jwtService: tokenService}
}

func (auth *UseCase) Register(name string, email string, rawPassword string) (*entity.User, error) {
	user, err := entity.NewUser(name, email, rawPassword)
	if err != nil {
		return nil, err
	}
	err = auth.repository.Create(user)
	if err != nil {
		if errors.As(err, &repo.ErrEmailUniqConstraint{}) {
			return nil, ErrUserAlreadyExists{email}
		}
		return nil, err
	}

	return user, nil
}

func (auth *UseCase) Login(email string, password string) (string, error) {
	user, err := auth.repository.FindBy(email)
	if err != nil {
		if user == nil {
			return "", ErrUserNotFound{email}
		} else {
			return "", err
		}
	}

	if !user.IsPasswordMatch(password) {
		return "", ErrUserNotFound{email}
	}

	tokenString, err := auth.jwtService.CreateTokenStringWithES256(user.Id)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
