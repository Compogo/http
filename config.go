package http

import (
	"time"

	"github.com/Compogo/compogo/configurator"
)

const (
	InterfaceFieldName       = "server.http.interface"
	PortFieldName            = "server.http.port"
	ShutdownTimeoutFieldName = "server.http.timeout.shutdown"

	InterfaceDefault       = "0.0.0.0"
	PortDefault            = uint16(8080)
	ShutdownTimeoutDefault = 30 * time.Second
)

type Config struct {
	Interface       string
	Port            uint16
	ShutdownTimeout time.Duration
}

func NewConfig() *Config {
	return &Config{}
}

func Configuration(config *Config, configurator configurator.Configurator) *Config {
	if config.Interface == "" || config.Interface == InterfaceDefault {
		configurator.SetDefault(InterfaceFieldName, InterfaceDefault)
		config.Interface = configurator.GetString(InterfaceFieldName)
	}

	if config.Port == 0 || config.Port == PortDefault {
		configurator.SetDefault(PortFieldName, PortDefault)
		config.Port = configurator.GetUint16(PortFieldName)
	}

	if config.ShutdownTimeout == 0 || config.ShutdownTimeout == ShutdownTimeoutDefault {
		configurator.SetDefault(ShutdownTimeoutFieldName, ShutdownTimeoutDefault)
		config.ShutdownTimeout = configurator.GetDuration(ShutdownTimeoutFieldName)
	}

	return config
}
