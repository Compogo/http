package websocket

import (
	"github.com/Compogo/compogo"
	"github.com/Compogo/compogo/flag"
)

// Component — компонент WebSocket для Compogo.
// Регистрирует конфигурацию и Upgrader в DI-контейнере.
var Component = compogo.Component{
	Name: "http.server.websocket",
	Init: compogo.StepFunc(func(container compogo.Container) error {
		return container.Provides(NewConfig, NewUpgrader)
	}),
	BindFlags: compogo.BindFlags(func(flagSet flag.FlagSet, container compogo.Container) error {
		return container.Invoke(func(config *Config) {
			flagSet.IntVar(&config.ReadBufferSize, ReadBufferSizeFieldName, ReadBufferSizeDefault, "")
			flagSet.IntVar(&config.WriteBufferSize, WriteBufferSizeFieldName, WriteBufferSizeDefault, "")
			flagSet.IntVar(&config.ClientEventBufferSize, ClientEventBufferSizeFieldName, ClientEventBufferSizeDefault, "")

			flagSet.DurationVar(&config.WriteTimeout, WriteTimeoutFieldName, WriteTimeoutDefault, "")
			flagSet.DurationVar(&config.PingTimeout, PingTimeoutFieldName, PingTimeoutDefault, "")

			flagSet.StringSliceVar(&config.origins, OriginsFieldName, nil, "")
		})
	}),
	Configuration: compogo.StepFunc(func(container compogo.Container) error {
		return container.Invoke(Configuration)
	}),
}
