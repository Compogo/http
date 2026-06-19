# Compogo HTTP Server

[![Go Reference](https://pkg.go.dev/badge/github.com/Compogo/http_server.svg)](https://pkg.go.dev/github.com/Compogo/http_server)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

HTTP-сервер и набор middleware для фреймворка [Compogo](https://github.com/Compogo/compogo).

Предоставляет:

* HTTP-сервер с graceful shutdown
* Маршрутизацию с поддержкой групп и middleware
* Аутентификацию (Basic Auth, Token Auth)
* Логирование запросов и ответов
* Метрики Prometheus (количество запросов, длительность)
* Извлечение и валидацию параметров запроса
* WebSocket (на базе gorilla/websocket)

## Установка

```shell
go get github.com/Compogo/http
```

## Быстрый старт

```go
package main

import (
    "net/http"

    "github.com/Compogo/compogo"
	httpServer "github.com/Compogo/http_server"
)

func main() {
    app := compogo.NewApp("myapp",
        // Базовые компоненты
        compogo.WithComponents(&httpServer.Component),
    )

    // Настройка роутера
    app.AddComponents(&compogo.Component{
        Name: "router",
        Init: compogo.StepFunc(func(container compogo.Container) error {
            return container.Invoke(func(server httpServer.Server) error {
                router := // ваша реализация Router (например, chi)
                server.SetRouter(router)
                
                router.Get("/health", func(w http.ResponseWriter, r *http.Request) {
                    w.Write([]byte("ok"))
                })
                
                return nil
            })
        }),
    })

    if err := app.Serve(); err != nil {
        panic(err)
    }
}
```

## Компоненты

### HTTP Server

HTTP-сервер с поддержкой graceful shutdown.

```go
// Компонент
compogo.WithComponents(&httpServer.Component)

// Настройка через флаги
// --server.http.interface=0.0.0.0
// --server.http.port=8080
// --server.http.timeout.shutdown=30s
```

### Аутентификация

#### Basic Auth

```go
import "github.com/Compogo/http_server/middleware/basic"

// Подключение компонента
app.AddComponents(&basic.Component)

// Настройка через флаги
// --server.http.auth.basic.creds=admin:password,user:123
// --server.http.auth.basic.filepath=/path/to/creds.txt

// Использование в роутере
router.Group(func(r http.Router) {
    r.Use(auth)
    r.Get("/admin", adminHandler)
})
```

#### Token Auth

```go
import "github.com/Compogo/http_server/middleware/token"

// Подключение компонента
app.AddComponents(&token.Component)

// Настройка через флаги
// --server.http.auth.token.tokens=abc123,def456
// --server.http.auth.token.header=X-Auth-Token
// --server.http.auth.token.filepath=/path/to/tokens.txt
```

### Логирование

```go
import "github.com/Compogo/http_server/middleware/logger"

// Подключение компонентов
app.AddComponents(&logger.RequestComponent)
app.AddComponents(&logger.ResponseComponent)

// Логирует путь и тело запроса/ответа на уровне Debug
```

### Метрики Prometheus

```go
import "github.com/Compogo/http_server/middleware/metric"

// Подключение компонентов
app.AddComponents(&metric.RequestCountComponent)
app.AddComponents(&metric.DurationComponent)

// Метрики:
// compogo_http_server_requests_total{app="myapp", code="200", endpoint="/api/users"}
// compogo_http_server_duration_seconds{app="myapp", endpoint="/api/users"}
```

### Параметры запроса

```go
import "github.com/Compogo/http_server/middleware/param"

// Создание параметра
userId := param.NewParamInt("user_id", logger,
    param.WithUriGetter(),
    param.AddValidator(param.GteValidator(1)),
)

// Использование в роутере
router.Get("/users/:id", userId.Middleware(handler))

// Получение из контекста
userID := r.Context().Value("user_id").(int)
```

### WebSocket

```go
import "github.com/Compogo/http_server/websocket"

// Подключение компонента
app.AddComponents(&websocket.Component)

// Создание обработчика
handler := websocket.NewHandler(config, upgrader, logger, closer)

// Подписка на события
handler.OnClientConnection.Subscribe(func(ctx context.Context, client *websocket.Client) {
    logger.Info("client connected")
})

handler.OnMessage.Subscribe(func(ctx context.Context, event *websocket.Event) {
    // обработка входящего сообщения
})

// Регистрация в роутере
router.Get("/ws", handler.ServeHTTP)

// Отправка всем клиентам
handler.Send(websocket.NewEvent("message", payload))
```

## Router

Пакет определяет интерфейс `Router`, который может быть реализован любой библиотекой маршрутизации:

```go
type Router interface {
    http.Handler
    Use(middlewares ...Middleware)
    Group(fn func(r Router))
    Route(pattern string, fn func(r Router))
    Mount(pattern string, h http.Handler)
    Handle(pattern string, h http.Handler)
    HandleFunc(pattern string, h http.HandlerFunc)
    Method(method, pattern string, h http.Handler)
    MethodFunc(method, pattern string, h http.HandlerFunc)
    Get(pattern string, h http.HandlerFunc)
    Post(pattern string, h http.HandlerFunc)
    Put(pattern string, h http.HandlerFunc)
    Delete(pattern string, h http.HandlerFunc)
    // ... и другие методы
}
```

## Middleware

Написание своего middleware

```go
type MyMiddleware struct {
    logger compogo.Logger
}

func (m *MyMiddleware) Middleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // before
        m.logger.Info("request started")
        
        next.ServeHTTP(w, r)
        
        // after
        m.logger.Info("request finished")
    })
}

// Использование
router.Use(&MyMiddleware{logger: logger})
```

## Helper

```go
// Отправка ошибки в JSON или plain text
helper.WriteError(w, r, "error message", http.StatusBadRequest)

// Отправка ошибки в JSON
helper.JSONError(w, helper.NewError("error message"), http.StatusBadRequest)
```

## Зависимости

* [Compogo](https://github.com/Compogo/compogo) — основной фреймворк
* [Compogo Runner](https://github.com/Compogo/runner) — управление фоновыми процессами
* [gorilla/websocket](https://github.com/gorilla/websocket) — WebSocket
* [prometheus/client_golang](https://github.com/prometheus/client_golang) — метрики
* [spf13/cast](https://github.com/spf13/cast) — приведение типов

## Лицензия

```plantuml
MIT License

Copyright (c) 2026 Compogo

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.

```
