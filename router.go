package http

import "net/http"

// Middleware defines the interface for HTTP middleware.
// Middlewares wrap http.Handler to add cross-cutting concerns
// such as logging, authentication, metrics, or recovery.
//
// Example:
//
//	func LoggerMiddleware(next http.Handler) http.Handler {
//	    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//	        log.Printf("%s %s", r.Method, r.URL.Path)
//	        next.ServeHTTP(w, r)
//	    })
//	}
type Middleware interface {
	// Middleware wraps the next handler and returns a new handler.
	// The returned handler should call next.ServeHTTP() to continue the chain.
	Middleware(next http.Handler) http.Handler
}

// MiddlewareFunc is a function adapter that allows ordinary functions
// to be used as Middleware implementations.
type MiddlewareFunc func(next http.Handler) http.Handler

// Middleware implements the Middleware interface by calling the underlying function.
func (m MiddlewareFunc) Middleware(next http.Handler) http.Handler {
	return m(next)
}

// Router defines the interface for HTTP routing.
// It's designed to be compatible with popular routers like chi, gorilla/mux,
// or the standard library's ServeMux.
//
// The interface includes:
//   - Standard HTTP handler methods (Get, Post, Put, etc.)
//   - Middleware support (Use)
//   - Route grouping (Group, Route)
//   - Subtree mounting (Mount)
type Router interface {
	// http.Handler makes Router usable as a standard HTTP handler.
	http.Handler

	// Use adds middleware to the router.
	// Middlewares are applied in the order they are added.
	Use(middlewares ...Middleware)

	// Group creates a new sub-router with inherited middleware.
	// All routes defined inside the function will share the same prefix.
	Group(fn func(r Router))

	// Route creates a new sub-router with the given prefix.
	// Similar to Group but without inheriting parent middleware.
	Route(pattern string, fn func(r Router))

	// Mount attaches another http.Handler at the specified pattern.
	// Useful for mounting sub-applications or third-party handlers.
	Mount(pattern string, h http.Handler)

	// Handle registers a handler for the given pattern with all HTTP methods.
	Handle(pattern string, h http.Handler)

	// HandleFunc registers a handler function for the given pattern with all HTTP methods.
	HandleFunc(pattern string, h http.HandlerFunc)

	// Method registers a handler for the given pattern with a specific HTTP method.
	Method(method, pattern string, h http.Handler)

	// MethodFunc registers a handler function for the given pattern with a specific HTTP method.
	MethodFunc(method, pattern string, h http.HandlerFunc)

	// Connect registers a handler for CONNECT method.
	Connect(pattern string, h http.HandlerFunc)

	// Delete registers a handler for DELETE method.
	Delete(pattern string, h http.HandlerFunc)

	// Get registers a handler for GET method.
	Get(pattern string, h http.HandlerFunc)

	// Head registers a handler for HEAD method.
	Head(pattern string, h http.HandlerFunc)

	// Options registers a handler for OPTIONS method.
	Options(pattern string, h http.HandlerFunc)

	// Patch registers a handler for PATCH method.
	Patch(pattern string, h http.HandlerFunc)

	// Post registers a handler for POST method.
	Post(pattern string, h http.HandlerFunc)

	// Put registers a handler for PUT method.
	Put(pattern string, h http.HandlerFunc)

	// Trace registers a handler for TRACE method.
	Trace(pattern string, h http.HandlerFunc)

	// NotFound sets the handler for routes that don't match any pattern.
	NotFound(h http.HandlerFunc)

	// MethodNotAllowed sets the handler for routes that match the path
	// but not the HTTP method.
	MethodNotAllowed(h http.HandlerFunc)
}
