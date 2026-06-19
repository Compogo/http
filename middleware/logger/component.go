package logger

import (
	"github.com/Compogo/compogo"
)

var (
	// RequestComponent — компонент middleware логирования запросов для Compogo.
	// Регистрирует Request в DI-контейнере.
	RequestComponent = compogo.Component{
		Init: compogo.StepFunc(func(container compogo.Container) error {
			return container.Provide(NewRequest)
		}),
	}

	// ResponseComponent — компонент middleware логирования ответов для Compogo.
	// Регистрирует Response в DI-контейнере.
	ResponseComponent = compogo.Component{
		Init: compogo.StepFunc(func(container compogo.Container) error {
			return container.Provide(NewResponse)
		}),
	}
)
