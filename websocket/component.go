package websocket

import (
	"github.com/Compogo/compogo/component"
	"github.com/Compogo/compogo/container"
)

var Component = &component.Component{
	Name: "http.server.websocket",
	Init: component.StepFunc(func(container container.Container) error {
		return container.Provides(NewConfig, NewUpgrader)
	}),
	Configuration: component.StepFunc(func(container container.Container) error {
		return container.Invoke(Configuration)
	}),
}
