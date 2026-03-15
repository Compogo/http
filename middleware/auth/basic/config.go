package basic

import (
	"bufio"
	"os"
	"strings"

	"github.com/Compogo/compogo/configurator"
	"github.com/Compogo/types/mapper"
)

const (
	// CredsFieldName is the command-line flag for inline credentials.
	// Format: "username1:password1,username2:password2"
	CredsFieldName = "server.http.auth.basic.creds"

	// FilePathFieldName is the command-line flag for credentials file.
	// File should contain one "username:password" per line.
	FilePathFieldName = "server.http.auth.basic.filepath"

	pairsSeparator = ":"
)

// Cred represents a single username/password pair for basic authentication.
type Cred struct {
	UserName string
	Password string
}

// String returns the username, implementing fmt.Stringer.
func (c *Cred) String() string {
	return c.UserName
}

// Config holds the basic authentication configuration.
// It can be populated from command-line flags, config files, or a credentials file.
type Config struct {
	FilePath string

	creds []string
	Creds mapper.Mapper[*Cred]
}

// NewConfig creates a new Config instance with default values.
func NewConfig() *Config {
	return &Config{}
}

// Configuration applies configuration values to the Config struct.
// It reads from configurator and optionally from a credentials file.
// Returns an error if file reading fails.
func Configuration(config *Config, configurator configurator.Configurator) (*Config, error) {
	if len(config.creds) == 0 {
		config.creds = configurator.GetStringSlice(CredsFieldName)
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

			config.creds = append(config.creds, line)
		}

		if err := scanner.Err(); err != nil {
			return nil, err
		}
	}

	for _, cred := range config.creds {
		userNamePassword := strings.SplitN(cred, pairsSeparator, 2)

		config.Creds.Add(&Cred{UserName: userNamePassword[0], Password: userNamePassword[1]})
	}

	return config, nil
}
