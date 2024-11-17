package middleware

import (
  "log"
  "net/http"
  "time"
)

type Middleware func(http.Handler) http.Handler

func Chain(handler http.Handler, middlewares ...Middleware) http.Handler {
  for _, middleware := range middlewares {
    handler = middleware(handler)
  }
  return handler
}

  func Logger(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()

        // Call the next handler
        next(w, r)

        // Log the request details
        log.Printf(
            "%s %s %s %v",
            r.Method,
            r.RequestURI,
            r.RemoteAddr,
            time.Since(start),
        )
    }
}
