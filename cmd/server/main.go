package main

import (
  "log"
  "time"
  "net/http"
  "github.com/AjayPoshak/url-shortener/internal/handlers"
  "github.com/AjayPoshak/url-shortener/internal/middleware"
)

func main() {
  // Create a new router
  mux := http.NewServeMux()

  // Initialize the handlers
  handlers := handlers.NewHandlers()

  // Register routes with middleware
  mux.Handle("/home", middleware.Logger(handlers.HomeHandler))
  mux.Handle("/health", middleware.Logger(handlers.HealthHandler))

  port := ":8095"
  // Create a new server
  server := &http.Server{
    Addr:        port,
    Handler:      mux,
    ReadTimeout:  5 * time.Second,
    WriteTimeout: 10 * time.Second,
  }

  // Start the server
  log.Println("Starting server on ", port)
  if err := server.ListenAndServe(); err != nil {
    log.Fatalf("Server failed to start: %v", err)
  }
}
// test comment
