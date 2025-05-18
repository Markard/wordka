package game

import (
	"errors"
	"fmt"
	"github.com/Markard/wordka/internal/controller/http/v1/game/currentgame"
	"github.com/Markard/wordka/internal/controller/http/v1/game/guess"
	"github.com/Markard/wordka/internal/entity"
	"github.com/Markard/wordka/internal/infra/middleware/jwt"
	"github.com/Markard/wordka/internal/usecase/game"
	"github.com/Markard/wordka/pkg/httpserver"
	"github.com/Markard/wordka/pkg/httpserver/response"
	"github.com/Markard/wordka/pkg/logger"
	"github.com/go-chi/render"
	"net/http"
)

type Controller struct {
	useCase   *game.UseCase
	logger    logger.Interface
	validator httpserver.ProjectValidator
}

func NewController(useCase *game.UseCase, logger logger.Interface, validator httpserver.ProjectValidator) *Controller {
	return &Controller{useCase: useCase, logger: logger, validator: validator}
}

func (c *Controller) GetCurrentGame(w http.ResponseWriter, r *http.Request) {
	currentUser, _ := r.Context().Value(jwt.CurrentUserCtxKey).(*entity.User)
	currentGame, err := c.useCase.FindCurrentGame(currentUser)
	if err != nil {
		if errors.As(err, &game.ErrCurrentGameNotFound{}) {
			_ = render.Render(w, r, response.ErrNotFound(err))
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
		if errors.As(err, &game.ErrCurrentGameAlreadyExists{}) {
			_ = render.Render(w, r, response.ErrConflict(err))
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
	guessReq, converterErr := converter.ValidateAndApply(r)
	if converterErr != nil {
		_ = render.Render(w, r, converterErr)
		return
	}

	currentUser, _ := r.Context().Value(jwt.CurrentUserCtxKey).(*entity.User)
	fmt.Printf(currentUser.Email)
	fmt.Printf(guessReq.Word)

	currentGame, err := c.useCase.Guess(currentUser, guessReq.Word)
	if err != nil {
		if errors.As(err, &game.ErrCurrentGameNotFound{}) {
			_ = render.Render(w, r, response.ErrNotFound(err))
			return
		} else if errors.As(err, &game.ErrIncorrectWord{}) {
			errIncorrectWord := response.NewCustomValidationErrs(
				"word",
				"The word must be a Russian noun consisting of exactly 5 letters",
			)
			_ = render.Render(w, r, response.ErrValidation(errIncorrectWord))
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
