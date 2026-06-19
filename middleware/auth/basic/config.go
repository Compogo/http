package basic

import (
	"bufio"
	"os"
	"strings"

	"github.com/Compogo/compogo"
	"github.com/Compogo/types/mapper"
)

const (
	// CredsFieldName — имя поля для списка учетных данных (логин:пароль).
	CredsFieldName = "server.http.auth.basic.creds"

	// FilePathFieldName — имя поля для пути к файлу с учетными данными.
	FilePathFieldName = "server.http.auth.basic.filepath"

	// pairsSeparator — разделитель между логином и паролем.
	pairsSeparator = ":"
)

// Cred представляет учетные данные пользователя.
// Реализует fmt.Stringer для использования в mapper.
type Cred struct {
	UserName string
	Password string
}

// String возвращает имя пользователя.
// Реализует интерфейс fmt.Stringer для mapper.
func (c *Cred) String() string {
	return c.UserName
}

// Config содержит конфигурацию Basic Auth.
type Config struct {
	FilePath string

	creds []string
	Creds mapper.Mapper[*Cred]
}

// NewConfig создаёт новую конфигурацию.
func NewConfig() *Config {
	return &Config{}
}

// Configuration загружает конфигурацию из Configurator.
// Поддерживает загрузку учетных данных из:
//   - Параметра creds (срез строк "логин:пароль")
//   - Файла (построчное чтение, каждая строка = "логин:пароль")
//
// Пустые строки в файле игнорируются.
func Configuration(config *Config, configurator compogo.Configurator) (*Config, error) {
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
