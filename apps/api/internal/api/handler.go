package api

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/merttmutlu/notes-stack/apps/api/internal/note"
	"github.com/merttmutlu/notes-stack/apps/api/internal/store"
)

type Handler struct {
	logger *slog.Logger
	store  *store.Store
}

func (h Handler) Healthz(w http.ResponseWriter, r *http.Request) {
	_ = r
	writeJSON(w, http.StatusOK, map[string]string{
		"status": "ok",
	})
}

func (h Handler) Readyz(w http.ResponseWriter, r *http.Request) {
	if err := h.store.Ping(r.Context()); err != nil {
		if h.logger != nil {
			h.logger.Warn("readiness check failed", "error", err)
		}

		writeJSON(w, http.StatusServiceUnavailable, map[string]string{
			"status": "not_ready",
		})
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{
		"status": "ready",
	})
}

func (h Handler) ListNotes(w http.ResponseWriter, r *http.Request) {
	notes, err := h.store.ListNotes(r.Context())
	if err != nil {
		if h.logger != nil {
			h.logger.Error("failed to list notes", "error", err)
		}

		writeJSON(w, http.StatusInternalServerError, map[string]string{
			"error": "internal server error",
		})
		return
	}

	writeJSON(w, http.StatusOK, struct {
		Data []note.Note `json:"data"`
	}{
		Data: notes,
	})
}

func (h Handler) GetNote(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimSpace(r.PathValue("id"))
	if id == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{
			"error": "note id is required",
		})
		return
	}

	item, err := h.store.GetNote(r.Context(), id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			writeJSON(w, http.StatusNotFound, map[string]string{
				"error": "note not found",
			})
			return
		}

		if h.logger != nil {
			h.logger.Error("failed to get note", "error", err, "id", id)
		}

		writeJSON(w, http.StatusInternalServerError, map[string]string{
			"error": "internal server error",
		})
		return
	}

	writeJSON(w, http.StatusOK, struct {
		Data note.Note `json:"data"`
	}{
		Data: item,
	})
}

func (h Handler) CreateNote(w http.ResponseWriter, r *http.Request) {
	var req note.CreateRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{
			"error": "invalid request body",
		})
		return
	}

	if err := validateCreateRequest(req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
		return
	}

	created, err := h.store.CreateNote(r.Context(), req)
	if err != nil {
		if h.logger != nil {
			h.logger.Error("failed to create note", "error", err)
		}

		writeJSON(w, http.StatusInternalServerError, map[string]string{
			"error": "internal server error",
		})
		return
	}

	writeJSON(w, http.StatusCreated, struct {
		Data note.Note `json:"data"`
	}{
		Data: created,
	})
}

func (h Handler) UpdateNote(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimSpace(r.PathValue("id"))
	if id == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{
			"error": "note id is required",
		})
		return
	}

	var req note.UpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{
			"error": "invalid request body",
		})
		return
	}

	if err := validateUpdateRequest(req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
		return
	}

	updated, err := h.store.UpdateNote(r.Context(), id, req)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			writeJSON(w, http.StatusNotFound, map[string]string{
				"error": "note not found",
			})
			return
		}

		if h.logger != nil {
			h.logger.Error("failed to update note", "error", err, "id", id)
		}

		writeJSON(w, http.StatusInternalServerError, map[string]string{
			"error": "internal server error",
		})
		return
	}

	writeJSON(w, http.StatusOK, struct {
		Data note.Note `json:"data"`
	}{
		Data: updated,
	})
}

func (h Handler) DeleteNote(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimSpace(r.PathValue("id"))
	if id == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{
			"error": "note id is required",
		})
		return
	}

	deleted, err := h.store.DeleteNote(r.Context(), id)
	if err != nil {
		if h.logger != nil {
			h.logger.Error("failed to delete note", "error", err, "id", id)
		}

		writeJSON(w, http.StatusInternalServerError, map[string]string{
			"error": "internal server error",
		})
		return
	}

	if !deleted {
		writeJSON(w, http.StatusNotFound, map[string]string{
			"error": "note not found",
		})
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func validateCreateRequest(req note.CreateRequest) error {
	if strings.TrimSpace(req.Title) == "" || strings.TrimSpace(req.Content) == "" {
		return errors.New("title and content are required")
	}

	return nil
}

func validateUpdateRequest(req note.UpdateRequest) error {
	if strings.TrimSpace(req.Title) == "" || strings.TrimSpace(req.Content) == "" {
		return errors.New("title and content are required")
	}

	return nil
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(payload); err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
}
