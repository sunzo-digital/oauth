package main

import (
	"net/http"
	"time"
)

func main() {
	handler := &Handler{
		allowedClients: map[string]*Client{
			"1111": &Client{
				id:     "1111",
				secret: "0000",
			},
		},
	}

	server := &http.Server{
		Addr:         ":8080",
		Handler:      handler,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	panic(server.ListenAndServe())
}

// Handler auth + resources server handler
type Handler struct {
	allowedClients map[string]*Client
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/grant":
		h.Grant(w, r)
	case "/token":
		h.Token(w, r)
	default:
		h.NotFound(w, r)
	}
}

func (h *Handler) Grant(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (h *Handler) Token(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (h *Handler) NotFound(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
}
