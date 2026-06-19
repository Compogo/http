package http_server

import (
	"github.com/Compogo/compogo"
	"github.com/Compogo/compogo/flag"
	"github.com/Compogo/runner"
)

// Component — компонент HTTP-сервера для Compogo.
// Регистрирует сервер в DI-контейнере и запускает его через Runner.
var Component = compogo.Component{
	Name: "http.server",
	Dependencies: compogo.Components{
		&runner.Component,
	},
	Init: compogo.StepFunc(func(container compogo.Container) error {
		return container.Provides(
			NewConfig,
			NewServer,
		)
	}),
	BindFlags: compogo.BindFlags(func(flagSet flag.FlagSet, container compogo.Container) error {
		return container.Invoke(func(config *Config) {
			flagSet.StringVar(&config.Interface, InterfaceFieldName, InterfaceDefault, "interface for listening to incoming requests")
			flagSet.Uint16Var(&config.Port, PortFieldName, PortDefault, "port for listening to incoming requests")
			flagSet.DurationVar(
				&config.ShutdownTimeout,
				ShutdownTimeoutFieldName,
				ShutdownTimeoutDefault,
				"timeout during which the HTTP server must be turned off after receiving a shutdown signal",
			)
		})
	}),
	Configuration: compogo.StepFunc(func(container compogo.Container) error {
		return container.Invoke(Configuration)
	}),
	PreWait: compogo.StepFunc(func(container compogo.Container) error {
		return container.Invoke(func(r runner.Runner, server Server) error {
			return r.RunProcess(server)
		})
	}),
}
