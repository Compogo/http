package basic

import (
	"github.com/Compogo/compogo/component"
	"github.com/Compogo/compogo/container"
)

// Component is a ready-to-use Compogo component that provides
// HTTP Basic Authentication middleware.
//
// Usage:
//
//	compogo.WithComponents(
//	    http.Component,
//	    basic.Component,
//	)
//
// Then in your router setup:
//
//	router.Use(basicAuth.Middleware)
var Component = &component.Component{
	Init: component.StepFunc(func(container container.Container) error {
		return container.Provides(
			NewConfig,
			NewAuth,
		)
	}),
	PreRun: component.StepFunc(func(container container.Container) error {
		return container.Invoke(Configuration)
	}),
}
