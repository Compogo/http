package param

import (
	"fmt"
	"net/http"

	"github.com/spf13/cast"
)

// WithUriGetter добавляет getter для извлечения параметра из URL-query.
// Использует имя параметра, заданное при создании Param.
func WithUriGetter() Option {
	return func(param *Param) *Param {
		param.getters = append(
			param.getters, func(request *http.Request) string {
				return request.URL.Query().Get(param.Name())
			},
		)

		return param
	}
}

// WithHeaderGetter добавляет getter для извлечения параметра из заголовка.
// Использует имя параметра, заданное при создании Param.
func WithHeaderGetter() Option {
	return func(param *Param) *Param {
		param.getters = append(
			param.getters, func(request *http.Request) string {
				return request.Header.Get(param.Name())
			},
		)

		return param
	}
}

// WithCookieGetter добавляет getter для извлечения параметра из cookie.
// Использует имя параметра, заданное при создании Param.
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

// WithUriGetterByName добавляет getter для извлечения параметра из URL-query
// с указанием имени параметра (полезно, когда имя в запросе отличается от имени в контексте).
func WithUriGetterByName(name string) Option {
	return AddGetter(func(request *http.Request) string {
		return request.URL.Query().Get(name)
	})
}

// WithHeaderGetterByName добавляет getter для извлечения параметра из заголовка
// с указанием имени параметра.
func WithHeaderGetterByName(name string) Option {
	return AddGetter(func(request *http.Request) string {
		return request.Header.Get(name)
	})
}

// WithCookieGetterByName добавляет getter для извлечения параметра из cookie
// с указанием имени параметра.
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

// AddGetter добавляет произвольный getter для извлечения параметра.
func AddGetter(getter Getter) Option {
	return func(param *Param) *Param {
		param.getters = append(param.getters, getter)

		return param
	}
}

// AddValidator добавляет валидатор для проверки значения параметра.
func AddValidator(validator Validator) Option {
	return func(param *Param) *Param {
		param.validators = append(param.validators, validator)

		return param
	}
}

// WithDefault устанавливает значение по умолчанию для параметра.
func WithDefault(value any) Option {
	return WithDefaultFunc(func() any {
		return value
	})
}

// WithDefaultFunc устанавливает функцию для получения значения по умолчанию.
// Функция вызывается каждый раз, когда параметр отсутствует в запросе.
func WithDefaultFunc(defaultFunc DefaultValueFunc) Option {
	return func(param *Param) *Param {
		param.defaultValue = defaultFunc

		return param
	}
}

// GteValidator создаёт валидатор, проверяющий что значение >= указанного.
// Работает с числовыми типами (int, float, и т.д.).
//
// Пример:
//
//	param := NewParamInt("age", logger, AddValidator(GteValidator(18)))
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
