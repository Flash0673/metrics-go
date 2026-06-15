package get_all_metrics

import (
	"fmt"
	"net/http"
	"strings"

	models "github.com/Flash0673/metrics-go/internal/model"
)

var respTmpl = `
<html>
<head>
<ul>
%s
<ul>
`

var listElTmpl = `<li>%s</li>`

type Service interface {
	GetAll() []*models.Metrics
}

type Handler struct {
	svc Service
}

func NewHandler(svc Service) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	metrics := h.svc.GetAll()
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	res := make([]string, 0, len(metrics))
	for _, m := range metrics {
		res = append(res, fmt.Sprintf(listElTmpl, fmt.Sprintf("%s: %v", m.GetName(), m.GetValue())))
	}
	w.Write([]byte(fmt.Sprintf(respTmpl, strings.Join(res, "\n"))))
}
