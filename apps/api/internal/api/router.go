package api

import (
	"log/slog"
	"net/http"

	"github.com/merttmutlu/notes-stack/apps/api/internal/store"
)

type Dependencies struct {
	Logger *slog.Logger
	Store  *store.Store
}

func NewRouter(deps Dependencies) http.Handler {
	handler := Handler{
		logger: deps.Logger,
		store:  deps.Store,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", handler.Healthz)
	mux.HandleFunc("GET /readyz", handler.Readyz)
	mux.HandleFunc("GET /api/v1/notes", handler.ListNotes)
	mux.HandleFunc("GET /api/v1/notes/{id}", handler.GetNote)
	mux.HandleFunc("POST /api/v1/notes", handler.CreateNote)
	mux.HandleFunc("PUT /api/v1/notes/{id}", handler.UpdateNote)
	mux.HandleFunc("DELETE /api/v1/notes/{id}", handler.DeleteNote)

	return mux
}
