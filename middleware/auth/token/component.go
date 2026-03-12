package token

import (
	"github.com/Compogo/compogo/component"
	"github.com/Compogo/compogo/container"
)

// Component is a ready-to-use Compogo component that provides
// HTTP Token Authentication middleware.
//
// Usage:
//
//	compogo.WithComponents(
//	    http.Component,
//	    token.Component,
//	)
//
// Then in your router setup:
//
//	router.Use(tokenAuth.Middleware)
var Component = &component.Component{
	Init: component.StepFunc(func(container container.Container) error {
		return container.Provides(
			NewConfig,
			NewAuth,
		)
	}),
	Configuration: component.StepFunc(func(container container.Container) error {
		return container.Invoke(Configuration)
	}),
}
