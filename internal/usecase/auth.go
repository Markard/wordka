package usecase

import (
	"errors"
	"github.com/Markard/wordka/internal/entity"
	"github.com/Markard/wordka/internal/repo"
	"github.com/Markard/wordka/pkg/jwtauth"
)

type IAuthRepository interface {
	Create(user *entity.User) error
	FindBy(email string) (*entity.User, error)
}

type Auth struct {
	repository   IAuthRepository
	tokenService *jwtauth.TokenService
}

func NewAuth(repository IAuthRepository, tokenService *jwtauth.TokenService) *Auth {
	return &Auth{repository: repository, tokenService: tokenService}
}

func (auth *Auth) Register(name string, email string, rawPassword string) (*entity.User, error) {
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

func (auth *Auth) Login(email string, password string) (string, error) {
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

	tokenString, err := auth.tokenService.CreateTokenStringWithES256(user.Id)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
