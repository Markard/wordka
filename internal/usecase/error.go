package usecase

import (
	"fmt"
)

type ErrUserNotFound struct {
	email string
}

func (e ErrUserNotFound) Error() string {
	return fmt.Sprintf("User with emai: %s not found", e.email)
}
