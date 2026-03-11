package logger

import (
	"github.com/Compogo/compogo/component"
	"github.com/Compogo/compogo/container"
)

var (
	// RequestComponent is a Compogo component that provides request logging middleware.
	RequestComponent = &component.Component{
		Init: component.StepFunc(func(container container.Container) error {
			return container.Provide(NewRequest)
		}),
	}

	// ResponseComponent is a Compogo component that provides response logging middleware.
	ResponseComponent = &component.Component{
		Init: component.StepFunc(func(container container.Container) error {
			return container.Provide(NewResponse)
		}),
	}
)
