package http

import (
	"net/http"
)

func NewRouter(h *Handler) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /orders", h.CreateOrder)
	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusOK) })
	return mux
}

