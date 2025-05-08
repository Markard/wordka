package usecase

import (
	"github.com/Markard/wordka/internal/entity"
)

type UserCreator interface {
	Create(user *entity.User) error
}

type Auth struct {
	creator UserCreator
}

func NewAuth(creator UserCreator) *Auth {
	return &Auth{
		creator: creator,
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
