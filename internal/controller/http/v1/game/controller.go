package game

import (
	"errors"
	"github.com/Markard/wordka/internal/controller/http/v1/game/currentgame"
	"github.com/Markard/wordka/internal/entity"
	"github.com/Markard/wordka/internal/infra/middleware/jwt"
	"github.com/Markard/wordka/internal/usecase/game"
	"github.com/Markard/wordka/pkg/httpserver/response"
	"github.com/Markard/wordka/pkg/logger"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"net/http"
)

type Controller struct {
	useCase   *game.UseCase
	logger    logger.Interface
	validator *validator.Validate
}

func NewController(useCase *game.UseCase, logger logger.Interface, validator *validator.Validate) *Controller {
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
