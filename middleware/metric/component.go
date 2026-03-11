package metric

import (
	"github.com/Compogo/compogo/component"
	"github.com/Compogo/compogo/container"
)

var (
	// RequestCountComponent is a Compogo component that provides
	// HTTP request count metrics middleware.
	RequestCountComponent = &component.Component{
		Init: component.StepFunc(func(container container.Container) error {
			return container.Provide(NewRequestCount)
		}),
	}

	// DurationComponent is a Compogo component that provides
	// HTTP request duration metrics middleware.
	DurationComponent = &component.Component{
		Init: component.StepFunc(func(container container.Container) error {
			return container.Provide(NewDuration)
		}),
	}
)
