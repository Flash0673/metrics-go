package update_metrics

import (
	"net/http"
	"strconv"

	models "github.com/Flash0673/metrics-go/internal/model"
	"github.com/go-chi/chi/v5"
)

type Service interface {
	Set(name, mType string, value any) (*models.Metrics, error)
}

type Handler struct {
	svc Service
}

func NewHandler(svc Service) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	t := chi.URLParam(r, "type")
	n := chi.URLParam(r, "name")
	v := chi.URLParam(r, "value")

	if n == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	var value any
	var err error
	switch t {
	case "gauge":
		value, err = strconv.ParseFloat(v, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	case "counter":
		value, err = strconv.ParseInt(v, 10, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	default:
		w.WriteHeader(http.StatusBadRequest)
		return

	}

	_, err = h.svc.Set(n, t, value)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
