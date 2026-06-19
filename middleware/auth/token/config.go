package token

import (
	"bufio"
	"os"

	"github.com/Compogo/compogo"
	"github.com/Compogo/types/set"
)

const (
	// TokensFieldName — имя поля для списка токенов.
	TokensFieldName = "server.http.auth.token.tokens"

	// HeaderNameFieldName — имя поля для названия заголовка.
	HeaderNameFieldName = "server.http.auth.token.header"

	// FilePathFieldName — имя поля для пути к файлу с токенами.
	FilePathFieldName = "server.http.auth.token.filepath"
)

// HeaderNameDefault — имя заголовка по умолчанию (X-Auth-Token).
var HeaderNameDefault = "X-Auth-Token"

// Config содержит конфигурацию Token Auth.
type Config struct {
	HeaderName string
	FilePath   string

	tokens []string
	Tokens set.Set[string]
}

// NewConfig создаёт новую конфигурацию.
func NewConfig() *Config {
	return &Config{}
}

// Configuration загружает конфигурацию из Configurator.
// Поддерживает загрузку токенов из:
//   - Параметра tokens (срез строк)
//   - Файла (построчное чтение, каждая строка = токен)
//
// Пустые строки в файле игнорируются.
func Configuration(config *Config, configurator compogo.Configurator) (*Config, error) {
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
