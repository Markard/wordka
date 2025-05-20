package registration

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
	registrationReq := &Request{}

	err := json.NewDecoder(r.Body).Decode(registrationReq)
	if err != nil {
		return nil, response.ErrInvalidJson(err)
	}

	if errVal := c.validator.Struct(registrationReq); errVal != nil {
		return nil, response.ErrValidation(errVal)
	}

	return registrationReq, nil
}
