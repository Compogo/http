package token

import (
	"github.com/Compogo/compogo"
)

// Component — компонент Token Auth для Compogo.
// Регистрирует конфигурацию и middleware в DI-контейнере.
var Component = compogo.Component{
	Init: compogo.StepFunc(func(container compogo.Container) error {
		return container.Provides(
			NewConfig,
			NewAuth,
		)
	}),
	Configuration: compogo.StepFunc(func(container compogo.Container) error {
		return container.Invoke(Configuration)
	}),
}
