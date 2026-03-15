package token

import (
	"bufio"
	"os"

	"github.com/Compogo/compogo/configurator"
	"github.com/Compogo/types/set"
)

const (
	// TokensFieldName is the command-line flag for inline tokens.
	// Format: "token1,token2,token3"
	TokensFieldName = "server.http.auth.token.tokens"

	// HeaderNameFieldName is the command-line flag for the header name.
	HeaderNameFieldName = "server.http.auth.token.header"

	// FilePathFieldName is the command-line flag for tokens file.
	// File should contain one token per line.
	FilePathFieldName = "server.http.auth.token.filepath"

	// HeaderNameDefault is the default header name for token authentication.
	HeaderNameDefault = "X-Auth-Token"
)

// Config holds the token authentication configuration.
// It can be populated from command-line flags, config files, or a tokens file.
type Config struct {
	// HeaderName is the HTTP header to look for the token.
	HeaderName string

	// FilePath is the path to a file containing allowed tokens.
	FilePath string

	tokens []string
	// Set of allowed tokens for O(1) lookup
	Tokens set.Set[string]
}

// NewConfig creates a new Config instance with default values.
func NewConfig() *Config {
	return &Config{}
}

// Configuration applies configuration values to the Config struct.
// It reads from configurator and optionally from a tokens file.
// Returns an error if file reading fails.
func Configuration(config *Config, configurator configurator.Configurator) (*Config, error) {
	if config.HeaderName == "" || config.HeaderName == HeaderNameDefault {
		configurator.SetDefault(HeaderNameFieldName, HeaderNameDefault)
		config.HeaderName = configurator.GetString(HeaderNameFieldName)
	}

	if len(config.tokens) == 0 {
		config.tokens = configurator.GetStringSlice(TokensFieldName)
	}

	for _, token := range config.tokens {
		config.Tokens.Add(token)
	}

	if config.FilePath == "" {
		config.FilePath = configurator.GetString(FilePathFieldName)
	}

	if config.FilePath != "" {
		f, err := os.Open(config.FilePath)
		if err != nil {
			return nil, err
		}

		defer f.Close()

		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
			line := scanner.Text()
			if line == "" {
				continue
			}

			config.Tokens.Add(line)
		}

		if err := scanner.Err(); err != nil {
			return nil, err
		}
	}

	return config, nil
}
