package auth

import (
	"fmt"
)

type ErrUserNotFound struct {
	email string
}

func (e ErrUserNotFound) Error() string {
	return fmt.Sprintf("User with emai: %s not found", e.email)
}

type ErrUserAlreadyExists struct {
	email string
}

func (e ErrUserAlreadyExists) Error() string {
	return fmt.Sprintf("User with emai: %s already exists", e.email)
}
