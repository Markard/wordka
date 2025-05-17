package game

import (
	"fmt"
)

type ErrCurrentGameNotFound struct{}

func (e ErrCurrentGameNotFound) Error() string {
	return fmt.Sprintf("The current user is not playing any game now")
}

type ErrIncorrectWord struct{}

func (e ErrIncorrectWord) Error() string {
	return fmt.Sprintf("The word you entered is not a 5-letter noun")
}

type ErrCurrentGameAlreadyExists struct{}

func (e ErrCurrentGameAlreadyExists) Error() string {
	return fmt.Sprintf("The current user is already playing a game")
}
