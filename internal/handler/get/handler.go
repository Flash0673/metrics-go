package get

import (
	"errors"
	"net/http"

	models "github.com/Flash0673/metrics-go/internal/model"
	"github.com/Flash0673/metrics-go/internal/repository/repo_err"
	"github.com/go-chi/chi/v5"
)

type Service interface {
	Get(name, mType string) (*models.Metrics, error)
}

type Handler struct {
	svc Service
}

func NewHandler(svc Service) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t := chi.URLParam(r, "type")
	n := chi.URLParam(r, "name")

	if n == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	switch t {
	case "gauge":
	case "counter":
	default:
		w.WriteHeader(http.StatusBadRequest)
		return

	}

	metrics, err := h.svc.Get(n, t)
	if err != nil {
		if errors.Is(err, repo_err.ErrNotFound) {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Write([]byte(metrics.GetValue()))
}
