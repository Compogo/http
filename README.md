# Compogo HTTP 🌐

**Compogo HTTP** — минималистичный HTTP-сервер для Compogo, построенный
поверх [Runner](https://github.com/Compogo/runner). Сервер не содержит встроенного роутинга — вы подключаете любой
роутер, реализующий простой интерфейс `Router`.

## 🚀 Установка

```bash
go get github.com/Compogo/http
```

### 📦 Быстрый старт

```go
package main

import (
	"net/http"

	"github.com/Compogo/compogo"
	"github.com/Compogo/runner"
	"github.com/Compogo/http"
	"github.com/go-chi/chi/v5" // любой роутер
)

func main() {
	app := compogo.NewApp("myapp",
		compogo.WithOsSignalCloser(),
		runner.WithRunner(),
		http.WithServer(), // добавляем HTTP-сервер
		compogo.WithComponents(
			routerComponent,
		),
	)

	if err := app.Serve(); err != nil {
		panic(err)
	}
}

// Компонент с роутером
var routerComponent = &component.Component{
	Dependencies: component.Components{http.Component},
	Init: component.StepFunc(func(c container.Container) error {
		return c.Provide(newRouter)
	}),
	PostRun: component.StepFunc(func(c container.Container) error {
		return c.Invoke(func(r chi.Router, s http.Server) {
			// регистрируем маршруты
			r.Get("/hello", func(w http.ResponseWriter, r *http.Request) {
				w.Write([]byte("Hello, World!"))
			})

			// подключаем к серверу
			s.SetRouter(r)
		})
	}),
}

func newRouter() chi.Router {
	return chi.NewRouter()
}
```

### ✨ Возможности

#### 🎯 Чистый HTTP-сервер без лишнего

* Только слушает порт и отдаёт трафик в роутер
* Graceful shutdown с таймаутом
* Интеграция с Runner'ом
* Конфигурация через флаги

#### 🔌 Поддержка любых роутеров

Интерфейс `Router` совместим с:

* [chi](https://github.com/go-chi/chi)
* [gorilla/mux](https://github.com/gorilla/mux)
* http.ServeMux
* Любым другим, реализующим те же методы

#### 🧩 Middleware

```go
type LoggerMiddleware struct{}

func (m *LoggerMiddleware) Middleware(next http.Handler) http.Handler {
return http.HandlerFunc(func (w http.ResponseWriter, r *http.Request) {
log.Printf("%s %s", r.Method, r.URL.Path)
next.ServeHTTP(w, r)
})
}

// Использование в роутере
router.Use(&LoggerMiddleware{})
```

#### ⚙️ Конфигурация

```bash
./myapp \
    --server.http.interface=0.0.0.0 \
    --server.http.port=8080 \
    --server.http.timeout.shutdown=30s
```

| Флаг                           | По умолчанию | Описание                              |
|--------------------------------|--------------|---------------------------------------|
| --server.http.interface        | 0.0.0.0      | Интерфейс для прослушивания           |
| --server.http.port             | 8080         | Порт                                  |
| --server.http.timeout.shutdown | 30s          | Время на завершение активных запросов |

