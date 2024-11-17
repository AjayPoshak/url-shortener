package handlers

import (
  "net/http"
)

type Handlers struct {}

func NewHandlers() *Handlers {
  return &Handlers{}
}

func (h *Handlers) HomeHandler(w http.ResponseWriter, r *http.Request) {
  w.Write([]byte("Welcome to the home page 12345"))
}

func (h *Handlers) HealthHandler(w http.ResponseWriter, r *http.Request) {
  w.Write([]byte("OK"))
}
