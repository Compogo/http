package param

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/Compogo/compogo/logger"
	"github.com/Compogo/http/helper"
)

// Option configures a Param during creation.
type Option func(param *Param) *Param

// Caster converts a raw string value to the desired type.
// Returns the converted value or an error.
type Caster func(any) (any, error)

// Getter extracts a string value from an HTTP request.
type Getter func(request *http.Request) string

// Validator validates a typed value.
// Returns an error if validation fails.
type Validator func(any) error

// DefaultValueFunc provides a default value when none is found in the request.
type DefaultValueFunc func() any

// Param defines a request parameter with extraction, casting, and validation rules.
// It implements http.Middleware, injecting the parsed value into the request context.
type Param struct {
	name string

	caster       Caster
	defaultValue DefaultValueFunc
	getters      []Getter
	validators   []Validator

	logger logger.Logger
}

// NewParam creates a new Param with the given name, logger, caster, and options.
// The parameter value will be stored in the request context under the parameter name.
func NewParam(name string, logger logger.Logger, caster Caster, options ...Option) *Param {
	param := &Param{
		name:   name,
		logger: logger.GetLogger("http.server.param." + name),
		caster: caster,
	}

	for _, option := range options {
		option(param)
	}

	return param
}

// Middleware implements http.Middleware.
// It extracts, casts, validates the parameter and injects it into the request context.
// If any step fails, it returns an appropriate HTTP error response.
func (param *Param) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var requestValue string

		for _, getter := range param.getters {
			requestValue = getter(r)
			if strings.TrimSpace(requestValue) != "" {
				break
			}
		}

		if requestValue == "" && param.defaultValue == nil {
			err := fmt.Sprintf("name - '%s' empty", param.name)
			helper.WriteError(w, r, err, http.StatusBadRequest)
			param.logger.Error(err)
			return
		}

		value := any(requestValue)
		if requestValue == "" {
			value = param.defaultValue()
		}

		var err error

		value, err = param.caster(value)
		if err != nil {
			helper.WriteError(w, r, err.Error(), http.StatusBadRequest)
			param.logger.Error(err)
			return
		}

		for _, validator := range param.validators {
			if err = validator(value); err != nil {
				helper.WriteError(w, r, err.Error(), http.StatusBadRequest)
				param.logger.Error(err)
				return
			}
		}

		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), param.name, value)))
	})
}

func (param *Param) Name() string {
	return param.name
}
