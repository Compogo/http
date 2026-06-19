package http_server

import (
	"time"

	"github.com/Compogo/compogo"
)

const (
	// InterfaceFieldName — имя поля для сетевого интерфейса.
	InterfaceFieldName = "server.http.interface"

	// PortFieldName — имя поля для порта.
	PortFieldName = "server.http.port"

	// ShutdownTimeoutFieldName — имя поля для таймаута завершения.
	ShutdownTimeoutFieldName = "server.http.timeout.shutdown"
)

var (
	// InterfaceDefault — сетевой интерфейс по умолчанию (все интерфейсы).
	InterfaceDefault = "0.0.0.0"

	// PortDefault — порт по умолчанию.
	PortDefault = uint16(8080)

	// ShutdownTimeoutDefault — таймаут завершения по умолчанию (30 секунд).
	ShutdownTimeoutDefault = 30 * time.Second
)

// Config содержит конфигурацию HTTP-сервера.
type Config struct {
	Interface       string
	Port            uint16
	ShutdownTimeout time.Duration
}

// NewConfig создаёт новую конфигурацию.
func NewConfig() *Config {
	return &Config{}
}

// Configuration загружает конфигурацию из Configurator.
// Если значения не заданы, устанавливаются значения по умолчанию.
func Configuration(config *Config, configurator compogo.Configurator) *Config {
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
