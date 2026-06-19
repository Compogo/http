package param

import (
	"fmt"
	"net"
	"net/http"
	"strings"

	"github.com/Compogo/compogo"
	"github.com/araddon/dateparse"
	"github.com/spf13/cast"
)

const (
	// HeaderRealIp — заголовок с реальным IP клиента (nginx).
	HeaderRealIp = "X-REAL-IP"

	// HeaderForwardedFor — заголовок с цепочкой IP (X-Forwarded-For).
	HeaderForwardedFor = "X-FORWARDED-FOR"

	// ipSeparator — разделитель IP-адресов в заголовке X-Forwarded-For.
	ipSeparator = ","
)

// NewParamString создаёт параметр строкового типа.
func NewParamString(name string, logger compogo.Logger, options ...Option) *Param {
	return NewParam(
		name,
		logger,
		func(value any) (any, error) { return cast.ToStringE(value) },
		options...,
	)
}

// NewParamInt создаёт параметр типа int.
func NewParamInt(name string, logger compogo.Logger, options ...Option) *Param {
	return NewParam(
		name,
		logger,
		func(value any) (any, error) { return cast.ToIntE(value) },
		options...,
	)
}

// NewParamInt8 создаёт параметр типа int8.
func NewParamInt8(name string, logger compogo.Logger, options ...Option) *Param {
	return NewParam(
		name,
		logger,
		func(value any) (any, error) { return cast.ToInt8E(value) },
		options...,
	)
}

// NewParamInt16 создаёт параметр типа int16.
func NewParamInt16(name string, logger compogo.Logger, options ...Option) *Param {
	return NewParam(
		name,
		logger,
		func(value any) (any, error) { return cast.ToInt16E(value) },
		options...,
	)
}

// NewParamInt32 создаёт параметр типа int32.
func NewParamInt32(name string, logger compogo.Logger, options ...Option) *Param {
	return NewParam(
		name,
		logger,
		func(value any) (any, error) { return cast.ToInt32E(value) },
		options...,
	)
}

// NewParamInt64 создаёт параметр типа int64.
func NewParamInt64(name string, logger compogo.Logger, options ...Option) *Param {
	return NewParam(
		name,
		logger,
		func(value any) (any, error) { return cast.ToInt64E(value) },
		options...,
	)
}

// NewParamFloat32 создаёт параметр типа float32.
func NewParamFloat32(name string, logger compogo.Logger, options ...Option) *Param {
	return NewParam(
		name,
		logger,
		func(value any) (any, error) { return cast.ToFloat32E(value) },
		options...,
	)
}

// NewParamFloat64 создаёт параметр типа float64.
func NewParamFloat64(name string, logger compogo.Logger, options ...Option) *Param {
	return NewParam(
		name,
		logger,
		func(value any) (any, error) { return cast.ToFloat64E(value) },
		options...,
	)
}

// NewParamUint создаёт параметр типа uint.
func NewParamUint(name string, logger compogo.Logger, options ...Option) *Param {
	return NewParam(
		name,
		logger,
		func(value any) (any, error) { return cast.ToUintE(value) },
		options...,
	)
}

// NewParamUint8 создаёт параметр типа uint8.
func NewParamUint8(name string, logger compogo.Logger, options ...Option) *Param {
	return NewParam(
		name,
		logger,
		func(value any) (any, error) { return cast.ToUint8E(value) },
		options...,
	)
}

// NewParamUint16 создаёт параметр типа uint16.
func NewParamUint16(name string, logger compogo.Logger, options ...Option) *Param {
	return NewParam(
		name,
		logger,
		func(value any) (any, error) { return cast.ToUint16E(value) },
		options...,
	)
}

// NewParamUint32 создаёт параметр типа uint32.
func NewParamUint32(name string, logger compogo.Logger, options ...Option) *Param {
	return NewParam(
		name,
		logger,
		func(value any) (any, error) { return cast.ToUint32E(value) },
		options...,
	)
}

// NewParamUint64 создаёт параметр типа uint64.
func NewParamUint64(name string, logger compogo.Logger, options ...Option) *Param {
	return NewParam(
		name,
		logger,
		func(value any) (any, error) { return cast.ToUint64E(value) },
		options...,
	)
}

// NewParamBool создаёт параметр типа bool.
func NewParamBool(name string, logger compogo.Logger, options ...Option) *Param {
	return NewParam(
		name,
		logger,
		func(value any) (any, error) { return cast.ToBoolE(value) },
		options...,
	)
}

// NewParamDuration создаёт параметр типа time.Duration.
func NewParamDuration(name string, logger compogo.Logger, options ...Option) *Param {
	return NewParam(
		name,
		logger,
		func(value any) (any, error) { return cast.ToDurationE(value) },
		options...,
	)
}

// NewParamTime создаёт параметр типа time.Time.
// Поддерживает различные форматы дат (RFC3339, ISO8601, и т.д.).
func NewParamTime(name string, logger compogo.Logger, options ...Option) *Param {
	return NewParam(
		name,
		logger,
		func(value any) (any, error) {
			val, err := cast.ToStringE(value)
			if err != nil {
				return nil, err
			}

			return dateparse.ParseStrict(val)
		},
		options...,
	)
}

// NewIp создаёт параметр для извлечения IP-адреса клиента.
// Последовательно проверяет заголовки:
//   - X-REAL-IP (nginx)
//   - X-FORWARDED-FOR (proxy/load balancer)
//   - RemoteAddr (прямое соединение)
//
// Возвращает net.IP.
func NewIp(name string, logger compogo.Logger) *Param {
	return NewParam(
		name,
		logger,
		IpCaster,
		WithHeaderGetterByName(HeaderRealIp),
		WithHeaderGetterByName(HeaderForwardedFor),
		AddGetter(func(request *http.Request) string {
			ip, _, err := net.SplitHostPort(request.RemoteAddr)
			if err != nil {
				logger.Error(err)
			}

			return ip
		}),
	)
}

// IpCaster преобразует строку в net.IP.
// Поддерживает X-Forwarded-For с несколькими IP (берёт первый валидный).
func IpCaster(value any) (any, error) {
	val, err := cast.ToStringE(value)
	if err != nil {
		return nil, err
	}

	splitIps := strings.Split(val, ipSeparator)
	for _, item := range splitIps {
		ip := net.ParseIP(item)
		if ip != nil {
			return ip, nil
		}
	}

	return nil, fmt.Errorf("caster.ip: ip '%s' is invalid", value)
}
