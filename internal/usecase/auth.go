package usecase

import (
	"github.com/Markard/wordka/internal/entity"
	"github.com/Markard/wordka/pkg/jwtauth"
)

type UserCreator interface {
	Create(user *entity.User) error
}

type UserProvider interface {
	FindBy(email string) (*entity.User, error)
}

type Auth struct {
	creator      UserCreator
	provider     UserProvider
	tokenService *jwtauth.TokenService
}

func NewAuth(
	creator UserCreator,
	provider UserProvider,
	tokenService *jwtauth.TokenService,
) *Auth {
	return &Auth{
		creator:      creator,
		provider:     provider,
		tokenService: tokenService,
	}
}

func (auth *Auth) Register(name string, email string, rawPassword string) (*entity.User, error) {
	user, err := entity.NewUser(name, email, rawPassword)
	if err != nil {
		return nil, err
	}
	err = auth.creator.Create(user)

	return user, err
}

func (auth *Auth) Login(email string, password string) (string, error) {
	user, err := auth.provider.FindBy(email)
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
