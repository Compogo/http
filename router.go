package http_server

import "net/http"

// Middleware представляет функцию-обёртку для HTTP-обработчиков.
// Позволяет добавлять логирование, авторизацию, метрики и т.д.
type Middleware interface {
	// Middleware оборачивает http.Handler в цепочку обработки.
	Middleware(next http.Handler) http.Handler
}

// MiddlewareFunc — функциональный адаптер для Middleware.
// Позволяет использовать функции как middleware.
//
// Пример:
//
//	middleware := MiddlewareFunc(func(next http.Handler) http.Handler {
//	    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//	        // before
//	        next.ServeHTTP(w, r)
//	        // after
//	    })
//	})
type MiddlewareFunc func(next http.Handler) http.Handler

// Middleware реализует интерфейс Middleware для MiddlewareFunc.
func (m MiddlewareFunc) Middleware(next http.Handler) http.Handler {
	return m(next)
}

// Router определяет интерфейс HTTP-роутера с поддержкой:
//   - Группировки маршрутов
//   - Middleware
//   - Стандартных HTTP-методов (GET, POST, PUT, DELETE, и т.д.)
//   - Вложенных роутеров (Mount)
//
// Интерфейс абстрагирует конкретную реализацию роутера
// (например, chi, gorilla/mux).
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
	Connect(pattern string, h http.HandlerFunc)
	Delete(pattern string, h http.HandlerFunc)
	Get(pattern string, h http.HandlerFunc)
	Head(pattern string, h http.HandlerFunc)
	Options(pattern string, h http.HandlerFunc)
	Patch(pattern string, h http.HandlerFunc)
	Post(pattern string, h http.HandlerFunc)
	Put(pattern string, h http.HandlerFunc)
	Trace(pattern string, h http.HandlerFunc)
	NotFound(h http.HandlerFunc)
	MethodNotAllowed(h http.HandlerFunc)
}
