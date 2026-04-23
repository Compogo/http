package websocket

import (
	"github.com/Compogo/compogo/component"
	"github.com/Compogo/compogo/container"
	"github.com/Compogo/compogo/flag"
)

var Component = &component.Component{
	Name: "http.server.websocket",
	Init: component.StepFunc(func(container container.Container) error {
		return container.Provides(NewConfig, NewUpgrader)
	}),
	BindFlags: component.BindFlags(func(flagSet flag.FlagSet, container container.Container) error {
		return container.Invoke(func(config *Config) {
			flagSet.IntVar(&config.ReadBufferSize, ReadBufferSizeFieldName, ReadBufferSizeDefault, "")
			flagSet.IntVar(&config.WriteBufferSize, WriteBufferSizeFieldName, WriteBufferSizeDefault, "")
			flagSet.IntVar(&config.ClientEventBufferSize, ClientEventBufferSizeFieldName, ClientEventBufferSizeDefault, "")

			flagSet.DurationVar(&config.WriteTimeout, WriteTimeoutFieldName, WriteTimeoutDefault, "")
			flagSet.DurationVar(&config.PingTimeout, PingTimeoutFieldName, PingTimeoutDefault, "")

			flagSet.StringSliceVar(&config.origins, OriginsFieldName, nil, "")
		})
	}),
	Configuration: component.StepFunc(func(container container.Container) error {
		return container.Invoke(Configuration)
	}),
}
