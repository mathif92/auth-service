package handlers

import "net/http"

type Health struct{}

func NewHealth() *Health {
	return &Health{}
}

func (h *Health) Ping(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK"))
}
