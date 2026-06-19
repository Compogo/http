package websocket

import (
	"time"

	"github.com/Compogo/compogo"
	"github.com/Compogo/types/set"
)

const (
	// ReadBufferSizeFieldName -размер буфера чтения
	ReadBufferSizeFieldName = "server.http.websocket.size.buffer.read"

	//WriteBufferSizeFieldName - размер буфера записи
	WriteBufferSizeFieldName = "server.http.websocket.size.buffer.write"

	// OriginsFieldName - разрешённые origins
	OriginsFieldName = "server.http.websocket.origins"

	// PingTimeoutFieldName - таймаут ping
	PingTimeoutFieldName = "server.http.websocket.timeout.ping"

	// WriteTimeoutFieldName - таймаут записи
	WriteTimeoutFieldName = "server.http.websocket.timeout.write"

	// ClientEventBufferSizeFieldName - размер буфера событий клиента
	ClientEventBufferSizeFieldName = "server.http.websocket.size.buffer.events"

	// AnyOrigin - wildcard для разрешения любого origin
	AnyOrigin = "*"
)

var (
	ReadBufferSizeDefault        = int(1024)
	WriteBufferSizeDefault       = int(1024)
	ClientEventBufferSizeDefault = 10
	WriteTimeoutDefault          = 20 * time.Second
	PingTimeoutDefault           = 60 * time.Second
)

// Config содержит конфигурацию WebSocket.
type Config struct {
	ReadBufferSize        int
	WriteBufferSize       int
	ClientEventBufferSize int

	WriteTimeout time.Duration
	PingTimeout  time.Duration
	PingInterval time.Duration

	origins []string
	Origins set.Set[string]
}

// NewConfig создаёт новую конфигурацию с значениями по умолчанию.
func NewConfig() *Config {
	return &Config{
		ReadBufferSize:  ReadBufferSizeDefault,
		WriteBufferSize: WriteBufferSizeDefault,
		PingTimeout:     PingTimeoutDefault,
	}
}

// Configuration загружает конфигурацию из Configurator.
// PingInterval вычисляется как 90% от PingTimeout для надёжности.
func Configuration(config *Config, configurator compogo.Configurator) *Config {
	if config.ReadBufferSize == 0 || config.ReadBufferSize == ReadBufferSizeDefault {
		configurator.SetDefault(ReadBufferSizeFieldName, ReadBufferSizeFieldName)
		config.ReadBufferSize = configurator.GetInt(ReadBufferSizeFieldName)
	}

	if config.WriteBufferSize == 0 || config.WriteBufferSize == WriteBufferSizeDefault {
		configurator.SetDefault(WriteBufferSizeFieldName, WriteBufferSizeDefault)
		config.WriteBufferSize = configurator.GetInt(WriteBufferSizeFieldName)
	}

	if config.ClientEventBufferSize == 0 || config.ClientEventBufferSize == ClientEventBufferSizeDefault {
		configurator.SetDefault(ClientEventBufferSizeFieldName, ClientEventBufferSizeDefault)
		config.ClientEventBufferSize = configurator.GetInt(ClientEventBufferSizeFieldName)
	}

	if config.PingTimeout == 0 || config.PingTimeout == PingTimeoutDefault {
		configurator.SetDefault(PingTimeoutFieldName, PingTimeoutDefault)
		config.PingTimeout = configurator.GetDuration(PingTimeoutFieldName)
	}

	if config.WriteTimeout == 0 || config.WriteTimeout == WriteTimeoutDefault {
		configurator.SetDefault(WriteTimeoutFieldName, WriteTimeoutDefault)
		config.WriteTimeout = configurator.GetDuration(WriteTimeoutFieldName)
	}

	config.PingInterval = (config.PingTimeout * 9) / 10

	if len(config.origins) == 0 {
		config.origins = configurator.GetStringSlice(OriginsFieldName)
	}

	config.Origins = set.NewSet[string](config.origins...)

	return config
}
