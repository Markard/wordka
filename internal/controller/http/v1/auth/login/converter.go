package login

import (
	"encoding/json"
	"github.com/Markard/wordka/pkg/http/response"
	"github.com/Markard/wordka/pkg/http/validator"
	"github.com/go-chi/render"
	"net/http"
)

type Converter struct {
	validator validator.ProjectValidator
}

func NewConverter(validator validator.ProjectValidator) *Converter {
	return &Converter{validator: validator}
}

func (c *Converter) ValidateAndApply(r *http.Request) (*Request, render.Renderer) {
	loginReq := &Request{}

	err := json.NewDecoder(r.Body).Decode(loginReq)
	if err != nil {
		return nil, response.ErrInvalidJson(err)
	}

	if errVal := c.validator.Struct(loginReq); errVal != nil {
		return nil, response.ErrValidation(errVal)
	}

	return loginReq, nil
}
