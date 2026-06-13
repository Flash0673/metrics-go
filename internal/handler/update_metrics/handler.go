package update_metrics

import (
	"net/http"
	"strconv"
)

type Handler struct{}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	t := r.PathValue("type")
	n := r.PathValue("name")
	v := r.PathValue("value")

	if n == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	switch t {
	case "gauge":
		_, err := strconv.ParseFloat(v, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	case "counter":
		_, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	default:
		w.WriteHeader(http.StatusBadRequest)
		return

	}

	w.WriteHeader(http.StatusOK)
}
