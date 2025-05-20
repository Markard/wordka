package game

import (
	"errors"
	"fmt"
	"github.com/Markard/wordka/internal/controller/http/v1/game/currentgame"
	"github.com/Markard/wordka/internal/controller/http/v1/game/guess"
	"github.com/Markard/wordka/internal/entity"
	"github.com/Markard/wordka/internal/infra/middleware/jwt"
	"github.com/Markard/wordka/internal/usecase/game"
	"github.com/Markard/wordka/pkg/http/response"
	"github.com/Markard/wordka/pkg/http/validator"
	"github.com/Markard/wordka/pkg/logger"
	"github.com/go-chi/render"
	"net/http"
)

type Controller struct {
	useCase   *game.UseCase
	logger    logger.Interface
	validator validator.ProjectValidator
}

func NewController(useCase *game.UseCase, logger logger.Interface, validator validator.ProjectValidator) *Controller {
	return &Controller{useCase: useCase, logger: logger, validator: validator}
}

func (c *Controller) GetCurrentGame(w http.ResponseWriter, r *http.Request) {
	currentUser, _ := r.Context().Value(jwt.CurrentUserCtxKey).(*entity.User)
	currentGame, err := c.useCase.FindCurrentGame(currentUser)
	if err != nil {
		if errors.Is(err, game.ErrCurrentGameNotFound) {
			response.ErrNotFound(w, err)
			return
		} else {
			c.logger.Error(err)
			return
		}
	}

	resp := currentgame.NewResponse(currentGame)
	render.Status(r, http.StatusOK)
	render.JSON(w, r, resp)
}

func (c *Controller) CreateGame(w http.ResponseWriter, r *http.Request) {
	currentUser, _ := r.Context().Value(jwt.CurrentUserCtxKey).(*entity.User)
	currentGame, err := c.useCase.CreateGame(currentUser)

	if err != nil {
		if errors.Is(err, game.ErrCurrentGameAlreadyExists) {
			response.ErrConflict(w, err)
			return
		} else {
			c.logger.Error(err)
			return
		}
	}

	resp := currentgame.NewResponse(currentGame)
	render.Status(r, http.StatusOK)
	render.JSON(w, r, resp)
}

func (c *Controller) Guess(w http.ResponseWriter, r *http.Request) {
	converter := guess.NewConverter(c.validator)
	guessReq, valErr := converter.ValidateAndApply(r)
	if valErr != nil {
		valErr.ErrValidation(w)
		return
	}

	currentUser, _ := r.Context().Value(jwt.CurrentUserCtxKey).(*entity.User)
	fmt.Printf(currentUser.Email)
	fmt.Printf(guessReq.Word)

	currentGame, err := c.useCase.Guess(currentUser, guessReq.Word)
	if err != nil {
		if errors.Is(err, game.ErrCurrentGameNotFound) {
			response.ErrNotFound(w, err)
			return
		} else if errors.Is(err, game.ErrIncorrectWord) {
			response.
				NewValidationError().
				AddFieldError("word", "The word must be a Russian noun consisting of exactly 5 letters").
				ErrValidation(w)
			return
		} else {
			c.logger.Error(err)
			return
		}
	}

	resp := currentgame.NewResponse(currentGame)
	render.Status(r, http.StatusCreated)
	render.JSON(w, r, resp)
}
