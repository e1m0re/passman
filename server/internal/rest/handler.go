package rest

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

type Handler struct {
}

func (h *Handler) NewRouter() *chi.Mux {
	r := chi.NewRouter()

	r.Route("/", func(r chi.Router) {
		r.Get("/", func(writer http.ResponseWriter, request *http.Request) {
			writer.Write([]byte("hello world"))
		})
	})

	return r
}

func NewHandler() *Handler {
	return &Handler{}
}
