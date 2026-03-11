package param

import (
	"fmt"
	"net"
	"net/http"
	"strings"

	"github.com/Compogo/compogo/logger"
	"github.com/araddon/dateparse"
	"github.com/spf13/cast"
)

const (
	// HeaderRealIp is the standard header for real client IP when behind proxies.
	HeaderRealIp = "X-REAL-IP"

	// HeaderForwardedFor is the standard header for forwarded client IPs.
	HeaderForwardedFor = "X-FORWARDED-FOR"

	// ipSeparator separates multiple IPs in X-Forwarded-For header.
	ipSeparator = ","
)

// NewParamString creates a new string parameter.
func NewParamString(name string, logger logger.Logger, options ...Option) *Param {
	return NewParam(
		name,
		logger,
		func(value any) (any, error) { return cast.ToStringE(value) },
		options...,
	)
}

// NewParamInt creates a new int parameter.
func NewParamInt(name string, logger logger.Logger, options ...Option) *Param {
	return NewParam(
		name,
		logger,
		func(value any) (any, error) { return cast.ToIntE(value) },
		options...,
	)
}

// NewParamInt8 creates a new int8 parameter.
func NewParamInt8(name string, logger logger.Logger, options ...Option) *Param {
	return NewParam(
		name,
		logger,
		func(value any) (any, error) { return cast.ToInt8E(value) },
		options...,
	)
}

// NewParamInt16 creates a new int16 parameter.
func NewParamInt16(name string, logger logger.Logger, options ...Option) *Param {
	return NewParam(
		name,
		logger,
		func(value any) (any, error) { return cast.ToInt16E(value) },
		options...,
	)
}

// NewParamInt32 creates a new int32 parameter.
func NewParamInt32(name string, logger logger.Logger, options ...Option) *Param {
	return NewParam(
		name,
		logger,
		func(value any) (any, error) { return cast.ToInt32E(value) },
		options...,
	)
}

// NewParamInt64 creates a new int64 parameter.
func NewParamInt64(name string, logger logger.Logger, options ...Option) *Param {
	return NewParam(
		name,
		logger,
		func(value any) (any, error) { return cast.ToInt64E(value) },
		options...,
	)
}

// NewParamFloat32 creates a new float32 parameter.
func NewParamFloat32(name string, logger logger.Logger, options ...Option) *Param {
	return NewParam(
		name,
		logger,
		func(value any) (any, error) { return cast.ToFloat32E(value) },
		options...,
	)
}

// NewParamFloat64 creates a new float64 parameter.
func NewParamFloat64(name string, logger logger.Logger, options ...Option) *Param {
	return NewParam(
		name,
		logger,
		func(value any) (any, error) { return cast.ToFloat64E(value) },
		options...,
	)
}

// NewParamUint creates a new uint parameter.
func NewParamUint(name string, logger logger.Logger, options ...Option) *Param {
	return NewParam(
		name,
		logger,
		func(value any) (any, error) { return cast.ToUintE(value) },
		options...,
	)
}

// NewParamUint8 creates a new uint8 parameter.
func NewParamUint8(name string, logger logger.Logger, options ...Option) *Param {
	return NewParam(
		name,
		logger,
		func(value any) (any, error) { return cast.ToUint8E(value) },
		options...,
	)
}

// NewParamUint16 creates a new uint16 parameter.
func NewParamUint16(name string, logger logger.Logger, options ...Option) *Param {
	return NewParam(
		name,
		logger,
		func(value any) (any, error) { return cast.ToUint16E(value) },
		options...,
	)
}

// NewParamUint32 creates a new uint32 parameter.
func NewParamUint32(name string, logger logger.Logger, options ...Option) *Param {
	return NewParam(
		name,
		logger,
		func(value any) (any, error) { return cast.ToUint32E(value) },
		options...,
	)
}

// NewParamUint64 creates a new uint64 parameter.
func NewParamUint64(name string, logger logger.Logger, options ...Option) *Param {
	return NewParam(
		name,
		logger,
		func(value any) (any, error) { return cast.ToUint64E(value) },
		options...,
	)
}

// NewParamBool creates a new boolean parameter.
func NewParamBool(name string, logger logger.Logger, options ...Option) *Param {
	return NewParam(
		name,
		logger,
		func(value any) (any, error) { return cast.ToBoolE(value) },
		options...,
	)
}

// NewParamDuration creates a new time.Duration parameter.
func NewParamDuration(name string, logger logger.Logger, options ...Option) *Param {
	return NewParam(
		name,
		logger,
		func(value any) (any, error) { return cast.ToDurationE(value) },
		options...,
	)
}

// NewParamTime creates a new time.Time parameter.
// Supports flexible date parsing via github.com/araddon/dateparse.
func NewParamTime(name string, logger logger.Logger, options ...Option) *Param {
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

// NewIp creates a parameter for extracting client IP addresses.
// It checks X-REAL-IP, X-FORWARDED-FOR, and finally RemoteAddr.
// Returns a net.IP value.
func NewIp(name string, logger logger.Logger) *Param {
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

// IpCaster converts a string to net.IP.
// Handles comma-separated lists (X-Forwarded-For) and returns the first valid IP.
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
