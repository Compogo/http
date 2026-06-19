package param

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/Compogo/compogo"
	"github.com/Compogo/http_server/helper"
)

// Option — функция настройки параметра.
type Option func(param *Param) *Param

// Caster — функция преобразования строкового значения в целевой тип.
type Caster func(any) (any, error)

// Getter — функция извлечения строкового значения из HTTP-запроса.
type Getter func(request *http.Request) string

// Validator — функция проверки корректности значения.
type Validator func(any) error

// DefaultValueFunc — функция для получения значения по умолчанию.
type DefaultValueFunc func() any

// Param — middleware для извлечения, приведения и валидации параметров из HTTP-запроса.
//
// Поддерживает:
//   - Извлечение из URL-query, заголовков, cookies
//   - Приведение к различным типам (string, int, float, bool, time, duration, ip)
//   - Валидацию значений
//   - Значения по умолчанию
//   - Сохранение в контексте запроса
//
// Пример:
//
//	userIdParam := NewParamInt("user_id", logger,
//	    WithUriGetter(),
//	    AddValidator(GteValidator(1)),
//	)
//	router.Get("/users/:id", userIdParam.Middleware(http.HandlerFunc(handler)))
type Param struct {
	name string

	caster       Caster
	defaultValue DefaultValueFunc
	getters      []Getter
	validators   []Validator

	logger compogo.Logger
}

// NewParam создаёт новый параметр с указанными опциями.
// Принимает имя, логгер, функцию приведения и опции.
func NewParam(name string, logger compogo.Logger, caster Caster, options ...Option) *Param {
	param := &Param{
		name:   name,
		logger: logger.GetLogger("http").GetLogger("server").GetLogger("param").GetLogger(name),
		caster: caster,
	}

	for _, option := range options {
		option(param)
	}

	return param
}

// Middleware реализует интерфейс http_server.Middleware.
// Извлекает параметр, приводит к целевому типу, валидирует и сохраняет в контекст.
// В случае ошибки возвращает 400 Bad Request.
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

// Name возвращает имя параметра.
func (param *Param) Name() string {
	return param.name
}
