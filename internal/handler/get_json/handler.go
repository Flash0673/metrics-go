package get_json

import (
	"errors"
	"io"
	"net/http"

	models "github.com/Flash0673/metrics-go/internal/model"
	"github.com/Flash0673/metrics-go/internal/repository/repo_err"
	"github.com/mailru/easyjson"
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
	defer r.Body.Close()
	raw, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	m := &models.Metrics{}
	err = easyjson.Unmarshal(raw, m)
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	if m.ID == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	switch m.MType {
	case "gauge":
	case "counter":
	default:
		w.WriteHeader(http.StatusBadRequest)
		return

	}

	metrics, err := h.svc.Get(m.ID, m.MType)
	if err != nil {
		if errors.Is(err, repo_err.ErrNotFound) {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	resp, err := easyjson.Marshal(metrics)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(resp)
}
