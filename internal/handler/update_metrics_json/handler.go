package update_metrics_json

import (
	"io"
	"net/http"

	models "github.com/Flash0673/metrics-go/internal/model"
	"github.com/mailru/easyjson"
)

//go:generate mockgen -destination=mocks/mocks.go -package=mocks . Service

type Service interface {
	Set(name, mType string, value any) error
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

	var value any
	switch m.MType {
	case models.Gauge:
		if m.Value == nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		value = *m.Value
	case models.Counter:
		if m.Delta == nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		value = *m.Delta
	default:
		w.WriteHeader(http.StatusBadRequest)
		return

	}

	err = h.svc.Set(m.ID, m.MType, value)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
