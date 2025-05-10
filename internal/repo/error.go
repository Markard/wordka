package repo

import (
	"fmt"
)

type ErrEmailUniqConstraint struct {
	email string
}

func (e ErrEmailUniqConstraint) Error() string {
	return fmt.Sprintf("Email: %s already exists", e.email)
}
