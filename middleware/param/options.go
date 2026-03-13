package param

import (
	"fmt"
	"net/http"

	"github.com/spf13/cast"
)

// WithUriGetter adds a getter that extracts the parameter from URL query string.
// Uses the parameter's name as the query key.
func WithUriGetter() Option {
	return func(param *Param) *Param {
		param.getters = append(
			param.getters, func(request *http.Request) string {
				return request.URL.Query().Get(param.name)
			},
		)

		return param
	}
}

// WithHeaderGetter adds a getter that extracts the parameter from HTTP headers.
// Uses the parameter's name as the header key.
func WithHeaderGetter() Option {
	return func(param *Param) *Param {
		param.getters = append(
			param.getters, func(request *http.Request) string {
				return request.Header.Get(param.name)
			},
		)

		return param
	}
}

// WithCookieGetter adds a getter that extracts the parameter from cookies.
// Uses the parameter's name as the cookie name.
func WithCookieGetter() Option {
	return func(param *Param) *Param {
		param.getters = append(
			param.getters, func(request *http.Request) string {
				for _, cookie := range request.Cookies() {
					if cookie.Name == param.name {
						return cookie.Value
					}
				}

				return ""
			},
		)

		return param
	}
}

// WithUriGetterByName adds a getter that extracts from URL query with a custom key.
func WithUriGetterByName(name string) Option {
	return AddGetter(func(request *http.Request) string {
		return request.URL.Query().Get(name)
	})
}

// WithHeaderGetterByName adds a getter that extracts from headers with a custom key.
func WithHeaderGetterByName(name string) Option {
	return AddGetter(func(request *http.Request) string {
		return request.Header.Get(name)
	})
}

// WithCookieGetterByName adds a getter that extracts from cookies with a custom key.
func WithCookieGetterByName(name string) Option {
	return AddGetter(func(request *http.Request) string {
		for _, cookie := range request.Cookies() {
			if cookie.Name == name {
				return cookie.Value
			}
		}

		return ""
	})
}

// AddGetter adds a custom getter function to the parameter.
func AddGetter(getter Getter) Option {
	return func(param *Param) *Param {
		param.getters = append(param.getters, getter)

		return param
	}
}

// AddValidator adds a custom validator function to the parameter.
func AddValidator(validator Validator) Option {
	return func(param *Param) *Param {
		param.validators = append(param.validators, validator)

		return param
	}
}

// WithDefault sets a static default value for the parameter.
func WithDefault(value any) Option {
	return WithDefaultFunc(func() any {
		return value
	})
}

// WithDefaultFunc sets a dynamic default value provider for the parameter.
func WithDefaultFunc(defaultFunc DefaultValueFunc) Option {
	return func(param *Param) *Param {
		param.defaultValue = defaultFunc

		return param
	}
}

// GteValidator returns a validator that checks if a numeric value is >= the given threshold.
func GteValidator(val float64) Validator {
	return func(value any) error {
		rval, err := cast.ToFloat64E(value)
		if err != nil {
			return err
		}

		if rval < val {
			return fmt.Errorf("validator.gte: %.4g < %.4g", rval, val)
		}

		return nil
	}
}
