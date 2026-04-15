package http

import (
	"github.com/Compogo/compogo/component"
	"github.com/Compogo/compogo/container"
	"github.com/Compogo/compogo/flag"
	"github.com/Compogo/runner"
)

// Component is a ready-to-use Compogo component that provides an HTTP server.
// It automatically:
//   - Registers Config and Server in the DI container
//   - Adds command-line flags for server configuration
//   - Configures the server during PreRun phase
//   - Starts the server as a runner task during PostRun phase
//   - Performs graceful shutdown during Stop phase
//
// Usage:
//
//	compogo.WithComponents(
//	    runner.Component,
//	    http.Component,
//	    myRouterComponent,  // must implement http.Router
//	)
var Component = &component.Component{
	Dependencies: component.Components{
		runner.Component,
	},
	Init: component.StepFunc(func(container container.Container) error {
		return container.Provides(
			NewConfig,
			NewServer,
		)
	}),
	BindFlags: component.BindFlags(func(flagSet flag.FlagSet, container container.Container) error {
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
	Configuration: component.StepFunc(func(container container.Container) error {
		return container.Invoke(Configuration)
	}),
	PreWait: component.StepFunc(func(container container.Container) error {
		return container.Invoke(func(r runner.Runner, server Server) error {
			return r.RunTask(runner.NewTask("server.http", server))
		})
	}),
}
