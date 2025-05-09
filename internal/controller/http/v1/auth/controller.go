package auth

import (
	"errors"
	"github.com/Markard/wordka/internal/controller/http/v1/auth/login"
	"github.com/Markard/wordka/internal/controller/http/v1/auth/registration"
	"github.com/Markard/wordka/internal/usecase"
	"github.com/Markard/wordka/pkg/httpserver/response"
	"github.com/Markard/wordka/pkg/logger"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"net/http"
)

type Controller struct {
	auth      *usecase.Auth
	logger    logger.Interface
	validator *validator.Validate
}

func NewController(auth *usecase.Auth, logger logger.Interface, validator *validator.Validate) *Controller {
	return &Controller{auth: auth, logger: logger, validator: validator}
}

func (c *Controller) Register(w http.ResponseWriter, r *http.Request) {
	converter := registration.NewConverter(c.validator)
	registrationReq, errRenderer := converter.ValidateAndApply(r)
	if errRenderer != nil {
		render.JSON(w, r, errRenderer)
		return
	}

	user, err := c.auth.Register(registrationReq.Name, registrationReq.Email, registrationReq.Password)
	if err != nil {
		render.JSON(w, r, response.ErrInvalidRequest(err))
		return
	}

	resp := registration.NewResponse(user)
	render.Status(r, http.StatusCreated)
	render.JSON(w, r, resp)
}

func (c *Controller) Login(w http.ResponseWriter, r *http.Request) {
	converter := login.NewConverter(c.validator)
	loginRequest, converterErr := converter.ValidateAndApply(r)
	if converterErr != nil {
		render.JSON(w, r, converterErr)
		return
	}

	tokenString, err := c.auth.Login(loginRequest.Email, loginRequest.Password)
	if err != nil {
		if errors.As(err, &usecase.ErrUserNotFound{}) {
			render.JSON(w, r, response.ErrIncorrectCredentials(err))
		} else {
			render.JSON(w, r, response.ErrInternalServer(err))
			c.logger.Error(err)
		}
		return
	}

	resp := login.NewResponse(tokenString)
	render.Status(r, http.StatusOK)
	render.JSON(w, r, resp)
}
