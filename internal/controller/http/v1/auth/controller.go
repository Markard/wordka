package auth

import (
	"errors"
	"github.com/Markard/wordka/internal/controller/http/v1/auth/login"
	"github.com/Markard/wordka/internal/controller/http/v1/auth/registration"
	"github.com/Markard/wordka/internal/usecase/auth"
	"github.com/Markard/wordka/pkg/http/response"
	"github.com/Markard/wordka/pkg/http/validator"
	"github.com/Markard/wordka/pkg/logger"
	"github.com/go-chi/render"
	"net/http"
)

type Controller struct {
	useCase   *auth.UseCase
	logger    logger.Interface
	validator validator.ProjectValidator
}

func NewController(useCase *auth.UseCase, logger logger.Interface, validator validator.ProjectValidator) *Controller {
	return &Controller{useCase: useCase, logger: logger, validator: validator}
}

func (c *Controller) Register(w http.ResponseWriter, r *http.Request) {
	converter := registration.NewConverter(c.validator)
	regRequest, converterErr := converter.ValidateAndApply(r)
	if converterErr != nil {
		response.ErrValidation(w, converterErr)
		return
	}

	user, err := c.useCase.Register(regRequest.Name, regRequest.Email, regRequest.Password)
	if err != nil {
		if errors.As(err, &auth.ErrUserAlreadyExists{}) {
			response.ErrConflict(w, err)
			return
		}
		response.ErrInternalServer(w)
		c.logger.Error(err)
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
		response.ErrValidation(w, converterErr)
		return
	}

	tokenString, err := c.useCase.Login(loginRequest.Email, loginRequest.Password)
	if err != nil {
		if errors.As(err, &auth.ErrUserNotFound{}) {
			response.ErrHttpError(w, http.StatusUnauthorized, "The credentials provided are incorrect.")
		} else {
			response.ErrInternalServer(w)
			c.logger.Error(err)
		}
		return
	}

	resp := login.NewResponse(tokenString)
	render.Status(r, http.StatusOK)
	render.JSON(w, r, resp)
}
